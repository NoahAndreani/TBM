package handlers

import (
	"encoding/json"
	"net/http"
	"time"
)

type NewsItem struct {
	ID      int64     `json:"id"`
	Title   string    `json:"title"`
	Content string    `json:"content"`
	Date    time.Time `json:"date"`
	Source  string    `json:"source,omitempty"`
	Link    string    `json:"link,omitempty"`
}

// GetNews renvoie la liste des actualités
func GetNews(w http.ResponseWriter, r *http.Request) {
	// Pour l'instant, on renvoie des données de test
	news := []NewsItem{
		{
			ID:      1,
			Title:   "Nouveau service de location de vélos",
			Content: "TBC Vclub lance son service de location de vélos à Bordeaux. Profitez de nos offres de lancement !",
			Date:    time.Now().Add(-24 * time.Hour),
			Source:  "TBC Vclub",
		},
		{
			ID:      2,
			Title:   "Maintenance des stations",
			Content: "Une maintenance sera effectuée sur les stations du centre-ville ce week-end.",
			Date:    time.Now().Add(-2 * 24 * time.Hour),
			Source:  "TBC Vclub",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(news)
}
