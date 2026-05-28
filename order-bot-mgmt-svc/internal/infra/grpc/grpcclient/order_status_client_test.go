package grpcclient

import (
	"context"
	"net"
	"testing"
	"time"

	orderbotv1pb "github.com/RockkLee/order-bot/goproto/orderbot/v1"
	orderbotsvcpb "github.com/RockkLee/order-bot/goproto/orderbot/v1/order_bot_svc"
	"google.golang.org/grpc"
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
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer lis.Close()

	grpcServer := grpc.NewServer()
	svc := &testOrderStatusServer{}
	orderbotsvcpb.RegisterOrderStatusServiceServer(grpcServer, svc)
	defer grpcServer.Stop()

	go func() {
		_ = grpcServer.Serve(lis)
	}()

	client := NewOrderStatusClient(lis.Addr().String())
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := client.MarkOrderCompleted(ctx, "order-123"); err != nil {
		t.Fatalf("MarkOrderCompleted() error = %v", err)
	}
	if svc.lastOrderID != "order-123" {
		t.Fatalf("expected order ID order-123, got %q", svc.lastOrderID)
	}
}
