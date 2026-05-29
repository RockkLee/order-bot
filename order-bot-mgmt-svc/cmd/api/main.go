package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"order-bot-mgmt-svc/internal/config"
	"order-bot-mgmt-svc/internal/infra/grpc/grpcserver"
	"order-bot-mgmt-svc/internal/infra/httphdlr/httpserver"
	"order-bot-mgmt-svc/internal/infra/sqldb"
	"order-bot-mgmt-svc/internal/infra/sqldb/orderbotsqldb"
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

func newServices(db *sqldb.DB, orderBotDb *sqldb.DB, cfg config.Config) *services.Services {
	ctxFunc := util.NewCtxFunc(cfg.Others.QryCtxTimeout)
	return services.NewServices(
		func() *authsvc.Svc {
			return authsvc.NewSvc(db, ctxFunc, cfg, sqldb.NewUserStore(db))
		},
		func() *menusvc.Svc {
			menuStore := sqldb.NewMenuStore(db)
			menuItemStore := sqldb.NewMenuItemStore(db)
			publishedMenuStore := orderbotsqldb.NewPublishedMenuStore(orderBotDb)
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
			return ordersvc.NewSvc(ctxFunc, orderStore, orderItemStore)
		},
	)
}

func main() {

	// Set up logger level
	slog.SetLogLoggerLevel(slog.LevelInfo)

	cfg := config.Load()
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

	// Build the Gin-backed HTTP server explicitly so main owns startup and shutdown.
	httpServContainer := httpserver.NewServerContainer(cfg.App.Port, db, serviceContainer)
	httpAddr := fmt.Sprintf("%s:%d", cfg.App.Address, cfg.App.Port)
	httpSrv := httpserver.NewHTTPServer(httpServContainer, cfg.App.GinMode, httpAddr)

	// Create the gRPC listener before starting goroutines so bind failures surface immediately.
	grpcAddr := fmt.Sprintf("%s:%d", cfg.Grpc.Address, cfg.Grpc.Port)
	grpcSrv, grpcLis, err := grpcserver.NewListeningServer(grpcAddr, serviceContainer.Order.Get(), orderBotDb)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = grpcLis.Close() }()

	// Cancel on SIGINT/SIGTERM so both servers can drain gracefully.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 2)

	// HTTP and gRPC both block while serving, so they run concurrently in separate goroutines.
	go func() {
		if err := httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- fmt.Errorf("http server: %w", err)
		}
	}()
	go func() {
		if err := grpcSrv.Serve(grpcLis); err != nil {
			errCh <- fmt.Errorf("grpc server: %w", err)
		}
	}()

	// Stop when either the process receives a shutdown signal or one of the servers fails.
	var runErr error
	select {
	case <-ctx.Done():
		slog.Info("shutdown signal received")
	case runErr = <-errCh:
		stop()
	}

	shutdownHttpCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Stop accepting new HTTP requests and wait for in-flight work to finish.
	if err := httpSrv.Shutdown(shutdownHttpCtx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("http shutdown failed", "error", err)
		if runErr == nil {
			runErr = err
		}
	}

	// No goroutine sends a value on grpcStopped. Closing the channel is the signal
	// that GracefulStop has finished, which lets the main goroutine continue.
	grpcStopped := make(chan struct{})
	go func() {
		// GracefulStop blocks this goroutine until active RPCs complete.
		grpcSrv.GracefulStop()
		close(grpcStopped)
	}()

	select {
	// This receive blocks the main goroutine until grpcStopped is closed.
	case <-grpcStopped:
		// The goroutine above finished GracefulStop and unblocked this select by closing the channel.
	case <-shutdownHttpCtx.Done():
		// Fall back to an immediate stop if graceful shutdown exceeds the timeout.
		grpcSrv.Stop()
	}

	if runErr != nil {
		log.Fatal(runErr)
	}
}
