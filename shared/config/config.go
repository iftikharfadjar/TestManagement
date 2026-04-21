package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	DBType       string `json:"DB_TYPE"`
	Port         string `json:"PORT"`
	DBConnString string `json:"DB_CONN_STRING"`
	JWTSecret    string `json:"JWT_SECRET"`
}

func LoadConfig() *Config {
	cfg := &Config{}

	file, err := os.Open("config.json")
	if err == nil {
		defer file.Close()
		decoder := json.NewDecoder(file)
		if err := decoder.Decode(cfg); err != nil {
			log.Printf("Warning: failed to decode config.json: %v", err)
		}
	} else {
		log.Println("config.json not found, relying purely on environment variables")
	}

	if cfg.DBType == "" {
		cfg.DBType = os.Getenv("DB_TYPE")
	}
	if cfg.Port == "" {
		cfg.Port = os.Getenv("PORT")
	}
	if cfg.DBConnString == "" {
		cfg.DBConnString = os.Getenv("DB_CONN_STRING")
	}
	if cfg.JWTSecret == "" {
		cfg.JWTSecret = os.Getenv("JWT_SECRET")
	}

	if cfg.DBType == "" {
		cfg.DBType = "sqlite"
	}
	if cfg.Port == "" {
		cfg.Port = "3000"
	}

	return cfg
}

