package services

import (
	"context"
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
	mu               sync.Mutex
	userStore        UserStore
	refreshTokens    map[string]models.RefreshRecord
	accessSecret     []byte
	refreshSecret    []byte
	accessTokenTTL   time.Duration
	refreshTokenTTL  time.Duration
	userQueryTimeout time.Duration
}

type UserStore interface {
	Create(ctx context.Context, user models.User) error
	FindByEmail(ctx context.Context, email string) (models.User, error)
}

func NewService(userStore UserStore) *Service {
	accessSecret := os.Getenv("JWT_ACCESS_SECRET")
	refreshSecret := os.Getenv("JWT_REFRESH_SECRET")
	if accessSecret == "" {
		accessSecret = "dev-access-secret"
	}
	if refreshSecret == "" {
		refreshSecret = "dev-refresh-secret"
	}
	return &Service{
		userStore:        userStore,
		refreshTokens:    make(map[string]models.RefreshRecord),
		accessSecret:     []byte(accessSecret),
		refreshSecret:    []byte(refreshSecret),
		accessTokenTTL:   parseDurationEnv("JWT_ACCESS_TTL", 15*time.Minute),
		refreshTokenTTL:  parseDurationEnv("JWT_REFRESH_TTL", 7*24*time.Hour),
		userQueryTimeout: parseDurationEnv("USER_QUERY_TIMEOUT", 2*time.Second),
	}
}
