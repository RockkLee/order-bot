package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"order-bot-mgmt-svc/internal/infra/httphdlrs"
	"order-bot-mgmt-svc/internal/infra/postgres"
	postgresuser "order-bot-mgmt-svc/internal/infra/postgres/user"
	"os"
	"os/signal"
	"strconv"
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

func newServices() *services.Services {
	return services.NewServices(
		func() *services.AuthService {
			db := postgres.New()
			return services.NewAuthService(postgresuser.NewStore(db.Conn()))
		},
		func() *services.MenuService {
			return services.NewMenuService()
		},
	)
}

func main() {

	port, _ := strconv.Atoi(os.Getenv("PORT"))
	server := httphdlrs.NewServer(
		port,
		postgres.New,
		func() *services.Services { return newServices() },
	)

	// Create a done channel to signal when the shutdown is complete
	done := make(chan bool, 1)

	// Run graceful shutdown in a separate goroutine
	go gracefulShutdown(server, done)

	err := server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(fmt.Sprintf("http server error: %s", err))
	}

	// Wait for the graceful shutdown to complete
	<-done
	log.Println("Graceful shutdown complete.")
}
