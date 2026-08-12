package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	Port        string
	DatabaseDSN string
}

func Load() (Config, error) {
	cfg := Config{
		Port:        normalizePort(getEnv("PORT", "8080")),
		DatabaseDSN: strings.TrimSpace(os.Getenv("DB_DSN")),
	}
	if cfg.DatabaseDSN == "" {
		return Config{}, fmt.Errorf("DB_DSN is required")
	}
	return cfg, nil
}

func getEnv(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func normalizePort(port string) string {
	value := strings.TrimSpace(port)
	if value == "" {
		return ":8080"
	}
	return value
}
