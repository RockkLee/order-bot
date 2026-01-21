package httpserver

import (
	"fmt"
	"net/http"
	"order-bot-mgmt-svc/internal/infra/postgres"
	"order-bot-mgmt-svc/internal/services/authsvc"
	"order-bot-mgmt-svc/internal/services/menusvc"
	"time"

	"order-bot-mgmt-svc/internal/services"
)

type Server struct {
	port int

	db       postgres.Service
	services *services.Services
}

func NewServer(port int, db postgres.Service, services *services.Services) *http.Server {
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

func (s *Server) dbService() postgres.Service {
	return s.db
}

func (s *Server) AuthService() *authsvc.Svc {
	if s.services == nil {
		return nil
	}
	return s.services.Auth()
}

func (s *Server) menuService() *menusvc.Svc {
	if s.services == nil {
		return nil
	}
	return s.services.Menu()
}
