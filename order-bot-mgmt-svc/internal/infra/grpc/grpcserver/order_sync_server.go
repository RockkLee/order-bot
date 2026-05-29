package grpcserver

import (
	"context"
	"fmt"

	"order-bot-mgmt-svc/internal/infra/sqldb"
	"order-bot-mgmt-svc/internal/models/entities"
	"order-bot-mgmt-svc/internal/services/ordersvc"
	"order-bot-mgmt-svc/internal/store"

	orderbotv1pb "github.com/RockkLee/order-bot/goproto/orderbot/v1"
	orderbotmgmtsvcpb "github.com/RockkLee/order-bot/goproto/orderbot/v1/order_bot_mgmt_svc"
)

type OrderSyncServer struct {
	orderbotmgmtsvcpb.UnimplementedOrderSyncServiceServer
	orderSvc *ordersvc.Svc
	db       sqldb.Service
}

func NewOrderSyncServer(orderSvc *ordersvc.Svc, db sqldb.Service) *OrderSyncServer {
	return &OrderSyncServer{orderSvc: orderSvc, db: db}
}

func (s *OrderSyncServer) CreateOrder(ctx context.Context, req *orderbotv1pb.CreateOrderRequest) (*orderbotv1pb.CreateOrderResponse, error) {
	err := s.db.WithTx(ctx, func(ctx context.Context, tx store.Tx) error {
		order := entities.Order{
			ID:          req.OrderId,
			BotID:       req.BotId,
			CartID:      req.CartId,
			SessionID:   req.SessionId,
			TotalScaled: int(req.TotalScaled),
		}
		items := make([]entities.OrderItem, 0, len(req.Items))
		for _, item := range req.Items {
			items = append(items, entities.OrderItem{
				ID:               item.Id,
				OrderID:          req.OrderId,
				MenuItemID:       item.MenuItemId,
				Name:             item.Name,
				Quantity:         int(item.Quantity),
				UnitPriceScaled:  int(item.UnitPriceScaled),
				TotalPriceScaled: int(item.TotalPriceScaled),
			})
		}
		if err := s.orderSvc.CreateOrderWithItems(ctx, tx, order, items); err != nil {
			return fmt.Errorf("grpcserver.OrderSyncServer.CreateOrder: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &orderbotv1pb.CreateOrderResponse{Ok: true}, nil
}
