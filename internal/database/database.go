package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"tbcvclub/internal/models"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

// InitDB initialise la connexion à la base de données
func InitDB(dataPath string) error {
	if err := os.MkdirAll(filepath.Dir(dataPath), 0755); err != nil {
		return fmt.Errorf("erreur création dossier DB: %v", err)
	}

	var err error
	db, err = sql.Open("sqlite3", dataPath)
	if err != nil {
		return fmt.Errorf("erreur connexion DB: %v", err)
	}

	// Test de la connexion
	if err = db.Ping(); err != nil {
		return fmt.Errorf("erreur ping DB: %v", err)
	}

	// Création des tables
	schema, err := os.ReadFile("internal/database/schema.sql")
	if err != nil {
		return fmt.Errorf("erreur lecture schema: %v", err)
	}

	if _, err = db.Exec(string(schema)); err != nil {
		return fmt.Errorf("erreur création tables: %v", err)
	}

	// Création du compte admin par défaut s'il n'existe pas
	if err := createDefaultAdmin(); err != nil {
		return fmt.Errorf("erreur création admin: %v", err)
	}

	return nil
}

// createDefaultAdmin crée un compte administrateur par défaut s'il n'existe pas
func createDefaultAdmin() error {
	// Vérifie si un admin existe déjà
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE role = 'admin'").Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		// Création du compte admin par défaut
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		now := time.Now()
		admin := &models.User{
			Username:       "admin",
			Email:          "admin@tbcvclub.fr",
			HashedPassword: string(hashedPassword),
			Role:           "admin",
			Level:          1,
			CreatedAt:      now,
			LastConnection: now,
		}

		return CreateUser(admin)
	}

	return nil
}

// CreateUser crée un nouvel utilisateur
func CreateUser(user *models.User) error {
	query := `
		INSERT INTO users (
			username, email, hashed_password, role, level, experience,
			total_distance, total_ride_time, consecutive_days
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := db.Exec(query,
		user.Username,
		user.Email,
		user.HashedPassword,
		user.Role,
		user.Level,
		user.Experience,
		user.TotalDistance,
		user.TotalRideTime,
		user.ConsecutiveDays,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = id
	return nil
}

// GetUserByUsername récupère un utilisateur par son nom d'utilisateur
func GetUserByUsername(username string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, username, email, hashed_password, role, level, experience,
		       total_distance, total_ride_time, consecutive_days, created_at,
		       last_connection
		FROM users WHERE username = ?
	`

	err := db.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.HashedPassword,
		&user.Role,
		&user.Level,
		&user.Experience,
		&user.TotalDistance,
		&user.TotalRideTime,
		&user.ConsecutiveDays,
		&user.CreatedAt,
		&user.LastConnection,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByEmail récupère un utilisateur par son email
func GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, username, email, hashed_password, role, level, experience,
		       total_distance, total_ride_time, consecutive_days, created_at,
		       last_connection
		FROM users WHERE email = ?
	`

	err := db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.HashedPassword,
		&user.Role,
		&user.Level,
		&user.Experience,
		&user.TotalDistance,
		&user.TotalRideTime,
		&user.ConsecutiveDays,
		&user.CreatedAt,
		&user.LastConnection,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateUser met à jour les informations d'un utilisateur
func UpdateUser(user *models.User) error {
	query := `
		UPDATE users SET
			username = ?, email = ?, level = ?, experience = ?,
			total_distance = ?, total_ride_time = ?, consecutive_days = ?,
			last_connection = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := db.Exec(query,
		user.Username,
		user.Email,
		user.Level,
		user.Experience,
		user.TotalDistance,
		user.TotalRideTime,
		user.ConsecutiveDays,
		user.ID,
	)

	return err
}

// UpdatePassword met à jour le mot de passe d'un utilisateur
func UpdatePassword(userID int, hashedPassword string) error {
	query := "UPDATE users SET hashed_password = ? WHERE id = ?"
	_, err := db.Exec(query, hashedPassword, userID)
	return err
}

// Close ferme la connexion à la base de données
func Close() error {
	if db != nil {
		return db.Close()
	}
	return nil
}
