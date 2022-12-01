package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"html/template"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", showBlog)
	http.HandleFunc("/admin", postArticles)
	http.HandleFunc("/publication", publishArticle)

	styleHandler := http.FileServer(http.Dir("./css"))
	http.Handle("/css/", http.StripPrefix("/css", styleHandler))

	err := http.ListenAndServe("localhost:8888", nil)
	if err != nil {
		log.Fatalln(err)
	}
}

func showBlog(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/index.html")
	err := tmpl.Execute(w, nil)
	if err != nil {
		log.Fatalln(err)
	}
}

func postArticles(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "html/admin.html")
}

func publishArticle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, r.FormValue("message"))

	connStr := "user=bemmanue dbname=bemmanue sslmode=disable"
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
