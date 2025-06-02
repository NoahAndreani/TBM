package middleware

import (
	"context"
	"net/http"
	"strings"
	"tbcvclub/configs"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// AuthMiddleware vérifie le token JWT et ajoute les informations de l'utilisateur au contexte
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}

		tokenStr := bearerToken[1]
		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(configs.AppConfig.JWT.Secret), nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				http.Error(w, "Invalid token signature", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Ajoute les informations de l'utilisateur au contexte
		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "username", claims.Username)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserID récupère l'ID de l'utilisateur depuis le contexte
func GetUserID(r *http.Request) (int, bool) {
	userID, ok := r.Context().Value("user_id").(int)
	return userID, ok
}

// GetUsername récupère le nom d'utilisateur depuis le contexte
func GetUsername(r *http.Request) (string, bool) {
	username, ok := r.Context().Value("username").(string)
	return username, ok
}
