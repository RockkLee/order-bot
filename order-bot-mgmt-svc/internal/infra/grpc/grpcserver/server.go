package grpcserver

import (
	"fmt"
	"net"

	"order-bot-mgmt-svc/internal/infra/sqldb"
	"order-bot-mgmt-svc/internal/services/ordersvc"

	orderbotmgmtsvcpb "github.com/RockkLee/order-bot/goproto/orderbot/v1/order_bot_mgmt_svc"
	"google.golang.org/grpc"
)

func NewServer(orderSvc *ordersvc.Svc, db sqldb.Service) *grpc.Server {
	server := grpc.NewServer()
	orderbotmgmtsvcpb.RegisterOrderSyncServiceServer(server, NewOrderSyncServer(orderSvc, db))
	return server
}

func NewListeningServer(addr string, orderSvc *ordersvc.Svc, db sqldb.Service) (*grpc.Server, net.Listener, error) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, nil, fmt.Errorf("grpcserver.NewListeningServer listen %s: %w", addr, err)
	}

	return NewServer(orderSvc, db), lis, nil
}
