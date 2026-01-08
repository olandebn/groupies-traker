package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

var tmpl = template.Must(template.ParseGlob("templates/*.html"))

func main() {
	// Chargement des fichiers CSS (Dossier static)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Route Accueil : Affiche tous les artistes [cite: 5, 9]
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		var artists []Artist
		if err := FetchData("/artists", &artists); err != nil {
			http.Error(w, "Erreur lors de la récupération des artistes", 500)
			return
		}
		tmpl.ExecuteTemplate(w, "index.html", artists)
	})

	// Route Détails : Affiche les relations (Lieux et Dates) [cite: 8, 11]
	http.HandleFunc("/artist", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Redirect(w, r, "/", 302)
			return
		}

		var rel Relation
		if err := FetchData("/relation/"+id, &rel); err != nil {
			http.Error(w, "Relations introuvables", 404)
			return
		}

		// Récupération de l'artiste pour le nom et l'image
		var artists []Artist
		FetchData("/artists", &artists)
		var currentArtist Artist
		for _, a := range artists {
			if fmt.Sprintf("%d", a.ID) == id {
				currentArtist = a
				break
			}
		}

		data := struct {
			Artist   Artist
			Relation Relation
		}{currentArtist, rel}

		tmpl.ExecuteTemplate(w, "details.html", data)
	})

	log.Println("Serveur en ligne : http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
