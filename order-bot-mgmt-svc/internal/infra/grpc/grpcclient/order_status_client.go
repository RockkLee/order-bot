package grpcclient

import (
	"context"
	"fmt"

	orderbotv1pb "github.com/RockkLee/order-bot/goproto/orderbot/v1"
	orderbotsvcpb "github.com/RockkLee/order-bot/goproto/orderbot/v1/order_bot_svc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type OrderStatusClient struct {
	target string
}

func NewOrderStatusClient(target string) *OrderStatusClient {
	return &OrderStatusClient{target: target}
}

func (c *OrderStatusClient) MarkOrderCompleted(ctx context.Context, orderID string) error {
	conn, err := grpc.DialContext(ctx, c.target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("grpcclient.OrderStatusClient.MarkOrderCompleted dial: %w", err)
	}
	defer func() { _ = conn.Close() }()

	client := orderbotsvcpb.NewOrderStatusServiceClient(conn)
	resp, err := client.MarkOrderCompleted(ctx, &orderbotv1pb.MarkOrderCompletedRequest{OrderId: orderID})
	if err != nil {
		return fmt.Errorf("grpcclient.OrderStatusClient.MarkOrderCompleted rpc: %w", err)
	}
	if !resp.Ok {
		return fmt.Errorf("grpcclient.OrderStatusClient.MarkOrderCompleted: response ok=false")
	}
	return nil
}
