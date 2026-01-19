package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"os"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserExists         = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
)

type User struct {
	ID           string
	Email        string
	PasswordHash string
}

type RefreshRecord struct {
	UserID    string
	ExpiresAt time.Time
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Service struct {
	mu              sync.Mutex
	usersByEmail    map[string]User
	refreshTokens   map[string]RefreshRecord
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
		usersByEmail:    make(map[string]User),
		refreshTokens:   make(map[string]RefreshRecord),
		accessSecret:    []byte(accessSecret),
		refreshSecret:   []byte(refreshSecret),
		accessTokenTTL:  parseDurationEnv("JWT_ACCESS_TTL", 15*time.Minute),
		refreshTokenTTL: parseDurationEnv("JWT_REFRESH_TTL", 7*24*time.Hour),
	}
}

func (s *Service) Signup(email, password string) (TokenPair, error) {
	if email == "" || password == "" {
		return TokenPair{}, ErrInvalidCredentials
	}
	s.mu.Lock()
	if _, exists := s.usersByEmail[email]; exists {
		s.mu.Unlock()
		return TokenPair{}, ErrUserExists
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		s.mu.Unlock()
		return TokenPair{}, err
	}
	user := User{
		ID:           newID(),
		Email:        email,
		PasswordHash: string(hash),
	}
	s.usersByEmail[email] = user
	s.mu.Unlock()
	return s.issueTokens(user)
}

func (s *Service) Login(email, password string) (TokenPair, error) {
	s.mu.Lock()
	user, exists := s.usersByEmail[email]
	s.mu.Unlock()
	if !exists {
		return TokenPair{}, ErrInvalidCredentials
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return TokenPair{}, ErrInvalidCredentials
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

func (s *Service) issueTokens(user User) (TokenPair, error) {
	now := time.Now()
	accessClaims := Claims{
		Sub:   user.ID,
		Email: user.Email,
		Exp:   now.Add(s.accessTokenTTL).Unix(),
		Iat:   now.Unix(),
		Typ:   "access",
	}
	refreshClaims := Claims{
		Sub:   user.ID,
		Email: user.Email,
		Exp:   now.Add(s.refreshTokenTTL).Unix(),
		Iat:   now.Unix(),
		Typ:   "refresh",
	}
	accessToken, err := signJWT(s.accessSecret, accessClaims)
	if err != nil {
		return TokenPair{}, err
	}
	refreshToken, err := signJWT(s.refreshSecret, refreshClaims)
	if err != nil {
		return TokenPair{}, err
	}
	s.mu.Lock()
	s.refreshTokens[refreshToken] = RefreshRecord{
		UserID:    user.ID,
		ExpiresAt: now.Add(s.refreshTokenTTL),
	}
	s.mu.Unlock()
	return TokenPair{
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
