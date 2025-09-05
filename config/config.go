package config

import (
	"log"

	"github.com/joho/godotenv"
)

type Config struct {
	ChargingConfig ChargingConfig
}

type ChargingConfig struct {
	Port string
}

func LoadConfig() *Config {

	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	return &Config{
		ChargingConfig: ChargingConfig{
			//Port:      os.Getenv("PORT"),
			Port: "8080",
		},
	}
}
