package authsvc

import (
	"context"
	"errors"
	"fmt"
	"order-bot-mgmt-svc/internal/config"
	"order-bot-mgmt-svc/internal/infra/sqldb/pqsqldb"
	"order-bot-mgmt-svc/internal/models"
	"order-bot-mgmt-svc/internal/models/entities"
	"order-bot-mgmt-svc/internal/store"
	"order-bot-mgmt-svc/internal/util"
	"order-bot-mgmt-svc/internal/util/jwtutil"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Svc struct {
	db              *pqsqldb.DB
	ctxFunc         util.CtxFunc
	userStore       store.User
	accessSecret    []byte
	refreshSecret   []byte
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewSvc(db *pqsqldb.DB, ctxFunc util.CtxFunc, cfg config.Config, userStore store.User) *Svc {
	if userStore == nil || ctxFunc == nil {
		panic("authSvc.NewSvc(), userStore or ctxFunc is nil")
	}
	return &Svc{
		db:              db,
		ctxFunc:         ctxFunc,
		userStore:       userStore,
		accessSecret:    []byte(cfg.Auth.AccessSecret),
		refreshSecret:   []byte(cfg.Auth.RefreshSecret),
		accessTokenTTL:  cfg.Auth.AccessTokenTTL,
		refreshTokenTTL: cfg.Auth.RefreshTokenTTL,
	}
}

func (s *Svc) Signup(ctx context.Context, tx store.Tx, email, password string) (tokenPair models.TokenPair, userId string, err error) {
	if email == "" || password == "" {
		return models.TokenPair{}, "", fmt.Errorf("authsvc.Signup: %w", ErrInvalidCredentials)
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return models.TokenPair{}, "", fmt.Errorf("authsvc.Signup: %w", err)
	}
	newUser := entities.User{
		ID:           util.NewID(),
		Email:        email,
		PasswordHash: string(hash),
		AccessToken:  "",
		RefreshToken: "",
	}
	ctx, cancel := util.CallCtxFunc(ctx, s.ctxFunc)
	defer cancel()
	if err := s.userStore.Create(ctx, tx, newUser); err != nil {
		if errors.Is(err, store.ErrUserExists) {
			return models.TokenPair{}, "", fmt.Errorf("authsvc.Signup: %w", ErrUserExists)
		}
		return models.TokenPair{}, "", fmt.Errorf("authsvc.Signup: %w", err)
	}
	tokens, err := s.issueTokens(ctx, tx, newUser)
	if err != nil {
		return models.TokenPair{}, "", fmt.Errorf("authsvc.Signup: %w", err)
	}
	return tokens, newUser.ID, nil
}

func (s *Svc) Login(ctx context.Context, email, password string) (models.TokenPair, error) {
	ctx, cancel := util.CallCtxFunc(ctx, s.ctxFunc)
	defer cancel()
	user, errFindUsr := s.userStore.FindByEmail(ctx, nil, email)
	if errFindUsr != nil {
		if errors.Is(errFindUsr, store.ErrNotFound) {
			return models.TokenPair{}, fmt.Errorf("authsvc.Login: %w", ErrInvalidCredentials)
		}
		return models.TokenPair{}, fmt.Errorf("authsvc.Login: %w", errFindUsr)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return models.TokenPair{}, fmt.Errorf("authsvc.Login: %w", ErrInvalidCredentials)
	}
	tokens, err := s.issueTokens(ctx, nil, user)
	if err != nil {
		return models.TokenPair{}, fmt.Errorf("authsvc.Login: %w", err)
	}
	return tokens, nil
}

func (s *Svc) Logout(ctx context.Context, refreshToken string) error {
	userID, errValidation := s.ValidateRefreshToken(ctx, refreshToken)
	if errValidation != nil {
		return fmt.Errorf("authsvc.Logout: %w", errValidation)
	}
	ctx, cancel := util.CallCtxFunc(ctx, s.ctxFunc)
	defer cancel()
	if err := s.userStore.UpdateTokens(ctx, nil, userID, "", ""); err != nil {
		return fmt.Errorf("authsvc.Logout: %w", err)
	}
	return nil
}

func (s *Svc) ValidateAccessToken(_ context.Context, accessToken string) error {
	if accessToken == "" {
		return fmt.Errorf("authsvc.ValidateAccessToken(): accessToken is empty %w", jwtutil.ErrInvalidToken)
	}
	claims, err := jwtutil.ParseJWT(s.accessSecret, accessToken)
	if err != nil {
		return fmt.Errorf("authsvc.ValidateAccessToken(): %w", err)
	}
	if claims.Typ != "access" {
		return fmt.Errorf("authsvc.ValidateAccessToken(), claims.Typ != 'access': %w", jwtutil.ErrInvalidToken)
	}
	if claims.Exp > time.Now().UnixMilli() {
		return fmt.Errorf("authsvc.ValidateAccessToken(), accessToken is expired: %w", jwtutil.ErrExpiredToken)
	}
	return nil
}

func (s *Svc) ValidateRefreshToken(ctx context.Context, refreshToken string) (string, error) {
	if refreshToken == "" {
		return "", fmt.Errorf("authsvc.ValidateRefreshToken(): %w", jwtutil.ErrInvalidToken)
	}
	claims, err := jwtutil.ParseJWT(s.refreshSecret, refreshToken)
	if err != nil {
		return "", fmt.Errorf("authsvc.ValidateRefreshToken(): %w", err)
	}
	if claims.Typ != "refresh" {
		return "", fmt.Errorf("authsvc.ValidateRefreshToken(), claims.Typ != 'refresh': %w", err)
	}
	ctx, cancel := util.CallCtxFunc(ctx, s.ctxFunc)
	defer cancel()
	user, err := s.userStore.FindByID(ctx, nil, claims.Sub)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return "", fmt.Errorf("authsvc.ValidateRefreshToken: %w", err)
		}
		return "", fmt.Errorf("authsvc.ValidateRefreshToken: %w", err)
	}
	if user.RefreshToken != refreshToken {
		return "", fmt.Errorf("authsvc.ValidateRefreshToken: %w", jwtutil.ErrInvalidToken)
	}
	return user.ID, nil
}

func (s *Svc) issueTokens(ctx context.Context, tx store.Tx, user entities.User) (models.TokenPair, error) {
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
		return models.TokenPair{}, fmt.Errorf("authsvc.issueTokens: %w", err)
	}
	refreshToken, err := jwtutil.SignJWT(s.refreshSecret, refreshClaims)
	if err != nil {
		return models.TokenPair{}, fmt.Errorf("authsvc.issueTokens: %w", err)
	}
	ctx, cancel := util.CallCtxFunc(ctx, s.ctxFunc)
	defer cancel()
	if err := s.userStore.UpdateTokens(ctx, tx, user.ID, accessToken, refreshToken); err != nil {
		return models.TokenPair{}, fmt.Errorf("authsvc.issueTokens: %w", err)
	}
	return models.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
