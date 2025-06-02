package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"tbcvclub/internal/models"

	"github.com/gorilla/mux"
)

// GetStations récupère la liste de toutes les stations
func GetStations(w http.ResponseWriter, r *http.Request) {
	// Appel à l'API externe pour récupérer les stations
	resp, err := http.Get("http://10.33.70.223:3000/api/stations")
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des stations", http.StatusInternalServerError)
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
