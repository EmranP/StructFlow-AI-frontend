package configs

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort string

	DBHost string
	DBPort string
	DBUser string
	DBPass string
	DBName string

	OpenAIKey string

	JWTSecret string
}

func Load() (*Config, error) {
	err := godotenv.Load()

	if err != nil {
		return nil, err
	}

	return &Config{
		AppPort: os.Getenv("APP_PORT"),

		DBHost: os.Getenv("DB_HOST"),
		DBPort: os.Getenv("DB_PORT"),
		DBUser: os.Getenv("DB_USER"),
		DBPass: os.Getenv("DB_PASS"),
		DBName: os.Getenv("DB_NAME"),

		OpenAIKey: os.Getenv("OPENAI_API_KEY"),

		JWTSecret: os.Getenv("JWT_SECRET"),
	}, nil
}
