package httphdlrsold

import (
	"fmt"
	"net/http"
	"order-bot-mgmt-svc/internal/services/authsvc"
	"order-bot-mgmt-svc/internal/services/botsvc"
	"order-bot-mgmt-svc/internal/services/menusvc"
	"order-bot-mgmt-svc/internal/store"
	"time"

	"order-bot-mgmt-svc/internal/services"
)

type Server struct {
	port int

	idb      store.DB
	services *services.Services
}

func NewServer(port int, db store.DB, services *services.Services) *http.Server {
	srv := &Server{
		port: port,
		idb:  db,

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

func (s *Server) db() store.DB {
	return s.idb
}

func (s *Server) AuthService() *authsvc.Svc {
	return s.services.Auth
}

func (s *Server) MenuService() *menusvc.Svc {
	return s.services.Menu
}

func (s *Server) BotService() *botsvc.Svc {
	return s.services.Bot
}
