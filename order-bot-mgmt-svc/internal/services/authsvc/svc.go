package authsvc

import (
	"context"
	"errors"
	"order-bot-mgmt-svc/internal/config"
	"order-bot-mgmt-svc/internal/models/entities"
	"order-bot-mgmt-svc/internal/services"
	"order-bot-mgmt-svc/internal/store"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
	"order-bot-mgmt-svc/internal/models"
	"order-bot-mgmt-svc/internal/util"
)

type Svc struct {
	mu               sync.Mutex
	userStore        store.User
	refreshTokens    map[string]models.RefreshRecord
	accessSecret     []byte
	refreshSecret    []byte
	accessTokenTTL   time.Duration
	refreshTokenTTL  time.Duration
	userQueryTimeout time.Duration
}

func NewSvc(userStore store.User, cfg config.Auth) *Svc {
	return &Svc{
		userStore:        userStore,
		refreshTokens:    make(map[string]models.RefreshRecord),
		accessSecret:     []byte(cfg.AccessSecret),
		refreshSecret:    []byte(cfg.RefreshSecret),
		accessTokenTTL:   cfg.AccessTokenTTL,
		refreshTokenTTL:  cfg.RefreshTokenTTL,
		userQueryTimeout: cfg.UserQueryTimeout,
	}
}

func (s *Svc) Signup(email, password string) (models.TokenPair, error) {
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
	ctx, cancel := services.QueryContext(s.userQueryTimeout)
	defer cancel()
	if err := s.userStore.Create(ctx, newUser); err != nil {
		if errors.Is(err, store.ErrUserExists) {
			return models.TokenPair{}, ErrUserExists
		}
		return models.TokenPair{}, err
	}
	return s.issueTokens(newUser)
}

func (s *Svc) Login(email, password string) (models.TokenPair, error) {
	if s.userStore == nil {
		return models.TokenPair{}, errors.New("user store not configured")
	}
	ctx, cancel := services.QueryContext(s.userQueryTimeout)
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

func (s *Svc) Logout(refreshToken string) error {
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

func (s *Svc) issueTokens(user entities.User) (models.TokenPair, error) {
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
