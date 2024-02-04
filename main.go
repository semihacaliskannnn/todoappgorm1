package main

import (
	"html/template"
	"log"
	"net/http"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB
var err error

type Todo struct {
	ID        int
	Title     string
	Completed bool
}

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("veri tabanına bağlanılamadı.!")
	}

	db.AutoMigrate(&Todo{})

	http.HandleFunc("/", viewHandler)
	http.HandleFunc("/add", addHnadler)
	http.HandleFunc("/toggle", toggleHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	var todos []Todo
	db.Find(&todos)
	renderTemplate(w, "index", todos)
}

func addHnadler(w http.ResponseWriter, r *http.Request) {
	text := r.FormValue("text")
	db.Create(&Todo{Title: text})
	http.Redirect(w, r, "/", http.StatusFound)
}

func toggleHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	var todo Todo
	db.First(&todo, id)
	todo.Completed = !todo.Completed
	db.Save(&todo)
	http.Redirect(w, r, "/", http.StatusFound)
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	t, _ := template.ParseFiles("templates/" + tmpl + ".html")
	t.Execute(w, data)
}
