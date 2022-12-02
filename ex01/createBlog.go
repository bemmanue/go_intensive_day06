package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var db *sql.DB

type Article struct {
	Id      int
	Article string
}

func main() {
	driverName := "postgres"
	dataSourceName := "user=bemmanue dbname=bemmanue sslmode=disable"
	var err error

	db, err = sql.Open(driverName, dataSourceName)
	if err != nil {
		log.Fatalln(err)
	}

	defer db.Close()

	http.HandleFunc("/", showPage)
	http.HandleFunc("/admin", postArticle)
	http.HandleFunc("/posting", addArticle)

	styleHandler := http.FileServer(http.Dir("./css"))
	http.Handle("/css/", http.StripPrefix("/css", styleHandler))

	err = http.ListenAndServe("localhost:8888", nil)
	if err != nil {
		log.Fatalln(err)
	}
}

func showPage(w http.ResponseWriter, r *http.Request) {
	limit, offset := 3, 0

	if strings.HasPrefix(r.URL.RawQuery, "page=") {
		page, err := strconv.Atoi(strings.TrimPrefix(r.URL.RawQuery, "page="))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		offset = limit * (page - 1)
	}

	articles := getArticles(limit, offset)

	data := struct {
		Articles []Article
		Next     string
	}{Articles: articles, Next: "2"}

	tmpl, _ := template.ParseFiles("templates/index.html")

	err := tmpl.Execute(w, data)
	if err != nil {
		log.Fatalln(err)
	}
}

func getArticles(limit, offset int) []Article {
	query := "SELECT * FROM blog ORDER BY id LIMIT " + strconv.Itoa(limit) +
		" OFFSET " + strconv.Itoa(offset)

	rows, err := db.Query(query)
	if err != nil {
		log.Fatalln(err)
	}

	defer rows.Close()

	var articles []Article
	for rows.Next() {
		article := Article{}
		rows.Scan(&article.Id, &article.Article)
		articles = append(articles, article)
	}

	return articles
}

func postArticle(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "html/admin.html")
}

func addArticle(w http.ResponseWriter, r *http.Request) {
	article := r.FormValue("article")
	query := "insert into blog (article) values ('" + article + "')"

	_, err := db.Exec(query)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Fprintln(w, "The article was successfully added")
}
