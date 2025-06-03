package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"tbcvclub/configs"
	"tbcvclub/internal/database"
	"tbcvclub/internal/middleware"
	"tbcvclub/internal/models"
	"tbcvclub/internal/utils"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

type ErrorResponse struct {
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors,omitempty"`
}

// Register gère l'inscription d'un nouvel utilisateur
func Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Erreur décodage requête: %v", err)
		sendError(w, "Format de requête invalide", nil, http.StatusBadRequest)
		return
	}

	// Log des données reçues (sans le mot de passe)
	log.Printf("Tentative d'inscription pour l'utilisateur: %s, email: %s", req.Username, req.Email)

	// Validation des champs
	if validationErrors := utils.GetValidationErrors(req.Username, req.Email, req.Password); len(validationErrors) > 0 {
		log.Printf("Erreurs de validation pour l'utilisateur %s: %v", req.Username, validationErrors)
		sendError(w, "Erreurs de validation", validationErrors, http.StatusBadRequest)
		return
	}

	// Vérifie si l'utilisateur existe déjà
	existingUser, err := database.GetUserByUsername(req.Username)
	if err != nil {
		log.Printf("Erreur vérification utilisateur existant: %v", err)
		sendError(w, "Erreur serveur", nil, http.StatusInternalServerError)
		return
	}
	if existingUser != nil {
		log.Printf("Tentative d'inscription avec un nom d'utilisateur déjà utilisé: %s", req.Username)
		sendError(w, "Ce nom d'utilisateur est déjà utilisé", nil, http.StatusConflict)
		return
	}

	// Vérifie si l'email existe déjà
	existingEmail, err := database.GetUserByEmail(req.Email)
	if err != nil {
		log.Printf("Erreur vérification email existant: %v", err)
		sendError(w, "Erreur serveur", nil, http.StatusInternalServerError)
		return
	}
	if existingEmail != nil {
		log.Printf("Tentative d'inscription avec un email déjà utilisé: %s", req.Email)
		sendError(w, "Cette adresse email est déjà utilisée", nil, http.StatusConflict)
		return
	}

	// Hash du mot de passe
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Erreur hashage mot de passe: %v", err)
		sendError(w, "Erreur lors du traitement du mot de passe", nil, http.StatusInternalServerError)
		return
	}

	// Création du nouvel utilisateur
	now := time.Now()
	user := &models.User{
		Username:       req.Username,
		Email:          req.Email,
		HashedPassword: string(hashedPassword),
		Role:           "user",
		Level:          1,
		Experience:     0,
		CreatedAt:      now,
		LastConnection: now,
	}

	// Sauvegarde dans la base de données
	if err := database.CreateUser(user); err != nil {
		log.Printf("Erreur création utilisateur: %v", err)
		sendError(w, "Erreur lors de la création du compte", nil, http.StatusInternalServerError)
		return
	}

	// Génération du token JWT
	token, err := generateToken(user)
	if err != nil {
		log.Printf("Erreur génération token: %v", err)
		sendError(w, "Erreur lors de la génération du token", nil, http.StatusInternalServerError)
		return
	}

	// Préparation de la réponse
	response := AuthResponse{
		Token: token,
		User:  *user,
	}

	// Log de la réponse (sans données sensibles)
	log.Printf("Inscription réussie pour l'utilisateur: %s (ID: %d)", user.Username, user.ID)

	// Envoi de la réponse
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Erreur encodage réponse: %v, response: %+v", err, response)
		sendError(w, "Erreur lors de la génération de la réponse", nil, http.StatusInternalServerError)
		return
	}
}

// LoginHandler gère la connexion des utilisateurs
func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := database.GetUserByUsername(credentials.Username)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(credentials.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Créer le token JWT
	claims := &middleware.Claims{
		UserID:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(configs.AppConfig.JWT.Secret))
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	// Stocker le token dans un cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    tokenString,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Mettre à true en production avec HTTPS
		SameSite: http.SameSiteLaxMode,
		MaxAge:   86400, // 24 heures
	})

	// Renvoyer également le token dans la réponse JSON pour les clients API
	response := map[string]interface{}{
		"token": tokenString,
		"user": map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
			"role":     user.Role,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// generateToken génère un token JWT pour l'utilisateur
func generateToken(user *models.User) (string, error) {
	expirationTime := time.Now().Add(time.Duration(configs.AppConfig.JWT.ExpirationHours) * time.Hour)

	claims := &middleware.Claims{
		UserID:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(configs.AppConfig.JWT.Secret))
}

// Helpers pour envoyer les réponses
func sendJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Erreur lors de l'encodage JSON: %v, data: %+v", err, data)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			Message: "Erreur lors de la génération de la réponse",
		})
		return
	}
}

func sendError(w http.ResponseWriter, message string, errors map[string]string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(ErrorResponse{
		Message: message,
		Errors:  errors,
	}); err != nil {
		log.Printf("Erreur lors de l'encodage JSON de l'erreur: %v", err)
	}
}
