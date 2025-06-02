package models

import (
	"time"
)

type User struct {
	ID              int64     `json:"id" db:"id"`
	Username        string    `json:"username" db:"username"`
	Email           string    `json:"email" db:"email"`
	HashedPassword  string    `json:"-" db:"hashed_password"`
	Role            string    `json:"role" db:"role"`
	Level           int       `json:"level" db:"level"`
	Experience      int       `json:"experience" db:"experience"`
	TotalDistance   float64   `json:"total_distance" db:"total_distance"`
	TotalRideTime   int       `json:"total_ride_time" db:"total_ride_time"`
	ConsecutiveDays int       `json:"consecutive_days" db:"consecutive_days"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	LastConnection  time.Time `json:"last_connection" db:"last_connection"`
}

type UserStats struct {
	Level            int     `json:"level"`
	Experience       int     `json:"experience"`
	ExperienceToNext int     `json:"experience_to_next"`
	TotalDistance    float64 `json:"total_distance"`
	TotalRideTime    int     `json:"total_ride_time"`
	ConsecutiveDays  int     `json:"consecutive_days"`
}

// ExperienceForNextLevel calcule l'expérience nécessaire pour le niveau suivant
func (u *User) ExperienceForNextLevel() int {
	return u.Level * 1000
}

// AddRideExperience ajoute de l'expérience basée sur une course
func (u *User) AddRideExperience(distance float64, time int) {
	// Base XP : 100 points par km et 1 point par minute
	xp := int(distance*100) + time

	// Bonus basé sur le niveau actuel
	bonus := 1.0
	if u.Level >= 5 {
		bonus = 1.5
	}
	if u.Level >= 10 {
		bonus = 2.0
	}
	if u.Level >= 20 {
		bonus = 3.0
	}

	u.Experience += int(float64(xp) * bonus)
	u.TotalDistance += distance
	u.TotalRideTime += time

	// Vérification du passage de niveau
	for u.Experience >= u.ExperienceForNextLevel() {
		u.Level++
	}
}

// AddDailyLoginExperience ajoute de l'expérience pour la connexion quotidienne
func (u *User) AddDailyLoginExperience() {
	// Ajoute 10 points d'expérience par connexion quotidienne
	u.Experience += 10

	// Vérifie si l'utilisateur peut monter de niveau
	// On monte de niveau tous les 100 points d'expérience
	newLevel := (u.Experience / 100) + 1
	if newLevel > u.Level {
		u.Level = newLevel
	}
}
