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
	ElectricBikes  int     `json:"electricBikes"`
	ClassicBikes   int     `json:"classicBikes"`
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
			ElectricBikes:  8,
			ClassicBikes:   7,
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
			ElectricBikes:  5,
			ClassicBikes:   3,
			Status:         "operational",
		},
		{
			ID:             3,
			Name:           "Station Gambetta",
			Address:        "Place Gambetta",
			Latitude:       44.8389,
			Longitude:      -0.5775,
			TotalSlots:     25,
			AvailableBikes: 12,
			ElectricBikes:  6,
			ClassicBikes:   6,
			Status:         "operational",
		},
		{
			ID:             4,
			Name:           "Station Victoire",
			Address:        "Place de la Victoire",
			Latitude:       44.8300,
			Longitude:      -0.5800,
			TotalSlots:     15,
			AvailableBikes: 10,
			ElectricBikes:  5,
			ClassicBikes:   5,
			Status:         "operational",
		},
		{
			ID:             5,
			Name:           "Station Saint-Michel",
			Address:        "Place Saint-Michel",
			Latitude:       44.8350,
			Longitude:      -0.5750,
			TotalSlots:     18,
			AvailableBikes: 7,
			ElectricBikes:  4,
			ClassicBikes:   3,
			Status:         "operational",
		},
		{
			ID:             6,
			Name:           "Station Mériadeck",
			Address:        "Place Mériadeck",
			Latitude:       44.8400,
			Longitude:      -0.5850,
			TotalSlots:     22,
			AvailableBikes: 14,
			ElectricBikes:  8,
			ClassicBikes:   6,
			Status:         "operational",
		},
		{
			ID:             7,
			Name:           "Station Grand Théâtre",
			Address:        "Place de la Comédie",
			Latitude:       44.8420,
			Longitude:      -0.5720,
			TotalSlots:     20,
			AvailableBikes: 9,
			ElectricBikes:  5,
			ClassicBikes:   4,
			Status:         "operational",
		},
		{
			ID:             8,
			Name:           "Station Jardin Public",
			Address:        "Cours de Verdun",
			Latitude:       44.8450,
			Longitude:      -0.5800,
			TotalSlots:     16,
			AvailableBikes: 11,
			ElectricBikes:  6,
			ClassicBikes:   5,
			Status:         "operational",
		},
		{
			ID:             9,
			Name:           "Station Hôtel de Ville",
			Address:        "Place de l'Hôtel de Ville",
			Latitude:       44.8370,
			Longitude:      -0.5820,
			TotalSlots:     24,
			AvailableBikes: 13,
			ElectricBikes:  7,
			ClassicBikes:   6,
			Status:         "operational",
		},
		{
			ID:             10,
			Name:           "Station Porte de Bourgogne",
			Address:        "Place de la Porte de Bourgogne",
			Latitude:       44.8350,
			Longitude:      -0.5650,
			TotalSlots:     19,
			AvailableBikes: 8,
			ElectricBikes:  4,
			ClassicBikes:   4,
			Status:         "operational",
		},
		{
			ID:             11,
			Name:           "Station Saint-Pierre",
			Address:        "Place Saint-Pierre",
			Latitude:       44.8382,
			Longitude:      -0.5735,
			TotalSlots:     18,
			AvailableBikes: 12,
			ElectricBikes:  6,
			ClassicBikes:   6,
			Status:         "operational",
		},
		{
			ID:             12,
			Name:           "Station Saint-Seurin",
			Address:        "Place des Martyrs de la Résistance",
			Latitude:       44.8415,
			Longitude:      -0.5862,
			TotalSlots:     15,
			AvailableBikes: 9,
			ElectricBikes:  4,
			ClassicBikes:   5,
			Status:         "operational",
		},
		{
			ID:             13,
			Name:           "Station Capucins",
			Address:        "Place des Capucins",
			Latitude:       44.8325,
			Longitude:      -0.5708,
			TotalSlots:     20,
			AvailableBikes: 13,
			ElectricBikes:  7,
			ClassicBikes:   6,
			Status:         "operational",
		},
		{
			ID:             14,
			Name:           "Station Chartrons",
			Address:        "Rue Notre-Dame",
			Latitude:       44.8530,
			Longitude:      -0.5690,
			TotalSlots:     17,
			AvailableBikes: 10,
			ElectricBikes:  5,
			ClassicBikes:   5,
			Status:         "operational",
		},
		{
			ID:             15,
			Name:           "Station Bastide",
			Address:        "Place Stalingrad",
			Latitude:       44.8410,
			Longitude:      -0.5530,
			TotalSlots:     16,
			AvailableBikes: 11,
			ElectricBikes:  6,
			ClassicBikes:   5,
			Status:         "operational",
		},
		{
			ID:             16,
			Name:           "Station Nansouty",
			Address:        "Place Nansouty",
			Latitude:       44.8255,
			Longitude:      -0.5730,
			TotalSlots:     14,
			AvailableBikes: 8,
			ElectricBikes:  4,
			ClassicBikes:   4,
			Status:         "operational",
		},
		{
			ID:             17,
			Name:           "Station Saint-Genès",
			Address:        "Rue Saint-Genès",
			Latitude:       44.8268,
			Longitude:      -0.5825,
			TotalSlots:     13,
			AvailableBikes: 7,
			ElectricBikes:  3,
			ClassicBikes:   4,
			Status:         "operational",
		},
		{
			ID:             18,
			Name:           "Station Victoire Sud",
			Address:        "Cours de la Marne",
			Latitude:       44.8280,
			Longitude:      -0.5630,
			TotalSlots:     15,
			AvailableBikes: 9,
			ElectricBikes:  5,
			ClassicBikes:   4,
			Status:         "operational",
		},
		{
			ID:             19,
			Name:           "Station Sainte-Croix",
			Address:        "Place Pierre Renaudel",
			Latitude:       44.8295,
			Longitude:      -0.5555,
			TotalSlots:     12,
			AvailableBikes: 7,
			ElectricBikes:  3,
			ClassicBikes:   4,
			Status:         "operational",
		},
		{
			ID:             20,
			Name:           "Station Saint-Michel Sud",
			Address:        "Rue Camille Sauvageau",
			Latitude:       44.8320,
			Longitude:      -0.5635,
			TotalSlots:     13,
			AvailableBikes: 8,
			ElectricBikes:  4,
			ClassicBikes:   4,
			Status:         "operational",
		},
		{
			ID:             21,
			Name:           "Station Fondaudège",
			Address:        "Rue Fondaudège",
			Latitude:       44.8462,
			Longitude:      -0.5790,
			TotalSlots:     15,
			AvailableBikes: 10,
			ElectricBikes:  5,
			ClassicBikes:   5,
			Status:         "operational",
		},
		{
			ID:             22,
			Name:           "Station Saint-Augustin",
			Address:        "Place de l'Église Saint-Augustin",
			Latitude:       44.8418,
			Longitude:      -0.6100,
			TotalSlots:     14,
			AvailableBikes: 8,
			ElectricBikes:  4,
			ClassicBikes:   4,
			Status:         "operational",
		},
		{
			ID:             23,
			Name:           "Station Croix Blanche",
			Address:        "Place Croix Blanche",
			Latitude:       44.8425,
			Longitude:      -0.5950,
			TotalSlots:     13,
			AvailableBikes: 7,
			ElectricBikes:  3,
			ClassicBikes:   4,
			Status:         "operational",
		},
		{
			ID:             24,
			Name:           "Station Barrière Judaïque",
			Address:        "Boulevard du Président Wilson",
			Latitude:       44.8410,
			Longitude:      -0.6000,
			TotalSlots:     15,
			AvailableBikes: 9,
			ElectricBikes:  4,
			ClassicBikes:   5,
			Status:         "operational",
		},
		{
			ID:             25,
			Name:           "Station Place de la Bourse",
			Address:        "Place de la Bourse",
			Latitude:       44.8412,
			Longitude:      -0.5690,
			TotalSlots:     18,
			AvailableBikes: 12,
			ElectricBikes:  6,
			ClassicBikes:   6,
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
