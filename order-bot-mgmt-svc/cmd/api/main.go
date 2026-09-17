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
	"order-bot-mgmt-svc/internal/infra/httphdlr"
	"order-bot-mgmt-svc/internal/infra/sqldb"
	"order-bot-mgmt-svc/internal/infra/sqldb/orderbotsqldb"
	"order-bot-mgmt-svc/internal/services/authsvc"
	"order-bot-mgmt-svc/internal/services/botsvc"
	"order-bot-mgmt-svc/internal/services/menusvc"
	"order-bot-mgmt-svc/internal/services/ordersvc"
	"order-bot-mgmt-svc/internal/util/errutil"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"order-bot-mgmt-svc/internal/services"
)

func newServices(rsrc *config.Resource, cfg config.Config) *services.Services {
	db, ok := rsrc.DB.(*sqldb.DB)
	if !ok {
		panic("newServices: resource DB must be *sqldb.DB for GORM-backed stores")
	}
	orderBotDB, ok := rsrc.OrderBotDB.(*sqldb.DB)
	if !ok {
		panic("newServices: resource order bot DB must be *sqldb.DB for GORM-backed stores")
	}

	auth := authsvc.NewSvc(rsrc, cfg, sqldb.NewUserStore(db))

	menuStore := sqldb.NewMenuStore(db)
	menuItemStore := sqldb.NewMenuItemStore(db)
	publishedMenuStore := orderbotsqldb.NewPublishedMenuStore(orderBotDB)
	menu := menusvc.NewSvc(rsrc, menuStore, menuItemStore, publishedMenuStore)

	botStore := sqldb.NewBotStore(db)
	userBotStore := sqldb.NewUserBotStore(db)
	bot := botsvc.NewSvc(rsrc, cfg, botStore, userBotStore)

	orderStore := sqldb.NewOrderStore(orderBotDB)
	orderItemStore := sqldb.NewOrderItemStore(orderBotDB)
	order := ordersvc.NewSvc(orderStore, orderItemStore)

	return services.NewServices(
		auth,
		menu,
		bot,
		order,
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
	orderBotGrpcConn, err := config.NewOrderBotGRPCConn(cfg.OrderBotGrpcClient)
	if err != nil {
		log.Fatalf("failed to create order-bot grpc client connection: \n%v", errutil.FormatErrChain(err))
	}
	rsrc := config.New(db, orderBotDb, config.GrpcClientConn{OrderBot: orderBotGrpcConn})
	defer func() {
		if err := rsrc.Close(); err != nil {
			log.Printf("failed to close resources: \n%v", errutil.FormatErrChain(err))
		}
	}()
	svcs := newServices(rsrc, cfg)

	// Build the Gin-backed HTTP server explicitly so main owns startup and shutdown.
	httpServContainer := httphdlr.NewServerContainer(cfg.App.Port, rsrc.DB, svcs)
	httpAddr := fmt.Sprintf("%s:%d", cfg.App.Address, cfg.App.Port)
	httpSrv := httphdlr.NewHTTPServer(httpServContainer, cfg.App.GinMode, httpAddr)

	// Create the gRPC listener before starting goroutines so bind failures surface immediately.
	grpcAddr := fmt.Sprintf("%s:%d", cfg.GrpcServer.Address, cfg.GrpcServer.Port)
	grpcSrv, grpcLis, err := grpcserver.NewListeningServer(grpcAddr, svcs.Order, rsrc.OrderBotDB)
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

	// Waiting here until either the process receives a shutdown signal or one of the servers fails.
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
