package grpcserver

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"order-bot-mgmt-svc/internal/models/entities"
	"order-bot-mgmt-svc/internal/services/ordersvc"
	"order-bot-mgmt-svc/internal/store"
	"order-bot-mgmt-svc/internal/util"

	orderbotv1pb "github.com/RockkLee/order-bot/goproto/orderbot/v1"
)

type fakeOrderStore struct {
	inserted entities.Order
}

func (f *fakeOrderStore) FindByBotID(ctx context.Context, tx store.Tx, botId string) ([]entities.Order, error) {
	return nil, nil
}

func (f *fakeOrderStore) Insert(ctx context.Context, tx store.Tx, order entities.Order) error {
	f.inserted = order
	return nil
}

type fakeOrderItemStore struct {
	inserted []entities.OrderItem
}

func (f *fakeOrderItemStore) FindByOrderIDs(ctx context.Context, orderIDs []string) ([]entities.OrderItem, error) {
	return nil, nil
}

func (f *fakeOrderItemStore) InsertMany(ctx context.Context, tx store.Tx, items []entities.OrderItem) error {
	f.inserted = append([]entities.OrderItem(nil), items...)
	return nil
}

type fakeDb struct{}

func (s *fakeDb) Health() (map[string]string, error) { return nil, nil }
func (s *fakeDb) Close() error                       { return nil }
func (s *fakeDb) Conn() *sql.DB                      { return nil }
func (s *fakeDb) GetWithTx(ctx context.Context, fn func(ctx context.Context, tx store.Tx) (any, error)) (any, error) {
	return fn(ctx, struct{}{})
}
func (s *fakeDb) WithTx(ctx context.Context, fn func(ctx context.Context, tx store.Tx) error) error {
	return fn(ctx, struct{}{})
}

func TestOrderSyncServerCreateOrderMapsProtoRequest(t *testing.T) {
	orderStore := &fakeOrderStore{}
	orderItemStore := &fakeOrderItemStore{}
	orderSvc := ordersvc.NewSvc(util.NewCtxFunc(time.Second), orderStore, orderItemStore)
	server := NewOrderSyncServer(orderSvc, &fakeDb{})

	resp, err := server.CreateOrder(context.Background(), &orderbotv1pb.CreateOrderRequest{
		OrderId:     "order-1",
		BotId:       "bot-1",
		CartId:      "cart-1",
		SessionId:   "session-1",
		TotalScaled: 4200,
		Items: []*orderbotv1pb.OrderItem{
			{
				Id:               "item-1",
				MenuItemId:       "menu-item-1",
				Name:             "Latte",
				Quantity:         2,
				UnitPriceScaled:  2100,
				TotalPriceScaled: 4200,
			},
		},
	})
	if err != nil {
		t.Fatalf("CreateOrder() error = %v", err)
	}
	if !resp.Ok {
		t.Fatalf("CreateOrder() response ok = false")
	}
	if orderStore.inserted.ID != "order-1" || orderStore.inserted.CartID != "cart-1" {
		t.Fatalf("unexpected inserted order: %+v", orderStore.inserted)
	}
	if len(orderItemStore.inserted) != 1 {
		t.Fatalf("expected 1 inserted item, got %d", len(orderItemStore.inserted))
	}
	if got := orderItemStore.inserted[0]; got.OrderID != "order-1" || got.MenuItemID != "menu-item-1" || got.Quantity != 2 {
		t.Fatalf("unexpected inserted item: %+v", got)
	}
}
