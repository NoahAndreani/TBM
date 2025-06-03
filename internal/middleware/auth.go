package middleware

import (
	"context"
	"net/http"
	"strings"
	"tbcvclub/configs"
	"tbcvclub/internal/database"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// AuthMiddleware vérifie le token JWT et ajoute les informations de l'utilisateur au contexte
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var tokenString string

		// Vérifier d'abord le header Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			bearerToken := strings.Split(authHeader, " ")
			if len(bearerToken) == 2 {
				tokenString = bearerToken[1]
			}
		}

		// Si pas de token dans le header, vérifier le cookie
		if tokenString == "" {
			cookie, err := r.Cookie("auth_token")
			if err == nil {
				tokenString = cookie.Value
			}
		}

		// Si toujours pas de token, rediriger vers la page de connexion pour les pages HTML
		// ou renvoyer une erreur 401 pour les requêtes API
		if tokenString == "" {
			if strings.HasPrefix(r.URL.Path, "/api/") {
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
			} else {
				http.Redirect(w, r, "/", http.StatusSeeOther)
			}
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(configs.AppConfig.JWT.Secret), nil
		})

		if err != nil || !token.Valid {
			if strings.HasPrefix(r.URL.Path, "/api/") {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
			} else {
				http.Redirect(w, r, "/", http.StatusSeeOther)
			}
			return
		}

		// Ajoute les informations de l'utilisateur au contexte
		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "username", claims.Username)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserID récupère l'ID de l'utilisateur depuis le contexte
func GetUserID(r *http.Request) (int64, bool) {
	userID, ok := r.Context().Value("user_id").(int64)
	return userID, ok
}

// GetUsername récupère le nom d'utilisateur depuis le contexte
func GetUsername(r *http.Request) (string, bool) {
	username, ok := r.Context().Value("username").(string)
	return username, ok
}

// IsAdmin vérifie si l'utilisateur est un administrateur
func IsAdmin(r *http.Request) bool {
	userID, ok := GetUserID(r)
	if !ok {
		return false
	}

	user, err := database.GetUserByID(userID)
	if err != nil || user == nil {
		return false
	}

	return user.Role == "admin"
}
