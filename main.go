package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

var tmpls = template.Must(template.ParseGlob("templates/*.html"))

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/artist", artistHandler)

	fmt.Println("Serveur lancé sur : http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	var artists []Artist
	err := FetchData("/artists", &artists)
	if err != nil {
		http.Error(w, "Erreur chargement artistes", http.StatusInternalServerError)
		return
	}

	tmpls.ExecuteTemplate(w, "index.html", artists)
}

func artistHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	var artist Artist
	var rel Relation

	if FetchData("/artists/"+id, &artist) != nil ||
		FetchData("/relation/"+id, &rel) != nil {
		http.NotFound(w, r)
		return
	}

	data := struct {
		Artist   Artist
		Relation Relation
	}{artist, rel}

	tmpls.ExecuteTemplate(w, "details.html", data)
}
