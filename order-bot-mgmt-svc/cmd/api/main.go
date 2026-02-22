package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"order-bot-mgmt-svc/internal/config"
	"order-bot-mgmt-svc/internal/infra/grpcsync"
	"order-bot-mgmt-svc/internal/infra/httphdlr/httpserver"
	"order-bot-mgmt-svc/internal/infra/sqldb"
	"order-bot-mgmt-svc/internal/infra/sqldb/orderbotmgmtsqldb"
	"order-bot-mgmt-svc/internal/services/authsvc"
	"order-bot-mgmt-svc/internal/services/botsvc"
	"order-bot-mgmt-svc/internal/services/menusvc"
	"order-bot-mgmt-svc/internal/services/ordersvc"
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

func newServices(db *sqldb.DB, orderBotDb *sqldb.DB, cfg config.Config) *services.Services {
	ctxFunc := util.NewCtxFunc(cfg.Others.QryCtxTimeout)
	return services.NewServices(
		func() *authsvc.Svc {
			return authsvc.NewSvc(db, ctxFunc, cfg, sqldb.NewUserStore(db))
		},
		func() *menusvc.Svc {
			menuStore := sqldb.NewMenuStore(db)
			menuItemStore := sqldb.NewMenuItemStore(db)
			publishedMenuStore := orderbotmgmtsqldb.NewPublishedMenuStore(orderBotDb)
			return menusvc.NewSvc(db, orderBotDb, ctxFunc, menuStore, menuItemStore, publishedMenuStore)
		},
		func() *botsvc.Svc {
			botStore := sqldb.NewBotStore(db)
			userBotStore := sqldb.NewUserBotStore(db)
			return botsvc.NewSvc(db, ctxFunc, cfg, botStore, userBotStore)
		},
		func() *ordersvc.Svc {
			orderStore := sqldb.NewOrderStore(orderBotDb)
			orderItemStore := sqldb.NewOrderItemStore(orderBotDb)
			callbackClient := grpcsync.NewCallbackClient(cfg.Others.OrderSvcGRPCAddr)
			return ordersvc.NewSvc(orderBotDb, orderStore, orderItemStore, callbackClient)
		},
	)
}

func main() {

	// Set up logger level
	slog.SetLogLoggerLevel(slog.LevelDebug)

	cfg := config.Load()
	port := cfg.App.Port
	db, err := sqldb.New(cfg.Db)
	if err != nil {
		log.Fatalf("failed to connect to database: \n%v", errutil.FormatErrChain(err))
	}
	orderBotDb, orderBotDbErr := sqldb.New(cfg.OrderBotDb)
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
	serviceContainer := newServices(db, orderBotDb, cfg)

	grpcServer := grpcsync.NewServer(cfg.Others.MgmtSvcGRPCAddr, serviceContainer.Order.Get())
	if err := grpcServer.Start(); err != nil {
		log.Fatalf("failed to start grpc server: %v", err)
	}

	server := httpserver.NewServer(
		port,
		db,
		serviceContainer,
	)
	addr := fmt.Sprintf("%s:%d", cfg.App.Address, cfg.App.Port)
	httpserver.Run(server, cfg.App.GinMode, addr)

	// // Create a done channel to signal when the shutdown is complete
	// done := make(chan bool, 1)

	// // Run graceful shutdown in a separate goroutine
	// go gracefulShutdown(server, done)

	// slog.Info("running http.ListAndServer...")
	// err = server.ListenAndServe()
	// if err != nil && !errors.Is(err, http.ErrServerClosed) {
	// 	panic(fmt.Sprintf("http server error: %s", err))
	// }

	// // Wait for the graceful shutdown to complete
	// <-done
	// log.Println("Graceful shutdown complete.")
}
