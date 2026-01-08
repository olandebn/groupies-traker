package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

// Chargement des templates globalement pour éviter de les re-parser à chaque requête
var tmpls = template.Must(template.ParseGlob("templates/*.html"))

func main() {
	// Configuration du serveur de fichiers statiques (pour le CSS premium) [cite: 9]
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Route principale : Affichage de tous les artistes [cite: 3, 5]
	http.HandleFunc("/", homeHandler)

	// Route Détails : Action client-serveur pour les relations [cite: 8, 11, 12]
	http.HandleFunc("/artist", artistHandler)

	fmt.Println("Siuuu! Serveur lancé sur : http://localhost:8080")
	// Gestion des erreurs fatales pour que le serveur ne crash pas [cite: 18]
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Handler pour la page d'accueil
func homeHandler(w http.ResponseWriter, r *http.Request) {
	// Gestion de l'erreur 404 si l'URL est incorrecte [cite: 19]
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	var artists []Artist
	// Appel à l'API via data.go [cite: 4, 14]
	err := FetchData("/artists", &artists)
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des artistes", http.StatusInternalServerError)
		return
	}

	// Visualisation des données via le template index.html [cite: 9, 31]
	tmpls.ExecuteTemplate(w, "index.html", artists)
}

// Handler pour les détails (Événement client-serveur)
func artistHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// 1. Récupérer les infos de l'artiste [cite: 5]
	var artist Artist
	errArt := FetchData("/artists/"+id, &artist)

	// 2. Récupérer les relations (Dates et Lieux) [cite: 8, 14]
	var rel Relation
	errRel := FetchData("/relation/"+id, &rel)

	if errArt != nil || errRel != nil {
		http.Error(w, "Données introuvables", http.StatusNotFound)
		return
	}

	// On combine les données pour le template [cite: 27, 28]
	data := struct {
		Artist   Artist
		Relation Relation
	}{
		Artist:   artist,
		Relation: rel,
	}

	tmpls.ExecuteTemplate(w, "details.html", data)
}
