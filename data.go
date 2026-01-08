package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Structures mises à jour
type Artist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
}

type Relation struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

// Récupère la liste de tous les artistes [cite: 5]
func FetchArtists() ([]Artist, error) {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var artists []Artist
	err = json.NewDecoder(resp.Body).Decode(&artists)
	return artists, err
}

// Récupère les relations (dates + lieux) d'un artiste spécifique
func FetchRelation(id string) (Relation, error) {
	url := fmt.Sprintf("https://groupietrackers.herokuapp.com/api/relation/%s", id)
	resp, err := http.Get(url)
	if err != nil {
		return Relation{}, err
	}
	defer resp.Body.Close()
	var rel Relation
	err = json.NewDecoder(resp.Body).Decode(&rel)
	return rel, err
}
