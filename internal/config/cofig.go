package config

import (
	"os"

	"github.com/joho/godotenv"
)

var AppConfig *Config

type Config struct {
	APP_PORT    string
	DB_HOST     string
	DB_PORT     string
	DB_USER     string
	DB_PASS     string
	DB_NAME     string
	REDIS_ADDR  string
	REDIS_PASS  string
	SMTP_HOST   string
	SMTP_PORT   string
	SMTP_USER   string
	SMTP_PASS   string
	SMTP_SENDER string
	JWT_SECRET  string
}

func Load() *Config {
	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}

	envFile := ".env." + env

	_ = godotenv.Load(envFile)
	AppConfig = &Config{
		APP_PORT:    getEnv("APP_PORT", "8080"),
		DB_HOST:     getEnv("POSTGRES_HOST", "localhost"),
		DB_PORT:     getEnv("POSTGRES_PORT", "5432"),
		DB_USER:     getEnv("POSTGRES_USER", "admin"),
		DB_PASS:     getEnv("POSTGRES_PASSWORD", "admin123"),
		DB_NAME:     getEnv("POSTGRES_DB", "jadwalin"),
		REDIS_ADDR:  getEnv("REDIS_HOST", "localhost") + ":" + getEnv("REDIS_PORT", "6379"),
		REDIS_PASS:  getEnv("REDIS_PASSWORD", ""),
		SMTP_HOST:   getEnv("SMTP_HOST", ""),
		SMTP_PORT:   getEnv("SMTP_PORT", ""),
		SMTP_USER:   getEnv("SMTP_USER", ""),
		SMTP_PASS:   getEnv("SMTP_PASSWORD", ""),
		SMTP_SENDER: getEnv("SMTP_SENDER", ""),
		JWT_SECRET:  getEnv("JWT_SECRET", "your-secret-key"),
	}
	return AppConfig
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
