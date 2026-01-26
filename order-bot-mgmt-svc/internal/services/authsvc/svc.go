package authsvc

import (
	"errors"
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
		AccessToken:  "",
		RefreshToken: "",
	}
	ctx, cancel := s.ctxFunc()
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
	ctx, cancel := s.ctxFunc()
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
	userID, err := s.ValidateRefreshToken(refreshToken)
	if err != nil {
		return err
	}
	ctx, cancel := s.ctxFunc()
	defer cancel()
	if err := s.userStore.UpdateTokens(ctx, userID, "", ""); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return ErrInvalidToken
		}
		return err
	}
	return nil
}

func (s *Svc) ValidateRefreshToken(refreshToken string) (string, error) {
	if refreshToken == "" {
		return "", ErrInvalidToken
	}
	claims, err := parseJWT(s.refreshSecret, refreshToken)
	if err != nil || claims.Typ != "refresh" {
		return "", ErrInvalidToken
	}
	ctx, cancel := s.ctxFunc()
	defer cancel()
	user, err := s.userStore.FindByID(ctx, claims.Sub)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return "", ErrInvalidToken
		}
		return "", err
	}
	if user.RefreshToken != refreshToken {
		return "", ErrInvalidToken
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
		return models.TokenPair{}, err
	}
	refreshToken, err := signJWT(s.refreshSecret, refreshClaims)
	if err != nil {
		return models.TokenPair{}, err
	}
	ctx, cancel := s.ctxFunc()
	defer cancel()
	if err := s.userStore.UpdateTokens(ctx, user.ID, accessToken, refreshToken); err != nil {
		return models.TokenPair{}, err
	}
	return models.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
