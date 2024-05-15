package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

// Spot représente la structure du fichier JSON
type Spot struct {
	Records []SpotRecord `json:"records"`
	Offset  string       `json:"offset"`
}

// SpotRecord représente un enregistrement dans le fichier JSON
type SpotRecord struct {
	ID          string     `json:"id"`
	Fields      SpotFields `json:"fields"`
	CreatedTime string     `json:"createdTime"`
}

// SpotFields représente les champs de l'enregistrement Spot
type SpotFields struct {
	SurfBreak               []string    `json:"Surf Break"`
	DifficultyLevel         int         `json:"Difficulty Level"`
	Destination             string      `json:"Destination"`
	Geocode                 string      `json:"Geocode"`
	Influencers             []string    `json:"Influencers"`
	MagicSeaweedLink        string      `json:"Magic Seaweed Link"`
	Photos                  []SpotPhoto `json:"Photos"`
	PeakSurfSeasonBegins    string      `json:"Peak Surf Season Begins"`
	DestinationStateCountry string      `json:"Destination State/Country"`
	PeakSurfSeasonEnds      string      `json:"Peak Surf Season Ends"`
	Address                 string      `json:"Address"`
}

// SpotPhoto représente une photo dans les champs Photos de l'enregistrement SpotFields
type SpotPhoto struct {
	ID         string          `json:"id"`
	URL        string          `json:"url"`
	Filename   string          `json:"filename"`
	Size       int             `json:"size"`
	Type       string          `json:"type"`
	Thumbnails SpotPhotoThumbs `json:"thumbnails"`
}

// SpotPhotoThumbs représente les miniatures d'une photo
type SpotPhotoThumbs struct {
	Small SpotThumbnail `json:"small"`
	Large SpotThumbnail `json:"large"`
	Full  SpotThumbnail `json:"full"`
}

// SpotThumbnail représente une miniature d'une photo
type SpotThumbnail struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

func getSpots(w http.ResponseWriter, r *http.Request) {
	// Lecture du fichier JSON
	jsonData, err := os.ReadFile("spot.json")
	if err != nil {
		http.Error(w, "Failed to read JSON file", http.StatusInternalServerError)
		return
	}

	// Conversion du contenu JSON en une structure Spot
	var spots Spot
	err = json.Unmarshal(jsonData, &spots)
	if err != nil {
		http.Error(w, "Failed to parse JSON", http.StatusInternalServerError)
		return
	}

	// Encodage de la structure Spot en JSON et envoi en réponse
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(spots)
}

func main() {
	http.HandleFunc("/api/spots", getSpots)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
