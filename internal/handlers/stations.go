package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"tbcvclub/internal/models"

	"github.com/gorilla/mux"
)

type Station struct {
	ID             int64   `json:"id"`
	Name           string  `json:"name"`
	Address        string  `json:"address"`
	Latitude       float64 `json:"latitude"`
	Longitude      float64 `json:"longitude"`
	TotalSlots     int     `json:"totalSlots"`
	AvailableBikes int     `json:"availableBikes"`
	Status         string  `json:"status"`
}

// GetStations renvoie la liste des stations
func GetStations(w http.ResponseWriter, r *http.Request) {
	// Pour l'instant, on renvoie des données de test
	stations := []Station{
		{
			ID:             1,
			Name:           "Station Pey Berland",
			Address:        "Place Pey Berland",
			Latitude:       44.837789,
			Longitude:      -0.57918,
			TotalSlots:     20,
			AvailableBikes: 15,
			Status:         "operational",
		},
		{
			ID:             2,
			Name:           "Station Quinconces",
			Address:        "Place des Quinconces",
			Latitude:       44.843849,
			Longitude:      -0.574502,
			TotalSlots:     30,
			AvailableBikes: 8,
			Status:         "operational",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stations)
}

// GetStationByID récupère les détails d'une station spécifique
func GetStationByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	_, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID de station invalide", http.StatusBadRequest)
		return
	}

	// Appel à l'API externe pour récupérer les détails de la station
	resp, err := http.Get("http://10.33.70.223:3000/api/stations/" + vars["id"])
	if err != nil {
		http.Error(w, "Erreur lors de la récupération de la station", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var station models.Station
	if err := json.NewDecoder(resp.Body).Decode(&station); err != nil {
		http.Error(w, "Erreur lors du décodage des données", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(station)
}

// GetNearbyStations récupère les stations proches d'une position donnée
func GetNearbyStations(w http.ResponseWriter, r *http.Request) {
	lat := r.URL.Query().Get("lat")
	lon := r.URL.Query().Get("lon")
	radius := r.URL.Query().Get("radius")

	if lat == "" || lon == "" {
		http.Error(w, "Latitude et longitude requises", http.StatusBadRequest)
		return
	}

	if radius == "" {
		radius = "1000" // Rayon par défaut en mètres
	}

	// Appel à l'API externe pour récupérer les stations proches
	url := "http://10.33.70.223:3000/api/stations/nearby?lat=" + lat + "&lon=" + lon + "&radius=" + radius
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des stations proches", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var stations []models.Station
	if err := json.NewDecoder(resp.Body).Decode(&stations); err != nil {
		http.Error(w, "Erreur lors du décodage des données", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stations)
}

// GetStationStats récupère les statistiques d'une station
func GetStationStats(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	_, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID de station invalide", http.StatusBadRequest)
		return
	}

	// Appel à l'API externe pour récupérer les statistiques de la station
	resp, err := http.Get("http://10.33.70.223:3000/api/stations/" + vars["id"] + "/stats")
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des statistiques", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var stats models.StationStats
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		http.Error(w, "Erreur lors du décodage des données", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
