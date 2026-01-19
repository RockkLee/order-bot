package services

import (
	"crypto/rand"
	"encoding/hex"
	"golang.org/x/crypto/bcrypt"
	"order-bot-mgmt-svc/internal/models"
	"os"
	"time"
)

func (s *Service) Signup(email, password string) (models.TokenPair, error) {
	if email == "" || password == "" {
		return models.TokenPair{}, ErrInvalidCredentials
	}
	s.mu.Lock()
	if _, exists := s.usersByEmail[email]; exists {
		s.mu.Unlock()
		return models.TokenPair{}, ErrUserExists
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		s.mu.Unlock()
		return models.TokenPair{}, err
	}
	user := models.User{
		ID:           newID(),
		Email:        email,
		PasswordHash: string(hash),
	}
	s.usersByEmail[email] = user
	s.mu.Unlock()
	return s.issueTokens(user)
}

func (s *Service) Login(email, password string) (models.TokenPair, error) {
	s.mu.Lock()
	user, exists := s.usersByEmail[email]
	s.mu.Unlock()
	if !exists {
		return models.TokenPair{}, ErrInvalidCredentials
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

func newID() string {
	buf := make([]byte, 16)
	_, err := rand.Read(buf)
	if err != nil {
		return time.Now().UTC().Format("20060102150405.000000000")
	}
	return hex.EncodeToString(buf)
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
