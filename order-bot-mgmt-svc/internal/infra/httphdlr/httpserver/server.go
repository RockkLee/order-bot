package httpserver

import (
	"context"
	"net/http"
	"order-bot-mgmt-svc/internal/infra/sqldb"
	"order-bot-mgmt-svc/internal/services"
	"order-bot-mgmt-svc/internal/services/authsvc"
	"order-bot-mgmt-svc/internal/services/botsvc"
	"order-bot-mgmt-svc/internal/services/menusvc"
	"order-bot-mgmt-svc/internal/services/ordersvc"
	"order-bot-mgmt-svc/internal/store"
)

type ServerContainer struct {
	port int

	db       sqldb.Service
	services *services.Services
}

func NewServerContainer(port int, db sqldb.Service, services *services.Services) *ServerContainer {
	return &ServerContainer{port: port, db: db, services: services}
}

func (s *ServerContainer) dbService() sqldb.Service { return s.db }
func (s *ServerContainer) WithTx(ctx context.Context, fn func(ctx context.Context, tx store.Tx) error) error {
	return s.db.WithTx(ctx, fn)
}
func (s *ServerContainer) GetWithTx(ctx context.Context, fn func(ctx context.Context, tx store.Tx) (any, error)) (any, error) {
	return s.db.GetWithTx(ctx, fn)
}
func (s *ServerContainer) AuthService() *authsvc.Svc { return s.services.Auth.Get() }
func (s *ServerContainer) MenuService() *menusvc.Svc { return s.services.Menu.Get() }
func (s *ServerContainer) BotService() *botsvc.Svc   { return s.services.Bot.Get() }
func (s *ServerContainer) OrderService() *ordersvc.Svc {
	return s.services.Order.Get()
}

func NewHTTPServer(s *ServerContainer, ginMode string, addr string) *http.Server {
	return &http.Server{
		Addr:    addr,
		Handler: NewRouters(s, ginMode),
	}
}
