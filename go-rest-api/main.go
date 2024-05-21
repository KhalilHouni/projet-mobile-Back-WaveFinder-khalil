// API Surf Spots by Khalil Houni
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

// Définition des structures de données

type Spot struct {
	Records []SpotRecord `json:"records"`
	Offset  string       `json:"offset"`
}

type SpotRecord struct {
	ID          string     `json:"id"`
	Fields      SpotFields `json:"fields"`
	CreatedTime string     `json:"createdTime"`
}

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

type SpotPhoto struct {
	ID         string          `json:"id"`
	URL        string          `json:"url"`
	Filename   string          `json:"filename"`
	Size       int             `json:"size"`
	Type       string          `json:"type"`
	Thumbnails SpotPhotoThumbs `json:"thumbnails"`
}

type SpotPhotoThumbs struct {
	Small SpotThumbnail `json:"small"`
	Large SpotThumbnail `json:"large"`
	Full  SpotThumbnail `json:"full"`
}

type SpotThumbnail struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

// Fonctions utilitaires

func readJSONFile(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}

func parseJSONData(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// Gestion des requêtes

func getSpots(w http.ResponseWriter, r *http.Request) {
	// Lecture du fichier JSON
	jsonData, err := readJSONFile("spot.json")
	if err != nil {
		log.Println("Error reading JSON file:", err)
		http.Error(w, "Failed to read JSON file", http.StatusInternalServerError)
		return
	}

	// Conversion du contenu JSON en une structure Spot
	var spots Spot
	if err := parseJSONData(jsonData, &spots); err != nil {
		log.Println("Error parsing JSON data:", err)
		http.Error(w, "Failed to parse JSON", http.StatusInternalServerError)
		return
	}

	// Création d'une liste simplifiée pour l'affichage dans la liste
	var simplifiedList []SpotRecord
	simplifiedList = append(simplifiedList, spots.Records...)

	// Encodage de la liste simplifiée en JSON et envoi en réponse
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(simplifiedList); err != nil {
		log.Println("Error encoding JSON response:", err)
		http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
	}

}

func getOneSpot(w http.ResponseWriter, r *http.Request) {
	// Récupération de l'ID du spot depuis les paramètres de la requête
	vars := mux.Vars(r)
	spotID := vars["id"]

	// Lecture du fichier JSON
	jsonData, err := readJSONFile("spot.json")
	if err != nil {
		log.Println("Error reading JSON file:", err)
		http.Error(w, "Failed to read JSON file", http.StatusInternalServerError)
		return
	}

	// Conversion du contenu JSON en une structure Spot
	var spots Spot
	if err := parseJSONData(jsonData, &spots); err != nil {
		log.Println("Error parsing JSON data:", err)
		http.Error(w, "Failed to parse JSON", http.StatusInternalServerError)
		return
	}

	// Recherche du spot correspondant à l'ID spécifié
	var foundSpot SpotRecord
	for _, record := range spots.Records {
		if record.ID == spotID {
			foundSpot = record
			break
		}
	}

	// Vérification si le spot a été trouvé
	if foundSpot.ID == "" {
		http.Error(w, "Spot not found", http.StatusNotFound)
		return
	}

	// Encodage du spot trouvé en JSON et envoi en réponse
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(foundSpot); err != nil {
		log.Println("Error encoding JSON response:", err)
		http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
	}
}

// Fonction pour rajouter un nouveau Post
func addSpot(w http.ResponseWriter, r *http.Request) {
	var newSpot SpotRecord

	// Décodage de la requête JSON
	if err := json.NewDecoder(r.Body).Decode(&newSpot); err != nil {
		log.Println("Error decoding JSON request:", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Lecture du fichier JSON
	jsonData, err := readJSONFile("spot.json")
	if err != nil {
		log.Println("Error reading JSON file:", err)
		http.Error(w, "Failed to read JSON file", http.StatusInternalServerError)
		return
	}

	// Conversion du contenu JSON en une structure Spot
	var spots Spot
	if err := parseJSONData(jsonData, &spots); err != nil {
		log.Println("Error parsing JSON data:", err)
		http.Error(w, "Failed to parse JSON", http.StatusInternalServerError)
		return
	}

	// Ajout du nouveau spot à la liste
	spots.Records = append(spots.Records, newSpot)

	// Encodage et sauvegarde de la nouvelle liste
	updatedData, err := json.Marshal(spots)
	if err != nil {
		log.Println("Error encoding JSON data:", err)
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		return
	}

	if err := os.WriteFile("spot.json", updatedData, 0644); err != nil {
		log.Println("Error writing JSON file:", err)
		http.Error(w, "Failed to write JSON file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func updateSpot(w http.ResponseWriter, r *http.Request) {
	// Récupération de l'ID du spot depuis les paramètres de la requête
	vars := mux.Vars(r)
	spotID := vars["id"]

	// Décodage de la requête JSON pour les champs à mettre à jour
	var updatedFields struct {
		SurfBreak []string    `json:"Surf Break"`
		Photos    []SpotPhoto `json:"Photos"`
		Address   string      `json:"Address"`
	}
	if err := json.NewDecoder(r.Body).Decode(&updatedFields); err != nil {
		log.Println("Error decoding JSON request:", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Lecture du fichier JSON
	jsonData, err := readJSONFile("spot.json")
	if err != nil {
		log.Println("Error reading JSON file:", err)
		http.Error(w, "Failed to read JSON file", http.StatusInternalServerError)
		return
	}

	// Conversion du contenu JSON en une structure Spot
	var spots Spot
	if err := parseJSONData(jsonData, &spots); err != nil {
		log.Println("Error parsing JSON data:", err)
		http.Error(w, "Failed to parse JSON", http.StatusInternalServerError)
		return
	}

	// Recherche et mise à jour des champs du spot correspondant à l'ID spécifié
	spotUpdated := false
	for i, record := range spots.Records {
		if record.ID == spotID {
			if updatedFields.SurfBreak != nil {
				spots.Records[i].Fields.SurfBreak = updatedFields.SurfBreak
			}
			if updatedFields.Photos != nil {
				spots.Records[i].Fields.Photos = updatedFields.Photos
			}
			if updatedFields.Address != "" {
				spots.Records[i].Fields.Address = updatedFields.Address
			}
			spotUpdated = true
			break
		}
	}

	// Vérification si le spot a été trouvé
	if !spotUpdated {
		http.Error(w, "Spot not found", http.StatusNotFound)
		return
	}

	// Encodage et sauvegarde de la nouvelle liste
	updatedData, err := json.Marshal(spots)
	if err != nil {
		log.Println("Error encoding JSON data:", err)
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		return
	}

	if err := os.WriteFile("spot.json", updatedData, 0644); err != nil {
		log.Println("Error writing JSON file:", err)
		http.Error(w, "Failed to write JSON file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func deleteSpot(w http.ResponseWriter, r *http.Request) {
	// Récupération de l'ID du spot depuis les paramètres de la requête
	vars := mux.Vars(r)
	spotID := vars["id"]

	// Lecture du fichier JSON
	jsonData, err := readJSONFile("spot.json")
	if err != nil {
		log.Println("Error reading JSON file:", err)
		http.Error(w, "Failed to read JSON file", http.StatusInternalServerError)
		return
	}

	// Conversion du contenu JSON en une structure Spot
	var spots Spot
	if err := parseJSONData(jsonData, &spots); err != nil {
		log.Println("Error parsing JSON data:", err)
		http.Error(w, "Failed to parse JSON", http.StatusInternalServerError)
		return
	}

	// Recherche du spot correspondant à l'ID spécifié
	index := -1
	for i, record := range spots.Records {
		if record.ID == spotID {
			index = i
			break
		}
	}

	// Vérification si le spot a été trouvé
	if index == -1 {
		http.Error(w, "Spot not found", http.StatusNotFound)
		return
	}

	// Suppression du spot de la liste
	spots.Records = append(spots.Records[:index], spots.Records[index+1:]...)

	// Encodage et sauvegarde de la nouvelle liste
	updatedData, err := json.Marshal(spots)
	if err != nil {
		log.Println("Error encoding JSON data:", err)
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		return
	}

	if err := os.WriteFile("spot.json", updatedData, 0644); err != nil {
		log.Println("Error writing JSON file:", err)
		http.Error(w, "Failed to write JSON file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/api/spots", getSpots).Methods("GET")
	r.HandleFunc("/api/spots/{id}", getOneSpot).Methods("GET")
	r.HandleFunc("/api/spots", addSpot).Methods("POST")
	r.HandleFunc("/api/spots/{id}", updateSpot).Methods("PUT")
	r.HandleFunc("/api/spots/{id}", deleteSpot).Methods("DELETE")
	log.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
