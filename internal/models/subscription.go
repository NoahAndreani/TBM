package models

import (
	"time"
)

type SubscriptionType string

const (
	Daily   SubscriptionType = "daily"
	Weekly  SubscriptionType = "weekly"
	Monthly SubscriptionType = "monthly"
	Yearly  SubscriptionType = "yearly"
)

type Subscription struct {
	ID        int64            `json:"id" db:"id"`
	UserID    int64            `json:"user_id" db:"user_id"`
	Type      SubscriptionType `json:"type" db:"type"`
	StartDate time.Time        `json:"start_date" db:"start_date"`
	EndDate   time.Time        `json:"end_date" db:"end_date"`
	Price     float64          `json:"price" db:"price"`
	IsActive  bool             `json:"is_active" db:"is_active"`
}

type SubscriptionPlan struct {
	ID          int64            `json:"id"`
	Type        SubscriptionType `json:"type"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Price       float64          `json:"price"`
	Duration    time.Duration    `json:"duration"`
	Benefits    []string         `json:"benefits"`
}

// Vérifie si l'abonnement est valide à une date donnée
func (s *Subscription) IsValidAt(date time.Time) bool {
	return s.IsActive && date.After(s.StartDate) && date.Before(s.EndDate)
}

// Vérifie si l'abonnement est expiré
func (s *Subscription) IsExpired() bool {
	return time.Now().After(s.EndDate)
}

// Calcule le temps restant de l'abonnement
func (s *Subscription) RemainingTime() time.Duration {
	if s.IsExpired() {
		return 0
	}
	return time.Until(s.EndDate)
}

// Renouvelle l'abonnement pour une nouvelle période
func (s *Subscription) Renew(plan SubscriptionPlan) {
	s.StartDate = time.Now()
	s.EndDate = s.StartDate.Add(plan.Duration)
	s.Price = plan.Price
	s.IsActive = true
}
