package grpcserver

import (
	"context"
	"net"
	"testing"
	"time"

	orderbotv1pb "github.com/RockkLee/order-bot/goproto/orderbot/v1"
	orderbotmgmtsvcpb "github.com/RockkLee/order-bot/goproto/orderbot/v1/order_bot_mgmt_svc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	"order-bot-mgmt-svc/internal/services/ordersvc"
	"order-bot-mgmt-svc/internal/util"
)

func TestNewServerServesOrderSyncService(t *testing.T) {
	orderStore := &fakeOrderStore{}
	orderItemStore := &fakeOrderItemStore{}
	orderSvc := ordersvc.NewSvc(util.NewCtxFunc(time.Second), orderStore, orderItemStore)

	lis := bufconn.Listen(1024 * 1024)
	defer func() { _ = lis.Close() }()

	grpcSrv := NewServer(orderSvc, &fakeDb{})
	defer grpcSrv.Stop()

	errCh := make(chan error, 1)
	go func() {
		errCh <- grpcSrv.Serve(lis)
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.NewClient(
		"passthrough:///bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("grpc.NewClient() error = %v", err)
	}
	defer func() { _ = conn.Close() }()

	client := orderbotmgmtsvcpb.NewOrderSyncServiceClient(conn)
	resp, err := client.CreateOrder(ctx, &orderbotv1pb.CreateOrderRequest{
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

	grpcSrv.Stop()
	select {
	case serveErr := <-errCh:
		if serveErr != nil {
			t.Fatalf("Serve() error = %v", serveErr)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for Serve() to stop")
	}
}
