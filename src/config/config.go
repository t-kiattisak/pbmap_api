package config

import (
	"os"
	"strconv"
)

type Config struct {
	AppPort   int
	DBHost    string
	DBUser    string
	DBPass    string
	DBName    string
	DBPort    string
	DBSSL     string
	DBZone    string
	JWTSecret string
}

func LoadConfig() *Config {
	portStr := os.Getenv("APP_PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		port = 3000
	}

	return &Config{
		AppPort:   port,
		DBHost:    getEnv("DB_HOST", "localhost"),
		DBUser:    getEnv("DB_USER", "postgres"),
		DBPass:    getEnv("DB_PASS", "password"),
		DBName:    getEnv("DB_NAME", "pbmap_db"),
		DBPort:    getEnv("DB_PORT", "5432"),
		DBSSL:     getEnv("DB_SSL", "disable"),
		DBZone:    getEnv("DB_ZONE", "Asia/Bangkok"),
		JWTSecret: getEnv("JWT_SECRET", "super-secret-key"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
