package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"order-bot-mgmt-svc/internal/config"
	"order-bot-mgmt-svc/internal/infra/httphdlrs/httpserver"
	"order-bot-mgmt-svc/internal/infra/postgres"
	postgresuser "order-bot-mgmt-svc/internal/infra/postgres/user"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"order-bot-mgmt-svc/internal/services"
)

func gracefulShutdown(apiServer *http.Server, done chan bool) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Listen for the interrupt signal.
	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")
	stop() // Allow Ctrl+C to force shutdown

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := apiServer.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server exiting")

	// Notify the main goroutine that the shutdown is complete
	done <- true
}

func newServices(db *postgres.DB, cfg config.Config) *services.Services {
	return services.NewServices(
		func() *services.AuthService {
			return services.NewAuthService(postgresuser.NewStore(db.Conn()), cfg.Auth)
		},
		func() *services.MenuService {
			return services.NewMenuService()
		},
	)
}

func main() {

	cfg := config.Load()
	port := cfg.App.Port
	db, err := postgres.New(cfg.DB)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("failed to close database: %v", err)
		}
	}()
	serviceContainer := newServices(db, cfg)

	server := httpserver.NewServer(
		port,
		db,
		serviceContainer,
	)

	// Create a done channel to signal when the shutdown is complete
	done := make(chan bool, 1)

	// Run graceful shutdown in a separate goroutine
	go gracefulShutdown(server, done)

	err = server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(fmt.Sprintf("http server error: %s", err))
	}

	// Wait for the graceful shutdown to complete
	<-done
	log.Println("Graceful shutdown complete.")
}
