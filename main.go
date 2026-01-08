package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func main() {
	tmplIndex := template.Must(template.ParseFiles("templates/index.html"))
	tmplArtist := template.Must(template.ParseFiles("templates/artist.html"))

	// Page d'accueil : Liste des artistes
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		artists, _ := FetchArtists()
		tmplIndex.Execute(w, artists)
	})

	// Page de détails : Relation Dates/Lieux (Action client-serveur) [cite: 12]
	http.HandleFunc("/artist", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		// On récupère les deux types de données pour cette page
		artists, _ := FetchArtists() // Pour le nom/image
		var currentArtist Artist
		for _, a := range artists {
			if fmt.Sprintf("%d", a.ID) == id {
				currentArtist = a
				break
			}
		}

		relation, _ := FetchRelation(id)

		data := struct {
			Artist   Artist
			Relation Relation
		}{currentArtist, relation}

		tmplArtist.Execute(w, data)
	})

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
