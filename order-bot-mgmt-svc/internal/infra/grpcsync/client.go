package grpcsync

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type CallbackClient struct {
	addr string
}

func NewCallbackClient(addr string) *CallbackClient { return &CallbackClient{addr: addr} }

func (c *CallbackClient) UpdateOrderStatus(ctx context.Context, req UpdateOrderStatusRequest) error {
	codec := JSONCodec{}
	conn, err := grpc.DialContext(ctx, c.addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultCallOptions(grpc.ForceCodec(codec)))
	if err != nil {
		return fmt.Errorf("grpcsync.CallbackClient.UpdateOrderStatus: %w", err)
	}
	defer conn.Close()

	var resp UpdateOrderStatusResponse
	if err := conn.Invoke(ctx, "/ordersync.OrderCallbackService/UpdateOrderStatus", &req, &resp); err != nil {
		return fmt.Errorf("grpcsync.CallbackClient.UpdateOrderStatus: %w", err)
	}
	return nil
}
