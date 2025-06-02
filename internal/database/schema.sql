-- Schéma de la base de données TBC Vclub

-- Table des utilisateurs
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username VARCHAR(20) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    hashed_password VARCHAR(255) NOT NULL,
    role VARCHAR(10) NOT NULL DEFAULT 'user',
    level INTEGER NOT NULL DEFAULT 1,
    experience INTEGER NOT NULL DEFAULT 0,
    total_distance REAL NOT NULL DEFAULT 0,
    total_ride_time INTEGER NOT NULL DEFAULT 0,
    consecutive_days INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_connection TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Création du compte admin par défaut
INSERT OR IGNORE INTO users (username, email, hashed_password, role)
VALUES (
    'admin',
    'admin@tbcvclub.fr',
    -- Le mot de passe par défaut est 'admin123'
    '$2a$10$8KzaNdKwZ8hQZXlmkfg5e.9uZ6OBgBkZ7eW5xvj.kn5uqzQGy7k6a',
    'admin'
); 