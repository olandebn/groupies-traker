package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings" // Nécessaire pour manipuler les chaînes de caractères
)

// --- BONUS : NETTOYAGE DES DONNÉES ---
// On définit ici des fonctions qu'on pourra utiliser directement dans le HTML.
// Celle-ci sert à rendre les noms de lieux plus jolis.
var funcMap = template.FuncMap{
	"cleanLocation": func(s string) string {
		// Remplace les tirets bas "_" par des espaces
		s = strings.ReplaceAll(s, "_", " ")
		// Remplace les tirets "-" par ", "
		s = strings.ReplaceAll(s, "-", ", ")
		// Met la première lettre de chaque mot en majuscule (simple implementation)
		s = strings.Title(strings.ToLower(s))
		return s
	},
}

// --- CHARGEMENT DES TEMPLATES ---
// IMPORTANT : On charge la FuncMap AVANT de parser les fichiers.
var tmpls = template.Must(template.New("").Funcs(funcMap).ParseGlob("templates/*.html"))

func main() {
	// Configuration du serveur de fichiers statiques (pour le CSS et images)
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Routes
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/artist", artistHandler)

	fmt.Println("Siuuu! Serveur lancé sur : http://localhost:8080")
	// Gestion des erreurs fatales
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Handler pour la page d'accueil
func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	var artists []Artist
	// Appel à l'API via data.go (assure-toi que ton fichier data.go est dans le même dossier)
	err := FetchData("/artists", &artists)
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des artistes", http.StatusInternalServerError)
		return
	}

	tmpls.ExecuteTemplate(w, "index.html", artists)
}

// Handler pour les détails
func artistHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// 1. Récupérer les infos de l'artiste
	var artist Artist
	errArt := FetchData("/artists/"+id, &artist)

	// 2. Récupérer les relations (Dates et Lieux)
	var rel Relation
	errRel := FetchData("/relation/"+id, &rel)

	if errArt != nil || errRel != nil {
		http.Error(w, "Données introuvables ou API indisponible", http.StatusNotFound)
		return
	}

	// On combine les données pour le template
	// Dans le HTML, on utilisera .Artist et .Relation
	data := struct {
		Artist   Artist
		Relation Relation
	}{
		Artist:   artist,
		Relation: rel,
	}

	tmpls.ExecuteTemplate(w, "details.html", data)
}
