package server

import (
	"fmt"
	"net/http"
	"time"

	"order-bot-svc/internal/database"
)

type Server struct {
	port int

	db database.Service
}

func NewServer(port int, db database.Service) *http.Server {
	srv := &Server{
		port: port,

		db: db,
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
