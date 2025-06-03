package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"tbcvclub/internal/database"
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

	// Récupération de l'utilisateur depuis la base de données
	user, err := database.GetUserByID(userID)
	if err != nil {
		log.Printf("Erreur lors de la récupération de l'utilisateur: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
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

	// Récupération de l'utilisateur actuel
	user, err := database.GetUserByID(userID)
	if err != nil {
		log.Printf("Erreur lors de la récupération de l'utilisateur: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Mise à jour des champs si fournis
	if req.Username != "" {
		// Vérifier si le nom d'utilisateur est déjà pris
		existingUser, err := database.GetUserByUsername(req.Username)
		if err != nil {
			log.Printf("Erreur lors de la vérification du nom d'utilisateur: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if existingUser != nil && existingUser.ID != userID {
			http.Error(w, "Username already taken", http.StatusConflict)
			return
		}
		user.Username = req.Username
	}

	if req.Email != "" {
		// Vérifier si l'email est déjà utilisé
		existingUser, err := database.GetUserByEmail(req.Email)
		if err != nil {
			log.Printf("Erreur lors de la vérification de l'email: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if existingUser != nil && existingUser.ID != userID {
			http.Error(w, "Email already in use", http.StatusConflict)
			return
		}
		user.Email = req.Email
	}

	// Mise à jour dans la base de données
	if err := database.UpdateUser(user); err != nil {
		log.Printf("Erreur lors de la mise à jour de l'utilisateur: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
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

	// Récupération de l'utilisateur depuis la base de données
	user, err := database.GetUserByID(userID)
	if err != nil {
		log.Printf("Erreur lors de la récupération de l'utilisateur: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	stats := models.UserStats{
		Level:            user.Level,
		Experience:       user.Experience,
		ExperienceToNext: user.ExperienceForNextLevel(),
		TotalDistance:    user.TotalDistance,
		TotalRideTime:    user.TotalRideTime,
		ConsecutiveDays:  user.ConsecutiveDays,
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

	// Récupération de l'utilisateur depuis la base de données
	user, err := database.GetUserByID(userID)
	if err != nil {
		log.Printf("Erreur lors de la récupération de l'utilisateur: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Ajout des statistiques de la course
	user.AddRideExperience(stats.Distance, stats.Time)

	// Sauvegarde des modifications dans la base de données
	if err := database.UpdateUser(user); err != nil {
		log.Printf("Erreur lors de la mise à jour des statistiques: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
