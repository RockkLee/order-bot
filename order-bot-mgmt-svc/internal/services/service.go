package services

import (
	"sync"

	"order-bot-mgmt-svc/internal/services/authsvc"
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
	auth *lazy[authsvc.Auth]
	menu *lazy[MenuService]
}

func NewServices(authInit func() *authsvc.Auth, menuInit func() *MenuService) *Services {
	return &Services{
		auth: newLazy(authInit),
		menu: newLazy(menuInit),
	}
}

func (s *Services) Auth() *authsvc.Auth {
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
