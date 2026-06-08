package configs

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort string

	DBHost string
	DBPort string
	DBUser string
	DBPass string
	DBName string

	AIKey           string
	ResendEmailKey  string
	ResendEmailFrom string

	SmtpHost     string
	SmtpPort     int
	SmtpEmail    string
	SmtpPassword string

	JWTSecret string
	OriginUrl string
}

func Load() (*Config, error) {
	err := godotenv.Load()

	if err != nil {
		return nil, err
	}

	smtpPort, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		return nil, err
	}

	return &Config{
		AppPort:   os.Getenv("APP_PORT"),
		OriginUrl: os.Getenv("CLIENT_URL"),

		DBHost: os.Getenv("DB_HOST"),
		DBPort: os.Getenv("DB_PORT"),
		DBUser: os.Getenv("DB_USER"),
		DBPass: os.Getenv("DB_PASS"),
		DBName: os.Getenv("DB_NAME"),

		AIKey:           os.Getenv("API_AI_KEY"),
		ResendEmailKey:  os.Getenv("API_RESEND_EMAIL_KEY"),
		ResendEmailFrom: os.Getenv("API_RESEND_EMAIL_FROM"),

		SmtpHost:     os.Getenv("SMTP_HOST"),
		SmtpPort:     smtpPort,
		SmtpEmail:    os.Getenv("SMTP_EMAIL"),
		SmtpPassword: os.Getenv("SMTP_PASSWORD"),

		JWTSecret: os.Getenv("JWT_SECRET"),
	}, nil
}
