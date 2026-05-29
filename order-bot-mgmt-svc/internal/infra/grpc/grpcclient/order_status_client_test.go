package grpcclient

import (
	"context"
	"net"
	"testing"
	"time"

	orderbotv1pb "github.com/RockkLee/order-bot/goproto/orderbot/v1"
	orderbotsvcpb "github.com/RockkLee/order-bot/goproto/orderbot/v1/order_bot_svc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

type testOrderStatusServer struct {
	orderbotsvcpb.UnimplementedOrderStatusServiceServer
	lastOrderID string
}

func (s *testOrderStatusServer) MarkOrderCompleted(ctx context.Context, req *orderbotv1pb.MarkOrderCompletedRequest) (*orderbotv1pb.MarkOrderCompletedResponse, error) {
	s.lastOrderID = req.OrderId
	return &orderbotv1pb.MarkOrderCompletedResponse{Ok: true}, nil
}

func TestOrderStatusClientMarkOrderCompleted(t *testing.T) {
	lis := bufconn.Listen(1024 * 1024)
	defer func() { _ = lis.Close() }()

	grpcServer := grpc.NewServer()
	svc := &testOrderStatusServer{}
	orderbotsvcpb.RegisterOrderStatusServiceServer(grpcServer, svc)
	defer grpcServer.Stop()

	go func() {
		_ = grpcServer.Serve(lis)
	}()

	client := newOrderStatusClient(
		"passthrough:///bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := client.MarkOrderCompleted(ctx, "order-123"); err != nil {
		t.Fatalf("MarkOrderCompleted() error = %v", err)
	}
	if svc.lastOrderID != "order-123" {
		t.Fatalf("expected order ID order-123, got %q", svc.lastOrderID)
	}
}
