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

type Message struct {
	Id      int
	Message string
}

func main() {
	http.HandleFunc("/", showBlog)
	http.HandleFunc("/admin", postArticles)
	http.HandleFunc("/publication", publishArticle)
	http.HandleFunc("/?page=", getPage)

	styleHandler := http.FileServer(http.Dir("./css"))
	http.Handle("/css/", http.StripPrefix("/css", styleHandler))

	err := http.ListenAndServe("localhost:8888", nil)
	if err != nil {
		log.Fatalln(err)
	}
}

func getPage(w http.ResponseWriter, r *http.Request) {
	_, err := strconv.Atoi(strings.TrimPrefix(r.URL.RawQuery, "page="))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

}

func showBlog(w http.ResponseWriter, r *http.Request) {
	connStr := "user=uliakulikova dbname=uliakulikova sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM blog ORDER BY id LIMIT 2 OFFSET 2")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var messages []Message

	for rows.Next() {
		p := Message{}
		err := rows.Scan(&p.Id, &p.Message)
		if err != nil {
			fmt.Println(err)
			continue
		}
		messages = append(messages, p)
	}

	data := struct {
		Msg []Message
	}{Msg: messages}

	tmpl, _ := template.ParseFiles("templates/index.html")
	err = tmpl.Execute(w, data)
	if err != nil {
		log.Fatalln(err)
	}
}

func postArticles(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "html/admin.html")
}

func publishArticle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, r.FormValue("message"))

	connStr := "user=uliakulikova dbname=uliakulikova sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	_, err = db.Exec("insert into blog (message) values ('" + r.FormValue("message") + "')")
	if err != nil {
		log.Fatalln(err)
	}
}
