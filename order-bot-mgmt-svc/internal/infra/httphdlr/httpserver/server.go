package httpserver

import (
	"context"
	"order-bot-mgmt-svc/internal/infra/sqldb"
	"order-bot-mgmt-svc/internal/services"
	"order-bot-mgmt-svc/internal/services/authsvc"
	"order-bot-mgmt-svc/internal/services/botsvc"
	"order-bot-mgmt-svc/internal/services/menusvc"
	"order-bot-mgmt-svc/internal/services/ordersvc"
	"order-bot-mgmt-svc/internal/store"
)

type Server struct {
	port int

	db       sqldb.Service
	services *services.Services
}

func NewServer(port int, db sqldb.Service, services *services.Services) *Server {
	return &Server{port: port, db: db, services: services}
}

func (s *Server) dbService() sqldb.Service { return s.db }
func (s *Server) WithTx(ctx context.Context, fn func(ctx context.Context, tx store.Tx) error) error {
	return s.db.WithTx(ctx, fn)
}
func (s *Server) GetWithTx(ctx context.Context, fn func(ctx context.Context, tx store.Tx) (any, error)) (any, error) {
	return s.db.GetWithTx(ctx, fn)
}
func (s *Server) AuthService() *authsvc.Svc   { return s.services.Auth.Get() }
func (s *Server) MenuService() *menusvc.Svc   { return s.services.Menu.Get() }
func (s *Server) BotService() *botsvc.Svc     { return s.services.Bot.Get() }
func (s *Server) OrderService() *ordersvc.Svc { return s.services.Order.Get() }
