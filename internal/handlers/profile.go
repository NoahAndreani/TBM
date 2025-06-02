package handlers

import (
	"encoding/json"
	"net/http"
	"tbcvclub/internal/middleware"
	"tbcvclub/internal/models"
)

type UpdateProfileRequest struct {
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

// GetProfile récupère les informations du profil de l'utilisateur
func GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// TODO: Récupérer l'utilisateur depuis la base de données
	// Pour l'exemple, on simule un utilisateur
	user := models.User{
		ID:            userID,
		Username:      "test_user",
		Email:         "test@example.com",
		Level:         5,
		Experience:    2500,
		TotalDistance: 150.5,
		TotalRideTime: 720, // 12 heures
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// UpdateProfile met à jour les informations du profil de l'utilisateur
func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// TODO: Mettre à jour l'utilisateur dans la base de données
	// Pour l'exemple, on simule la mise à jour
	user := models.User{
		ID:       userID,
		Username: req.Username,
		Email:    req.Email,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// GetProfileStats récupère les statistiques du profil de l'utilisateur
func GetProfileStats(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// TODO: Récupérer les statistiques depuis la base de données
	// Pour l'exemple, on simule les statistiques
	user := models.User{
		ID:            userID,
		Level:         5,
		Experience:    2500,
		TotalDistance: 150.5,
		TotalRideTime: 720,
	}

	stats := models.UserStats{
		Level:            user.Level,
		Experience:       user.Experience,
		ExperienceToNext: user.ExperienceForNextLevel(),
		TotalDistance:    user.TotalDistance,
		TotalRideTime:    user.TotalRideTime,
		ConsecutiveDays:  5,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// AddRideStats ajoute les statistiques d'une course
func AddRideStats(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	type RideStats struct {
		Distance float64 `json:"distance"`
		Time     int     `json:"time"` // en minutes
	}

	var stats RideStats
	if err := json.NewDecoder(r.Body).Decode(&stats); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// TODO: Récupérer l'utilisateur depuis la base de données
	user := models.User{
		ID:            userID,
		Level:         5,
		Experience:    2500,
		TotalDistance: 150.5,
		TotalRideTime: 720,
	}

	// Ajout des statistiques de la course
	user.AddRideExperience(stats.Distance, stats.Time)

	// TODO: Sauvegarder les modifications dans la base de données

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
