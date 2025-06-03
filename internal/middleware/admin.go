package middleware

import (
	"net/http"
)

// AdminMiddleware vérifie si l'utilisateur est un administrateur
func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Vérifier si l'utilisateur est authentifié et est un administrateur
		if !IsAdmin(r) {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}
