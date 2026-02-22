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
	QryCtxTimeout    time.Duration
	OrderSvcGRPCAddr string
	MgmtSvcGRPCAddr  string
}

type Config struct {
	App        App
	Db         Db
	OrderBotDb Db
	Auth       Auth
	Others     Others
}

func Load() Config {
	return Config{
		App: App{
			Address: getEnv("ADDRESS"),
			Port:    getIntEnv("PORT"),
			GinMode: getEnv("GIN_MODE"),
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
			QryCtxTimeout:    parseDurationEnv("QRY_CTX_TIMEOUT", 15*time.Second),
			OrderSvcGRPCAddr: envOrDefault("ORDER_SVC_GRPC_ADDR", "localhost:50052"),
			MgmtSvcGRPCAddr:  envOrDefault("MGMT_SVC_GRPC_ADDR", "0.0.0.0:50051"),
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
