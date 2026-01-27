package authsvc

import (
	"errors"
	"fmt"
	"order-bot-mgmt-svc/internal/config"
	"order-bot-mgmt-svc/internal/models"
	"order-bot-mgmt-svc/internal/models/entities"
	"order-bot-mgmt-svc/internal/store"
	"order-bot-mgmt-svc/internal/util"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Svc struct {
	userStore       store.User
	accessSecret    []byte
	refreshSecret   []byte
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
	ctxFunc         util.CtxFunc
}

func NewSvc(userStore store.User, cfg config.Config, ctxFunc util.CtxFunc) *Svc {
	return &Svc{
		userStore:       userStore,
		accessSecret:    []byte(cfg.Auth.AccessSecret),
		refreshSecret:   []byte(cfg.Auth.RefreshSecret),
		accessTokenTTL:  cfg.Auth.AccessTokenTTL,
		refreshTokenTTL: cfg.Auth.RefreshTokenTTL,
		ctxFunc:         ctxFunc,
	}
}

func (s *Svc) Signup(email, password string) (models.TokenPair, error) {
	if email == "" || password == "" {
		return models.TokenPair{}, fmt.Errorf("authsvc.Signup: %w", ErrInvalidCredentials)
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return models.TokenPair{}, fmt.Errorf("authsvc.Signup: %w", err)
	}
	newUser := entities.User{
		ID:           util.NewID(),
		Email:        email,
		PasswordHash: string(hash),
		AccessToken:  "",
		RefreshToken: "",
	}
	ctx, cancel := s.ctxFunc()
	defer cancel()
	if err := s.userStore.Create(ctx, newUser); err != nil {
		if errors.Is(err, store.ErrUserExists) {
			return models.TokenPair{}, fmt.Errorf("authsvc.Signup: %w", ErrUserExists)
		}
		return models.TokenPair{}, fmt.Errorf("authsvc.Signup: %w", err)
	}
	tokens, err := s.issueTokens(newUser)
	if err != nil {
		return models.TokenPair{}, fmt.Errorf("authsvc.Signup: %w", err)
	}
	return tokens, nil
}

func (s *Svc) Login(email, password string) (models.TokenPair, error) {
	if s.userStore == nil {
		return models.TokenPair{}, fmt.Errorf("authsvc.Login: %w", errors.New("user store not configured"))
	}
	ctx, cancel := s.ctxFunc()
	defer cancel()
	user, err := s.userStore.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return models.TokenPair{}, fmt.Errorf("authsvc.Login: %w", ErrInvalidCredentials)
		}
		return models.TokenPair{}, fmt.Errorf("authsvc.Login: %w", err)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return models.TokenPair{}, fmt.Errorf("authsvc.Login: %w", ErrInvalidCredentials)
	}
	tokens, err := s.issueTokens(user)
	if err != nil {
		return models.TokenPair{}, fmt.Errorf("authsvc.Login: %w", err)
	}
	return tokens, nil
}

func (s *Svc) Logout(refreshToken string) error {
	userID, err := s.ValidateRefreshToken(refreshToken)
	if err != nil {
		return fmt.Errorf("authsvc.Logout: %w", err)
	}
	ctx, cancel := s.ctxFunc()
	defer cancel()
	if err := s.userStore.UpdateTokens(ctx, userID, "", ""); err != nil {
		return fmt.Errorf("authsvc.Logout: %w", err)
	}
	return nil
}

func (s *Svc) ValidateAccessToken(accessToken string) error {
	if accessToken == "" {
		return fmt.Errorf("authsvc.ValidateAccessToken(): accessToken is empty %w", ErrInvalidToken)
	}
	claims, err := parseJWT(s.accessSecret, accessToken)
	if err != nil {
		return fmt.Errorf("authsvc.ValidateAccessToken(): %w", err)
	}
	if claims.Typ != "access" {
		return fmt.Errorf("authsvc.ValidateAccessToken(), claims.Typ != 'access': %w", ErrInvalidToken)
	}
	if claims.Exp > time.Now().UnixMilli() {
		return fmt.Errorf("authsvc.ValidateAccessToken(), accessToken is expired: %w", ErrExpiredToken)
	}
	return nil
}

func (s *Svc) ValidateRefreshToken(refreshToken string) (string, error) {
	if refreshToken == "" {
		return "", fmt.Errorf("authsvc.ValidateRefreshToken(): %w", ErrInvalidToken)
	}
	claims, err := parseJWT(s.refreshSecret, refreshToken)
	if err != nil {
		return "", fmt.Errorf("authsvc.ValidateRefreshToken(): %w", err)
	}
	if claims.Typ != "refresh" {
		return "", fmt.Errorf("authsvc.ValidateRefreshToken(), claims.Typ != 'refresh': %w", err)
	}
	ctx, cancel := s.ctxFunc()
	defer cancel()
	user, err := s.userStore.FindByID(ctx, claims.Sub)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return "", fmt.Errorf("authsvc.ValidateRefreshToken: %w", err)
		}
		return "", fmt.Errorf("authsvc.ValidateRefreshToken: %w", err)
	}
	if user.RefreshToken != refreshToken {
		return "", fmt.Errorf("authsvc.ValidateRefreshToken: %w", ErrInvalidToken)
	}
	return user.ID, nil
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
		return models.TokenPair{}, fmt.Errorf("authsvc.issueTokens: %w", err)
	}
	refreshToken, err := signJWT(s.refreshSecret, refreshClaims)
	if err != nil {
		return models.TokenPair{}, fmt.Errorf("authsvc.issueTokens: %w", err)
	}
	ctx, cancel := s.ctxFunc()
	defer cancel()
	if err := s.userStore.UpdateTokens(ctx, user.ID, accessToken, refreshToken); err != nil {
		return models.TokenPair{}, fmt.Errorf("authsvc.issueTokens: %w", err)
	}
	return models.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
