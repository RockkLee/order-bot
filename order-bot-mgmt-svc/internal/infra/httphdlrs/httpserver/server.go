package httpserver

import (
	"fmt"
	"net/http"
	"order-bot-mgmt-svc/internal/infra/sqldb/pqsqldb"
	"order-bot-mgmt-svc/internal/services/authsvc"
	"order-bot-mgmt-svc/internal/services/botsvc"
	"order-bot-mgmt-svc/internal/services/menusvc"
	"time"

	"order-bot-mgmt-svc/internal/services"
)

type Server struct {
	port int

	db       pqsqldb.Service
	services *services.Services
}

func NewServer(port int, db pqsqldb.Service, services *services.Services) *http.Server {
	srv := &Server{
		port: port,
		db:   db,

		services: services,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", srv.port),
		Handler:      srv.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}

func (s *Server) dbService() pqsqldb.Service {
	return s.db
}

func (s *Server) AuthService() *authsvc.Svc {
	return s.services.Auth.Get()
}

func (s *Server) MenuService() *menusvc.Svc {
	return s.services.Menu.Get()
}

func (s *Server) BotService() *botsvc.Svc {
	return s.services.Bot.Get()
}
