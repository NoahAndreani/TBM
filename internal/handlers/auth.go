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
		sendError(w, "Format de requête invalide", nil, http.StatusBadRequest)
		return
	}

	// Validation des champs
	if validationErrors := utils.GetValidationErrors(req.Username, req.Email, req.Password); len(validationErrors) > 0 {
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

	sendJSON(w, AuthResponse{
		Token: token,
		User:  *user,
	})
}

// Login gère la connexion d'un utilisateur
func Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "Format de requête invalide", nil, http.StatusBadRequest)
		return
	}

	// Validation du nom d'utilisateur
	if !utils.ValidateUsername(req.Username) {
		sendError(w, "Format de nom d'utilisateur invalide", nil, http.StatusBadRequest)
		return
	}

	log.Printf("Tentative de connexion pour l'utilisateur: %s", req.Username)

	// Récupération de l'utilisateur depuis la base de données
	user, err := database.GetUserByUsername(req.Username)
	if err != nil {
		log.Printf("Erreur récupération utilisateur: %v", err)
		sendError(w, "Erreur serveur", nil, http.StatusInternalServerError)
		return
	}
	if user == nil {
		sendError(w, "Nom d'utilisateur ou mot de passe incorrect", nil, http.StatusUnauthorized)
		return
	}

	// Vérification du mot de passe
	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(req.Password)); err != nil {
		log.Printf("Mot de passe incorrect pour l'utilisateur: %s", req.Username)
		sendError(w, "Nom d'utilisateur ou mot de passe incorrect", nil, http.StatusUnauthorized)
		return
	}

	log.Printf("Connexion réussie pour l'utilisateur: %s", req.Username)

	// Mise à jour de la dernière connexion et ajout de l'expérience
	user.LastConnection = time.Now()
	user.AddDailyLoginExperience()

	// Sauvegarde des modifications
	if err := database.UpdateUser(user); err != nil {
		log.Printf("Erreur mise à jour utilisateur: %v", err)
		sendError(w, "Erreur lors de la mise à jour du profil", nil, http.StatusInternalServerError)
		return
	}

	// Génération du token JWT
	token, err := generateToken(user)
	if err != nil {
		log.Printf("Erreur génération token: %v", err)
		sendError(w, "Erreur lors de la génération du token", nil, http.StatusInternalServerError)
		return
	}

	sendJSON(w, AuthResponse{
		Token: token,
		User:  *user,
	})
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
	json.NewEncoder(w).Encode(data)
}

func sendError(w http.ResponseWriter, message string, errors map[string]string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{
		Message: message,
		Errors:  errors,
	})
}
