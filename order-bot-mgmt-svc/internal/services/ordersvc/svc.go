package ordersvc

import (
	"context"
	"fmt"
	"log/slog"
	"order-bot-mgmt-svc/internal/infra/grpcsync"
	"order-bot-mgmt-svc/internal/models/entities"
	"order-bot-mgmt-svc/internal/store"
)

type Svc struct {
	db interface {
		WithTx(context.Context, func(context.Context, store.Tx) error) error
	}
	orderStore     store.Order
	orderItemStore store.OrderItem
	callbackClient *grpcsync.CallbackClient
}

func NewSvc(
	db interface {
		WithTx(context.Context, func(context.Context, store.Tx) error) error
	},
	orderStore store.Order,
	orderItemStore store.OrderItem,
	callbackClient *grpcsync.CallbackClient,
) *Svc {
	return &Svc{db: db, orderStore: orderStore, orderItemStore: orderItemStore, callbackClient: callbackClient}
}

func (s *Svc) ReceiveOrder(ctx context.Context, req grpcsync.SubmitOrderRequest) error {
	order := entities.Order{ID: req.OrderID, CartID: req.CartID, SessionID: req.SessionID, TotalScaled: req.TotalScaled}
	items := make([]entities.OrderItem, 0, len(req.Items))
	for _, item := range req.Items {
		items = append(items, entities.OrderItem{
			ID: item.ID, OrderID: req.OrderID, MenuItemID: item.MenuItemID, Name: item.Name,
			Quantity: item.Quantity, UnitPriceScaled: item.UnitPriceScaled, TotalPriceScaled: item.TotalPriceScaled,
		})
	}

	if err := s.db.WithTx(ctx, func(ctx context.Context, tx store.Tx) error {
		if err := s.orderStore.CreateOrder(ctx, tx, order); err != nil {
			return err
		}
		if err := s.orderItemStore.CreateOrderItems(ctx, tx, items); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return fmt.Errorf("ordersvc.Svc.ReceiveOrder: %w", err)
	}

	if err := s.callbackClient.UpdateOrderStatus(ctx, grpcsync.UpdateOrderStatusRequest{OrderID: req.OrderID, Status: "DONE"}); err != nil {
		slog.Error("ordersvc callback failed", "err", err)
		return err
	}
	return nil
}
