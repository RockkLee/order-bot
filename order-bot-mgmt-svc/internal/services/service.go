package services

import (
	"errors"
	"os"
	"sync"
	"time"
)

var (
	ErrUserExists         = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
)

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

type lazy[T any] struct {
	once sync.Once
	init func() *T
	val  *T
}

func (l *lazy[T]) Get() *T {
	if l == nil {
		return nil
	}
	l.once.Do(func() {
		if l.init != nil {
			l.val = l.init()
		}
	})
	return l.val
}

type Services struct {
	auth lazy[AuthService]
	menu lazy[MenuService]
}

func NewServices(authInit func() *AuthService, menuInit func() *MenuService) *Services {
	return &Services{
		auth: lazy[AuthService]{
			init: authInit,
		},
		menu: lazy[MenuService]{
			init: menuInit,
		},
	}
}

func (s *Services) Auth() *AuthService {
	return s.auth.Get()
}

func (s *Services) Menu() *MenuService {
	return s.menu.Get()
}
