package handlers

import (
	"fmt"
	"net/http"
	"time"

	"order-bot-mgmt-svc/internal/repository"
	"order-bot-mgmt-svc/internal/services"
)

type Server struct {
	port int

	db   repository.Service
	auth *services.Service
}

func NewServer(port int, db repository.Service, auth *services.Service) *http.Server {
	srv := &Server{
		port: port,

		db:   db,
		auth: auth,
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
