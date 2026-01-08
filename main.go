package main

import (
	"html/template"
	"log"
	"net/http"
)

func main() {
	tmpl := template.Must(template.ParseGlob("templates/*.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		var artists []Artist
		FetchData("https://groupietrackers.herokuapp.com/api/artists", &artists)
		tmpl.ExecuteTemplate(w, "index.html", artists)
	})

	http.HandleFunc("/details", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		var rel Relation
		FetchData("https://groupietrackers.herokuapp.com/api/relation/"+id, &rel)
		tmpl.ExecuteTemplate(w, "details.html", rel)
	})

	log.Println("Serveur lancé sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
