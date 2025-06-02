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
		Host     string `json:"host"`
		Port     string `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		DBName   string `json:"dbname"`
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
func LoadConfig(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	return decoder.Decode(&AppConfig)
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
