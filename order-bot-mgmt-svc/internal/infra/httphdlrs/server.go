package httphdlrs

import (
	"fmt"
	"net/http"
	"order-bot-mgmt-svc/internal/infra/postgres"
	"sync"
	"time"

	"order-bot-mgmt-svc/internal/services"
)

type Server struct {
	port int

	dbInit       func() postgres.Service
	dbOnce       sync.Once
	db           postgres.Service
	servicesInit func() *services.Services
	servicesOnce sync.Once
	services     *services.Services
}

func NewServer(port int, dbInit func() postgres.Service, servicesInit func() *services.Services) *http.
	Server {
	srv := &Server{
		port: port,

		dbInit:       dbInit,
		servicesInit: servicesInit,
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
	s.dbOnce.Do(func() {
		s.db = s.dbInit()
	})

	return s.db
}

func (s *Server) servicesContainer() *services.Services {
	s.servicesOnce.Do(func() {
		s.services = s.servicesInit()
	})

	return s.services
}

func (s *Server) authService() *services.AuthService {
	return s.servicesContainer().Auth()
}
