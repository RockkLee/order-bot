package services

import (
	"order-bot-mgmt-svc/internal/services/botsvc"
	"order-bot-mgmt-svc/internal/services/menusvc"
	"order-bot-mgmt-svc/internal/services/ordersvc"
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
	Auth  *lazy[authsvc.Svc]
	Menu  *lazy[menusvc.Svc]
	Bot   *lazy[botsvc.Svc]
	Order *lazy[ordersvc.Svc]
}

func NewServices(authInit func() *authsvc.Svc, menuInit func() *menusvc.Svc, botInit func() *botsvc.Svc, orderInit func() *ordersvc.Svc) *Services {
	return &Services{
		Auth:  newLazy(authInit),
		Menu:  newLazy(menuInit),
		Bot:   newLazy(botInit),
		Order: newLazy(orderInit),
	}
}
