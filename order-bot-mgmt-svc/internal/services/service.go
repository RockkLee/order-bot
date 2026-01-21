package services

import (
	"errors"
	"sync"
)

var (
	ErrUserExists         = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
)

type lazy[T any] struct {
	once func()
	val  *T
}

func newLazy[T any](init func() *T) *lazy[T] {
	l := &lazy[T]{}
	l.once = sync.OnceFunc(func() {
		if init != nil {
			l.val = init()
		}
	})
	return l
}

func (l *lazy[T]) Get() *T {
	if l == nil {
		return nil
	}
	l.once()
	return l.val
}

type Services struct {
	auth *lazy[AuthService]
	menu *lazy[MenuService]
}

func NewServices(authInit func() *AuthService, menuInit func() *MenuService) *Services {
	return &Services{
		auth: newLazy(authInit),
		menu: newLazy(menuInit),
	}
}

func (s *Services) Auth() *AuthService {
	if s == nil {
		return nil
	}
	return s.auth.Get()
}

func (s *Services) Menu() *MenuService {
	if s == nil {
		return nil
	}
	return s.menu.Get()
}
