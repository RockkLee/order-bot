package ordersvc

import (
	"context"
	"fmt"
	"order-bot-mgmt-svc/internal/models/entities"
	"order-bot-mgmt-svc/internal/store"
	"order-bot-mgmt-svc/internal/util"
)

type OrderWithItems struct {
	Order entities.Order
	Items []entities.OrderItem
}

type Svc struct {
	orderStore     store.Order
	orderItemStore store.OrderItem
	ctxFunc        util.CtxFunc
}

func NewSvc(ctxFunc util.CtxFunc, orderStore store.Order, orderItemStore store.OrderItem) *Svc {
	if orderStore == nil || orderItemStore == nil {
		panic("ordersvc.NewSvc(), orderStore or orderItemStore is nil")
	}
	return &Svc{orderStore: orderStore, orderItemStore: orderItemStore, ctxFunc: ctxFunc}
}

func (s *Svc) GetOrdersWithItems(ctx context.Context, botId string) ([]OrderWithItems, error) {
	ctx, cancel := util.CallCtxFunc(ctx, s.ctxFunc)
	defer cancel()

	orders, err := s.orderStore.FindByBotID(ctx, nil, botId)
	if err != nil {
		return nil, fmt.Errorf("ordersvc.GetOrdersWithItems: %w", err)
	}
	orderIDs := make([]string, 0, len(orders))
	for _, order := range orders {
		orderIDs = append(orderIDs, order.ID)
	}
	items, err := s.orderItemStore.FindByOrderIDs(ctx, orderIDs)
	if err != nil {
		return nil, fmt.Errorf("ordersvc.GetOrdersWithItems: %w", err)
	}
	itemsByOrderID := make(map[string][]entities.OrderItem, len(orderIDs))
	for _, item := range items {
		itemsByOrderID[item.OrderID] = append(itemsByOrderID[item.OrderID], item)
	}
	result := make([]OrderWithItems, 0, len(orders))
	for _, order := range orders {
		result = append(result, OrderWithItems{Order: order, Items: itemsByOrderID[order.ID]})
	}
	return result, nil
}
