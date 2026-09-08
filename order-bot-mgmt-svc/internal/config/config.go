package config

import (
	"os"
	"strconv"
	"time"
)

type App struct {
	Address string
	Port    int
	GinMode string
}

type Grpc struct {
	Address string
	Port    int
}

type Db struct {
	Database string
	Password string
	Username string
	Port     string
	Host     string
	Schema   string
}

type Auth struct {
	AccessSecret    string
	RefreshSecret   string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

type Others struct {
	QryCtxTimeout time.Duration
}

type Config struct {
	App          App
	Grpc         Grpc
	OrderBotGrpc Grpc
	Db           Db
	OrderBotDb   Db
	Auth         Auth
	Others       Others
}

func Load() Config {
	return Config{
		App: App{
			Address: getEnv("ADDRESS"),
			Port:    getIntEnv("PORT"),
			GinMode: getEnv("GIN_MODE"),
		},
		Grpc: Grpc{
			Address: getEnv("GRPC_ADDRESS"),
			Port:    getIntEnv("GRPC_PORT"),
		},
		OrderBotGrpc: Grpc{
			Address: envOrDefault("ORDER_BOT_GRPC_ADDRESS", getEnv("GRPC_ADDRESS")),
			Port:    getIntEnvOrDefault("ORDER_BOT_GRPC_PORT", getIntEnv("GRPC_PORT")),
		},
		Db: Db{
			Database: getEnv("BLUEPRINT_DB_DATABASE"),
			Password: getEnv("BLUEPRINT_DB_PASSWORD"),
			Username: getEnv("BLUEPRINT_DB_USERNAME"),
			Port:     getEnv("BLUEPRINT_DB_PORT"),
			Host:     getEnv("BLUEPRINT_DB_HOST"),
			Schema:   getEnv("BLUEPRINT_DB_SCHEMA"),
		},
		OrderBotDb: Db{
			Database: getEnv("BLUEPRINT_DB_DATABASE"),
			Password: getEnv("BLUEPRINT_DB_PASSWORD"),
			Username: getEnv("BLUEPRINT_DB_USERNAME"),
			Port:     getEnv("BLUEPRINT_DB_PORT"),
			Host:     getEnv("BLUEPRINT_DB_HOST"),
			Schema:   getEnv("BLUEPRINT_DB_ORDER_BOT_SCHEMA"),
		},
		Auth: Auth{
			AccessSecret:    envOrDefault("JWT_ACCESS_SECRET", "dev-access-secret"),
			RefreshSecret:   envOrDefault("JWT_REFRESH_SECRET", "dev-refresh-secret"),
			AccessTokenTTL:  parseDurationEnv("JWT_ACCESS_TTL", 30*time.Minute),
			RefreshTokenTTL: parseDurationEnv("JWT_REFRESH_TTL", 7*24*time.Hour),
		},
		Others: Others{
			QryCtxTimeout: parseDurationEnv("QRY_CTX_TIMEOUT", 15*time.Second),
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

func getIntEnvOrDefault(key string, fallback int) int {
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

func getEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		panic("config.getEnv(), env not found")
	}
	return val
}

func getIntEnv(key string) int {
	value := os.Getenv(key)
	parsed, err := strconv.Atoi(value)
	if err != nil {
		panic(err.Error())
	}
	return parsed
}
