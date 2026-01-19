package handlers

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"order-bot-mgmt-svc/internal/repository"
	"order-bot-mgmt-svc/internal/services"
)

type Server struct {
	port int

	dbInit   func() repository.Service
	authInit func() *services.Service
	dbOnce   sync.Once
	authOnce sync.Once
	db       repository.Service
	auth     *services.Service
}

func NewServer(port int, dbInit func() repository.Service, authInit func() *services.Service) *http.Server {
	srv := &Server{
		port: port,

		dbInit:   dbInit,
		authInit: authInit,
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

func (s *Server) dbService() repository.Service {
	s.dbOnce.Do(func() {
		s.db = s.dbInit()
	})

	return s.db
}

func (s *Server) authService() *services.Service {
	s.authOnce.Do(func() {
		s.auth = s.authInit()
	})

	return s.auth
}
