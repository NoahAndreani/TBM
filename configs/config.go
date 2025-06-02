package configs

import (
	"encoding/json"
	"os"
)

type Config struct {
	Server struct {
		Port string `json:"port"`
		Host string `json:"host"`
	} `json:"server"`

	Database struct {
		Type string `json:"type"`
		Path string `json:"path"`
	} `json:"database"`

	JWT struct {
		Secret          string `json:"secret"`
		ExpirationHours int    `json:"expiration_hours"`
	} `json:"jwt"`

	API struct {
		StationsEndpoint string `json:"stations_endpoint"`
		NewsEndpoint     string `json:"news_endpoint"`
	} `json:"api"`

	Experience struct {
		BaseLoginXP     int `json:"base_login_xp"`
		BaseDistanceXP  int `json:"base_distance_xp"`
		BaseTimeXP      int `json:"base_time_xp"`
		LevelMultiplier int `json:"level_multiplier"`
	} `json:"experience"`
}

var AppConfig Config

// LoadConfig charge la configuration depuis un fichier JSON
func LoadConfig() error {
	// Valeurs par défaut
	AppConfig = Config{}
	AppConfig.Server.Host = "localhost"
	AppConfig.Server.Port = "8080"
	AppConfig.Database.Type = "sqlite3"
	AppConfig.Database.Path = "data/tbcvclub.db"
	AppConfig.JWT.Secret = "votre_secret_jwt_super_securise"
	AppConfig.JWT.ExpirationHours = 24

	// Si le fichier de configuration existe, on le charge
	if _, err := os.Stat("configs/config.json"); err == nil {
		file, err := os.Open("configs/config.json")
		if err != nil {
			return err
		}
		defer file.Close()

		decoder := json.NewDecoder(file)
		if err := decoder.Decode(&AppConfig); err != nil {
			return err
		}
	}

	// Création du dossier data s'il n'existe pas
	if err := os.MkdirAll("data", 0755); err != nil {
		return err
	}

	return nil
}

// SaveConfig sauvegarde la configuration dans un fichier JSON
func SaveConfig(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	return encoder.Encode(AppConfig)
}
