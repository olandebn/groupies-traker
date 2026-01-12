package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

// On charge les templates normalement
var tmpls = template.Must(template.ParseGlob("templates/*.html"))

func main() {
	// On sert le dossier "static" pour que le CSS fonctionne
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Routes
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
	// Appel à ton data.go (FetchData)
	err := FetchData("/artists", &artists)
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des artistes", http.StatusInternalServerError)
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
	FetchData("/artists/"+id, &artist)

	var rel Relation
	FetchData("/relation/"+id, &rel)

	data := struct {
		Artist   Artist
		Relation Relation
	}{
		Artist:   artist,
		Relation: rel,
	}

	tmpls.ExecuteTemplate(w, "details.html", data)
}
