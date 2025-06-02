package handlers

import (
	"encoding/json"
	"net/http"
	"tbcvclub/internal/middleware"
	"tbcvclub/internal/models"
	"time"
)

var subscriptionPlans = []models.SubscriptionPlan{
	{
		ID:          1,
		Type:        models.Daily,
		Name:        "Pass Journée",
		Description: "Accès illimité pendant 24h",
		Price:       5.0,
		Duration:    24 * time.Hour,
		Benefits:    []string{"Accès illimité aux vélos", "Premiers 30 minutes gratuites"},
	},
	{
		ID:          2,
		Type:        models.Weekly,
		Name:        "Pass Semaine",
		Description: "Accès illimité pendant 7 jours",
		Price:       15.0,
		Duration:    7 * 24 * time.Hour,
		Benefits:    []string{"Accès illimité aux vélos", "Premiers 30 minutes gratuites", "Support prioritaire"},
	},
	{
		ID:          3,
		Type:        models.Monthly,
		Name:        "Pass Mensuel",
		Description: "Accès illimité pendant 30 jours",
		Price:       30.0,
		Duration:    30 * 24 * time.Hour,
		Benefits:    []string{"Accès illimité aux vélos", "Premiers 45 minutes gratuites", "Support prioritaire", "Bonus XP x1.5"},
	},
	{
		ID:          4,
		Type:        models.Yearly,
		Name:        "Pass Annuel",
		Description: "Accès illimité pendant 1 an",
		Price:       200.0,
		Duration:    365 * 24 * time.Hour,
		Benefits:    []string{"Accès illimité aux vélos", "Premiers 60 minutes gratuites", "Support VIP", "Bonus XP x2", "Assurance vol incluse"},
	},
}

// GetSubscriptionPlans renvoie la liste des forfaits disponibles
func GetSubscriptionPlans(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subscriptionPlans)
}

type SubscribeRequest struct {
	PlanID int `json:"plan_id"`
}

// Subscribe permet à un utilisateur de souscrire à un forfait
func Subscribe(w http.ResponseWriter, r *http.Request) {
	// Vérification de l'authentification
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req SubscribeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Recherche du forfait
	var plan *models.SubscriptionPlan
	for _, p := range subscriptionPlans {
		if p.ID == req.PlanID {
			plan = &p
			break
		}
	}

	if plan == nil {
		http.Error(w, "Plan not found", http.StatusNotFound)
		return
	}

	// Création de l'abonnement
	subscription := models.Subscription{
		UserID:    userID,
		Type:      plan.Type,
		StartDate: time.Now(),
		EndDate:   time.Now().Add(plan.Duration),
		Price:     plan.Price,
		IsActive:  true,
	}

	// TODO: Sauvegarder l'abonnement dans la base de données

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subscription)
}

// GetActiveSubscription récupère l'abonnement actif de l'utilisateur
func GetActiveSubscription(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// TODO: Récupérer l'abonnement actif depuis la base de données
	// Pour l'exemple, on simule un abonnement
	subscription := models.Subscription{
		UserID:    userID,
		Type:      models.Monthly,
		StartDate: time.Now().Add(-15 * 24 * time.Hour),
		EndDate:   time.Now().Add(15 * 24 * time.Hour),
		Price:     30.0,
		IsActive:  true,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subscription)
}
