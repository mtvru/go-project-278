package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL  string
	BaseURL      string
	SentryDSN    string
	AllowOrigins []string
}

func Load() Config {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		log.Printf("load .env: %v", err)
	}

	return Config{
		DatabaseURL:  os.Getenv("DATABASE_URL"),
		BaseURL:      getEnv("BASE_URL", "http://localhost:8080"),
		SentryDSN:    os.Getenv("SENTRY_DSN"),
		AllowOrigins: splitAndTrim(getEnv("ALLOW_ORIGINS", "http://localhost:5173")),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func splitAndTrim(s string) []string {
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
