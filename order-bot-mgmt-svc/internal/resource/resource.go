package resource

import (
	"errors"
	"fmt"
	"order-bot-mgmt-svc/internal/config"
	"order-bot-mgmt-svc/internal/infra/sqldb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcClientConn struct {
	OrderBot *grpc.ClientConn
}

type Resource struct {
	DB             *sqldb.DB
	OrderBotDB     *sqldb.DB
	GrpcClientConn GrpcClientConn
}

func New(db, orderBotDB *sqldb.DB, grpcConn GrpcClientConn) *Resource {
	return &Resource{
		DB:             db,
		OrderBotDB:     orderBotDB,
		GrpcClientConn: grpcConn,
	}
}

func (r *Resource) Close() error {
	if r == nil {
		return nil
	}

	var errs []error
	if r.GrpcClientConn.OrderBot != nil {
		if err := r.GrpcClientConn.OrderBot.Close(); err != nil {
			errs = append(errs, fmt.Errorf("resource.Resource.Close grpc order bot conn: %w", err))
		}
	}
	if r.DB != nil {
		if err := r.DB.Close(); err != nil {
			errs = append(errs, fmt.Errorf("resource.Resource.Close db: %w", err))
		}
	}
	if r.OrderBotDB != nil {
		if err := r.OrderBotDB.Close(); err != nil {
			errs = append(errs, fmt.Errorf("resource.Resource.Close order bot db: %w", err))
		}
	}
	return errors.Join(errs...)
}

func NewOrderBotGRPCConn(cfg config.Grpc) (*grpc.ClientConn, error) {
	target := fmt.Sprintf("%s:%d", cfg.Address, cfg.Port)
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("resource.NewOrderBotGRPCConn: %w", err)
	}
	return conn, nil
}
