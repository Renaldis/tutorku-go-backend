package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort        string
	AppEnv         string
	DBHost         string
	DBPort         string
	DBName         string
	DBUser         string
	DBPass         string
	JWTSecret      string
	JWTExpiresHour int
	N8NBaseURL     string
	N8NWebhooks    map[string]string
}

var Cfg *Config

func Load() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file, reading from environment")
	}

	expiresHour, _ := strconv.Atoi(getEnv("JWT_EXPIRES_HOUR", "72"))

	Cfg = &Config{
		AppPort:        getEnv("APP_PORT", "8080"),
		AppEnv:         getEnv("APP_ENV", "development"),
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnv("DB_PORT", "5432"),
		DBName:         getEnv("DB_NAME", "tutorku_db"),
		DBUser:         getEnv("DB_USER", "root"),
		DBPass:         getEnv("DB_PASS", "password"),
		JWTSecret:      getEnv("JWT_SECRET", "secret"),
		JWTExpiresHour: expiresHour,
		N8NBaseURL:     getEnv("N8N_BASE_URL", ""),
		N8NWebhooks: map[string]string{
			"ingest":    getEnv("N8N_WEBHOOK_INGEST", ""),
			"chat":      getEnv("N8N_WEBHOOK_CHAT", ""),
			"summarize": getEnv("N8N_WEBHOOK_SUMMARIZE", ""),
			"quiz":      getEnv("N8N_WEBHOOK_QUIZ", ""),
			"essay":     getEnv("N8N_WEBHOOK_ESSAY", ""),
		},
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
