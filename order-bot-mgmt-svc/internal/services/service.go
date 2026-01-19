package services

import (
	"errors"
	"os"
	"sync"
	"time"

	"order-bot-mgmt-svc/internal/models"
)

var (
	ErrUserExists         = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
)

type Service struct {
	mu              sync.Mutex
	usersByEmail    map[string]models.User
	refreshTokens   map[string]models.RefreshRecord
	accessSecret    []byte
	refreshSecret   []byte
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewService() *Service {
	accessSecret := os.Getenv("JWT_ACCESS_SECRET")
	refreshSecret := os.Getenv("JWT_REFRESH_SECRET")
	if accessSecret == "" {
		accessSecret = "dev-access-secret"
	}
	if refreshSecret == "" {
		refreshSecret = "dev-refresh-secret"
	}
	return &Service{
		usersByEmail:    make(map[string]models.User),
		refreshTokens:   make(map[string]models.RefreshRecord),
		accessSecret:    []byte(accessSecret),
		refreshSecret:   []byte(refreshSecret),
		accessTokenTTL:  parseDurationEnv("JWT_ACCESS_TTL", 15*time.Minute),
		refreshTokenTTL: parseDurationEnv("JWT_REFRESH_TTL", 7*24*time.Hour),
	}
}
