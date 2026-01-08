package main

import (
	"encoding/json"
	"net/http"
)

// Structures pour stocker les données de l'API [cite: 5, 8]
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

const apiBase = "https://groupietrackers.herokuapp.com/api"

// FetchData est une fonction générique pour les appels client-serveur
func FetchData(endpoint string, target interface{}) error {
	resp, err := http.Get(apiBase + endpoint)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(target)
}
