package services

import (
	"order-bot-mgmt-svc/internal/services/menusvc"
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
	auth *lazy[authsvc.Svc]
	menu *lazy[menusvc.Svc]
}

func NewServices(authInit func() *authsvc.Svc, menuInit func() *menusvc.Svc) *Services {
	return &Services{
		auth: newLazy(authInit),
		menu: newLazy(menuInit),
	}
}

func (s *Services) Auth() *authsvc.Svc {
	if s == nil {
		return nil
	}
	return s.auth.Get()
}

func (s *Services) Menu() *menusvc.Svc {
	if s == nil {
		return nil
	}
	return s.menu.Get()
}
