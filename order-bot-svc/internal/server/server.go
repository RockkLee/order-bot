package server

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"order-bot-svc/internal/database"
)

type Server struct {
	port int

	dbInit func() database.Service
	dbOnce sync.Once
	db     database.Service
}

func NewServer(port int, dbInit func() database.Service) *http.Server {
	srv := &Server{
		port: port,

		dbInit: dbInit,
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

func (s *Server) dbService() database.Service {
	s.dbOnce.Do(func() {
		s.db = s.dbInit()
	})

	return s.db
}
