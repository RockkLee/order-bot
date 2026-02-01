package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"order-bot-mgmt-svc/internal/config"
	"order-bot-mgmt-svc/internal/infra/httphdlrs/httpserver"
	"order-bot-mgmt-svc/internal/infra/sqldb"
	"order-bot-mgmt-svc/internal/infra/sqldb/pqsqldb"
	"order-bot-mgmt-svc/internal/services/authsvc"
	"order-bot-mgmt-svc/internal/services/botsvc"
	"order-bot-mgmt-svc/internal/services/menusvc"
	"order-bot-mgmt-svc/internal/util"
	"order-bot-mgmt-svc/internal/util/errutil"
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

func newServices(db *pqsqldb.DB, cfg config.Config) *services.Services {
	ctxFunc := util.NewCtxFunc(cfg.Others.QryCtxTimeout)
	return services.NewServices(
		func() *authsvc.Svc {
			return authsvc.NewSvc(db, ctxFunc, cfg, sqldb.NewUserStore(db))
		},
		func() *menusvc.Svc {
			menuStore := sqldb.NewMenuStore(db)
			menuItemStore := sqldb.NewMenuItemStore(db)
			return menusvc.NewSvc(db, ctxFunc, menuStore, menuItemStore)
		},
		func() *botsvc.Svc {
			botStore := sqldb.NewBotStore(db)
			userBotStore := sqldb.NewUserBotStore(db)
			return botsvc.NewSvc(botStore, userBotStore, db, ctxFunc)
		},
	)
}

func main() {

	// Set up logger level
	slog.SetLogLoggerLevel(slog.LevelDebug)

	cfg := config.Load()
	port := cfg.App.Port
	db, err := pqsqldb.New(cfg.Db)
	if err != nil {
		log.Fatalf("failed to connect to database: \n%v", errutil.FormatErrChain(err))
	}
	orderBotDb, orderBotDbErr := pqsqldb.New(cfg.OrderBotDb)
	if orderBotDbErr != nil {
		log.Fatalf("failed to connect to order-bot database: \n%v", orderBotDbErr)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("failed to close database: \n%v", errutil.FormatErrChain(err))
		}
		if err := orderBotDb.Close(); err != nil {
			log.Printf("failed to close order-bot database: \n%v", errutil.FormatErrChain(err))
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

	slog.Info("running http.ListAndServer...")
	err = server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(fmt.Sprintf("http server error: %s", err))
	}

	// Wait for the graceful shutdown to complete
	<-done
	log.Println("Graceful shutdown complete.")
}
