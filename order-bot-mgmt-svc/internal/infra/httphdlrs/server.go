package httphdlrs

import (
	"fmt"
	"net/http"
	"order-bot-mgmt-svc/internal/infra/postgres"
	"time"

	"order-bot-mgmt-svc/internal/services"
)

type Server struct {
	port int

	db      postgres.Service
	authSvc *services.AuthService
	menuSvc *services.MenuService
}

func NewServer(port int, db postgres.Service, authService *services.AuthService, menuService *services.MenuService) *http.Server {
	srv := &Server{
		port: port,
		db:   db,

		authSvc: authService,
		menuSvc: menuService,
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

func (s *Server) authService() *services.AuthService {
	return s.authSvc
}

func (s *Server) menuService() *services.MenuService {
	return s.menuSvc
}
