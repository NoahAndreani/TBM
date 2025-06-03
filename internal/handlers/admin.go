package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"tbcvclub/internal/database"
	"tbcvclub/internal/middleware"

	"github.com/gorilla/mux"
)

// AdminPage affiche la page d'administration
func AdminPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/admin.html")
}

// GetAllUsers retourne la liste de tous les utilisateurs
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	// Vérifier si l'utilisateur est admin
	if !middleware.IsAdmin(r) {
		http.Error(w, "Non autorisé", http.StatusUnauthorized)
		return
	}

	users, err := database.GetAllUsers()
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des utilisateurs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// DeleteUser supprime un utilisateur
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	// Vérifier si l'utilisateur est admin
	if !middleware.IsAdmin(r) {
		http.Error(w, "Non autorisé", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	userID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "ID utilisateur invalide", http.StatusBadRequest)
		return
	}

	// Vérifier que l'utilisateur n'est pas admin
	user, err := database.GetUserByID(userID)
	if err != nil {
		http.Error(w, "Erreur lors de la récupération de l'utilisateur", http.StatusInternalServerError)
		return
	}
	if user == nil {
		http.Error(w, "Utilisateur non trouvé", http.StatusNotFound)
		return
	}
	if user.Role == "admin" {
		http.Error(w, "Impossible de supprimer un administrateur", http.StatusForbidden)
		return
	}

	// Supprimer l'utilisateur
	if err := database.DeleteUser(userID); err != nil {
		http.Error(w, "Erreur lors de la suppression de l'utilisateur", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
