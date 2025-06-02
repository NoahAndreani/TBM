package models

import (
	"time"
)

type User struct {
	ID             int       `json:"id"`
	Username       string    `json:"username"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"-"`
	Level          int       `json:"level"`
	Experience     int       `json:"experience"`
	TotalDistance  float64   `json:"total_distance"`
	TotalRideTime  int       `json:"total_ride_time"` // en minutes
	LastConnection time.Time `json:"last_connection"`
	CreatedAt      time.Time `json:"created_at"`
}

type UserStats struct {
	Level            int     `json:"level"`
	Experience       int     `json:"experience"`
	ExperienceToNext int     `json:"experience_to_next"`
	TotalDistance    float64 `json:"total_distance"`
	TotalRideTime    int     `json:"total_ride_time"`
	ConsecutiveDays  int     `json:"consecutive_days"`
}

// Calcule l'expérience nécessaire pour le prochain niveau
func (u *User) ExperienceForNextLevel() int {
	return 1000 * u.Level // Formule simple : chaque niveau nécessite 1000 * niveau actuel
}

// Ajoute de l'expérience et met à jour le niveau si nécessaire
func (u *User) AddExperience(exp int) {
	u.Experience += exp
	for u.Experience >= u.ExperienceForNextLevel() {
		u.Experience -= u.ExperienceForNextLevel()
		u.Level++
	}
}

// Ajoute de l'expérience basée sur la distance parcourue
func (u *User) AddRideExperience(distance float64, timeInMinutes int) {
	// 10 XP par km + 1 XP par minute de trajet
	experienceGained := int(distance*10) + timeInMinutes
	u.AddExperience(experienceGained)
	u.TotalDistance += distance
	u.TotalRideTime += timeInMinutes
}

// Ajoute de l'expérience pour une connexion journalière
func (u *User) AddDailyLoginExperience() {
	if time.Since(u.LastConnection) >= 24*time.Hour {
		u.AddExperience(100) // 100 XP pour une connexion journalière
		u.LastConnection = time.Now()
	}
}
