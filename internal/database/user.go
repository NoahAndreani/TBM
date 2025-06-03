package database

import (
	"database/sql"
	"tbcvclub/internal/models"
)

// GetUserByID récupère un utilisateur par son ID
func GetUserByID(id int64) (*models.User, error) {
	var user models.User
	err := db.QueryRow(`
		SELECT id, username, email, hashed_password, role, level, experience,
		       total_distance, total_ride_time, consecutive_days, created_at, last_connection
		FROM users WHERE id = ?`, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.HashedPassword, &user.Role,
		&user.Level, &user.Experience, &user.TotalDistance, &user.TotalRideTime,
		&user.ConsecutiveDays, &user.CreatedAt, &user.LastConnection)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// ... existing code ...
