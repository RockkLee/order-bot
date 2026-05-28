package grpcserver

import (
	"context"
	"fmt"
	"order-bot-mgmt-svc/internal/infra/sqldb"
	"order-bot-mgmt-svc/internal/models/entities"
	"order-bot-mgmt-svc/internal/services/ordersvc"
	"order-bot-mgmt-svc/internal/store"
)

type OrderSyncServer struct {
	orderSvc *ordersvc.Svc
	txSvc    sqldb.Service
}

func NewOrderSyncServer(orderSvc *ordersvc.Svc, txSvc sqldb.Service) *OrderSyncServer {
	return &OrderSyncServer{orderSvc: orderSvc, txSvc: txSvc}
}

type CreateOrderRequest struct {
	OrderID     string
	BotID       string
	CartID      string
	SessionID   string
	TotalScaled int
	Items       []CreateOrderItem
}

type CreateOrderItem struct {
	ID               string
	MenuItemID       string
	Name             string
	Quantity         int
	UnitPriceScaled  int
	TotalPriceScaled int
}

func (s *OrderSyncServer) CreateOrder(ctx context.Context, req CreateOrderRequest) error {
	return s.txSvc.WithTx(ctx, func(ctx context.Context, tx store.Tx) error {
		order := entities.Order{ID: req.OrderID, BotID: req.BotID, CartID: req.CartID, SessionID: req.SessionID, TotalScaled: req.TotalScaled}
		items := make([]entities.OrderItem, 0, len(req.Items))
		for _, item := range req.Items {
			items = append(items, entities.OrderItem{ID: item.ID, OrderID: req.OrderID, MenuItemID: item.MenuItemID, Name: item.Name, Quantity: item.Quantity, UnitPriceScaled: item.UnitPriceScaled, TotalPriceScaled: item.TotalPriceScaled})
		}
		if err := s.orderSvc.CreateOrderWithItems(ctx, tx, order, items); err != nil {
			return fmt.Errorf("grpcserver.OrderSyncServer.CreateOrder: %w", err)
		}
		return nil
	})
}
