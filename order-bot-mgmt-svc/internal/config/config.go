package config

import (
	"os"
	"strconv"
	"time"
)

type App struct {
	Port int
}

type DB struct {
	Database string
	Password string
	Username string
	Port     string
	Host     string
	Schema   string
}

type Auth struct {
	AccessSecret     string
	RefreshSecret    string
	AccessTokenTTL   time.Duration
	RefreshTokenTTL  time.Duration
	UserQueryTimeout time.Duration
}

type Config struct {
	App  App
	DB   DB
	Auth Auth
}

func Load() Config {
	return Config{
		App: App{
			Port: parseIntEnv("PORT", 0),
		},
		DB: DB{
			Database: os.Getenv("BLUEPRINT_DB_DATABASE"),
			Password: os.Getenv("BLUEPRINT_DB_PASSWORD"),
			Username: os.Getenv("BLUEPRINT_DB_USERNAME"),
			Port:     os.Getenv("BLUEPRINT_DB_PORT"),
			Host:     os.Getenv("BLUEPRINT_DB_HOST"),
			Schema:   os.Getenv("BLUEPRINT_DB_SCHEMA"),
		},
		Auth: Auth{
			AccessSecret:     envOrDefault("JWT_ACCESS_SECRET", "dev-access-secret"),
			RefreshSecret:    envOrDefault("JWT_REFRESH_SECRET", "dev-refresh-secret"),
			AccessTokenTTL:   parseDurationEnv("JWT_ACCESS_TTL", 15*time.Minute),
			RefreshTokenTTL:  parseDurationEnv("JWT_REFRESH_TTL", 7*24*time.Hour),
			UserQueryTimeout: parseDurationEnv("USER_QUERY_TIMEOUT", 2*time.Second),
		},
	}
}

func envOrDefault(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func parseDurationEnv(key string, fallback time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func parseIntEnv(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}
