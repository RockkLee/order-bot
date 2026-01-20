package services

import (
	"context"
	"errors"
	"order-bot-mgmt-svc/internal/models/entities"
	"order-bot-mgmt-svc/internal/store"
	"order-bot-mgmt-svc/internal/util/jwtutil"
	"os"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
	"order-bot-mgmt-svc/internal/models"
	"order-bot-mgmt-svc/internal/util"
)

type AuthService struct {
	mu               sync.Mutex
	userStore        store.User
	refreshTokens    map[string]models.RefreshRecord
	accessSecret     []byte
	refreshSecret    []byte
	accessTokenTTL   time.Duration
	refreshTokenTTL  time.Duration
	userQueryTimeout time.Duration
}

func NewAuthService(userStore store.User) *AuthService {
	accessSecret := os.Getenv("JWT_ACCESS_SECRET")
	refreshSecret := os.Getenv("JWT_REFRESH_SECRET")
	if accessSecret == "" {
		accessSecret = "dev-access-secret"
	}
	if refreshSecret == "" {
		refreshSecret = "dev-refresh-secret"
	}
	return &AuthService{
		userStore:        userStore,
		refreshTokens:    make(map[string]models.RefreshRecord),
		accessSecret:     []byte(accessSecret),
		refreshSecret:    []byte(refreshSecret),
		accessTokenTTL:   parseDurationEnv("JWT_ACCESS_TTL", 15*time.Minute),
		refreshTokenTTL:  parseDurationEnv("JWT_REFRESH_TTL", 7*24*time.Hour),
		userQueryTimeout: parseDurationEnv("USER_QUERY_TIMEOUT", 2*time.Second),
	}
}

func (s *AuthService) Signup(email, password string) (models.TokenPair, error) {
	if email == "" || password == "" {
		return models.TokenPair{}, ErrInvalidCredentials
	}
	if s.userStore == nil {
		return models.TokenPair{}, errors.New("user store not configured")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return models.TokenPair{}, err
	}
	newUser := entities.User{
		ID:           util.NewID(),
		Email:        email,
		PasswordHash: string(hash),
	}
	ctx, cancel := s.userContext()
	defer cancel()
	if err := s.userStore.Create(ctx, newUser); err != nil {
		if errors.Is(err, store.ErrUserExists) {
			return models.TokenPair{}, ErrUserExists
		}
		return models.TokenPair{}, err
	}
	return s.issueTokens(newUser)
}

func (s *AuthService) Login(email, password string) (models.TokenPair, error) {
	if s.userStore == nil {
		return models.TokenPair{}, errors.New("user store not configured")
	}
	ctx, cancel := s.userContext()
	defer cancel()
	user, err := s.userStore.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return models.TokenPair{}, ErrInvalidCredentials
		}
		return models.TokenPair{}, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return models.TokenPair{}, ErrInvalidCredentials
	}
	return s.issueTokens(user)
}

func (s *AuthService) Logout(refreshToken string) error {
	if refreshToken == "" {
		return ErrInvalidToken
	}
	claims, err := jwtutil.ParseJWT(s.refreshSecret, refreshToken)
	if err != nil || claims.Typ != "refresh" {
		return ErrInvalidToken
	}
	s.mu.Lock()
	if _, exists := s.refreshTokens[refreshToken]; !exists {
		s.mu.Unlock()
		return ErrInvalidToken
	}
	delete(s.refreshTokens, refreshToken)
	s.mu.Unlock()
	return nil
}

func (s *AuthService) issueTokens(user entities.User) (models.TokenPair, error) {
	now := time.Now()
	accessClaims := models.Claims{
		Sub:   user.ID,
		Email: user.Email,
		Exp:   now.Add(s.accessTokenTTL).Unix(),
		Iat:   now.Unix(),
		Typ:   "access",
	}
	refreshClaims := models.Claims{
		Sub:   user.ID,
		Email: user.Email,
		Exp:   now.Add(s.refreshTokenTTL).Unix(),
		Iat:   now.Unix(),
		Typ:   "refresh",
	}
	accessToken, err := jwtutil.SignJWT(s.accessSecret, accessClaims)
	if err != nil {
		return models.TokenPair{}, err
	}
	refreshToken, err := jwtutil.SignJWT(s.refreshSecret, refreshClaims)
	if err != nil {
		return models.TokenPair{}, err
	}
	s.mu.Lock()
	s.refreshTokens[refreshToken] = models.RefreshRecord{
		UserID:    user.ID,
		ExpiresAt: now.Add(s.refreshTokenTTL),
	}
	s.mu.Unlock()
	return models.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) userContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), s.userQueryTimeout)
}
