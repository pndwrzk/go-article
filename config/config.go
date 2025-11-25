package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost string
	DBUser string
	DBPass string
	DBName string
	DBPort string
	Port   string
}

var AppConfig *Config

func LoadConfig() {
	_ = godotenv.Load()

	AppConfig = &Config{
		DBHost: getEnv("DB_HOST", "database"),
		DBUser: getEnv("DB_USER", "admin"),
		DBPass: getEnv("DB_PASS", "secret"),
		DBName: getEnv("DB_NAME", "article_db"),
		DBPort: getEnv("DB_PORT", "5432"),
		Port:   getEnv("APP_PORT", "8080"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
