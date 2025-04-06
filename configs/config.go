package configs

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	Db DbConfig
}

type DbConfig struct {
	Dsn string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
	return &Config{
		Db: DbConfig{
			Dsn: os.Getenv("DSN"),
		},
	}
}
