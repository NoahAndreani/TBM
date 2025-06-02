package handlers

import (
	"net/http"
)

// Home gère l'affichage de la page d'accueil
func Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, "templates/index.html")
}
