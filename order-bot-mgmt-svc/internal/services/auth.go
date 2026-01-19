package services

import (
	"context"
	"errors"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"
	"order-bot-mgmt-svc/internal/models"
	postgresuser "order-bot-mgmt-svc/internal/postgres/user"
	"order-bot-mgmt-svc/internal/util"
)

func (s *Service) Signup(email, password string) (models.TokenPair, error) {
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
	newUser := models.User{
		ID:           util.NewID(),
		Email:        email,
		PasswordHash: string(hash),
	}
	ctx, cancel := s.userContext()
	defer cancel()
	if err := s.userStore.Create(ctx, newUser); err != nil {
		if errors.Is(err, postgresuser.ErrUserExists) {
			return models.TokenPair{}, ErrUserExists
		}
		return models.TokenPair{}, err
	}
	return s.issueTokens(newUser)
}

func (s *Service) Login(email, password string) (models.TokenPair, error) {
	if s.userStore == nil {
		return models.TokenPair{}, errors.New("user store not configured")
	}
	ctx, cancel := s.userContext()
	defer cancel()
	user, err := s.userStore.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, postgresuser.ErrNotFound) {
			return models.TokenPair{}, ErrInvalidCredentials
		}
		return models.TokenPair{}, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return models.TokenPair{}, ErrInvalidCredentials
	}
	return s.issueTokens(user)
}

func (s *Service) Logout(refreshToken string) error {
	if refreshToken == "" {
		return ErrInvalidToken
	}
	claims, err := parseJWT(s.refreshSecret, refreshToken)
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

func (s *Service) issueTokens(user models.User) (models.TokenPair, error) {
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
	accessToken, err := signJWT(s.accessSecret, accessClaims)
	if err != nil {
		return models.TokenPair{}, err
	}
	refreshToken, err := signJWT(s.refreshSecret, refreshClaims)
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

func (s *Service) userContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), s.userQueryTimeout)
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
