package handlers

import (
	"encoding/json"
	"net/http"
	"tbcvclub/configs"
	"tbcvclub/internal/middleware"
	"tbcvclub/internal/models"
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
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

// Register gère l'inscription d'un nouvel utilisateur
func Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Vérification des champs requis
	if req.Username == "" || req.Email == "" || req.Password == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	// Hash du mot de passe
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error processing password", http.StatusInternalServerError)
		return
	}

	// Création du nouvel utilisateur
	user := models.User{
		Username:       req.Username,
		Email:          req.Email,
		HashedPassword: string(hashedPassword),
		Level:          1,
		Experience:     0,
		CreatedAt:      time.Now(),
		LastConnection: time.Now(),
	}

	// TODO: Sauvegarder l'utilisateur dans la base de données

	// Génération du token JWT
	token, err := generateToken(user)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	response := AuthResponse{
		Token: token,
		User:  user,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Login gère la connexion d'un utilisateur
func Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// TODO: Récupérer l'utilisateur depuis la base de données
	// Pour l'exemple, on simule un utilisateur
	user := models.User{
		ID:             1,
		Username:       "test",
		Email:          req.Email,
		HashedPassword: "$2a$10$...", // À remplacer par le vrai hash
	}

	// Vérification du mot de passe
	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(req.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Mise à jour de la dernière connexion
	user.LastConnection = time.Now()
	// TODO: Sauvegarder la mise à jour dans la base de données

	// Ajout de l'expérience pour la connexion journalière
	user.AddDailyLoginExperience()

	// Génération du token JWT
	token, err := generateToken(user)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	response := AuthResponse{
		Token: token,
		User:  user,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// generateToken génère un token JWT pour l'utilisateur
func generateToken(user models.User) (string, error) {
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
