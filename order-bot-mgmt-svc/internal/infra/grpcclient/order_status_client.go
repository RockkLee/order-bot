package grpcclient

import "context"

// NOTE: replace with generated OrderStatusServiceClient from proto/order_sync.proto.
type OrderStatusClient struct {
	target string
}

func NewOrderStatusClient(target string) *OrderStatusClient {
	return &OrderStatusClient{target: target}
}

func (c *OrderStatusClient) MarkOrderCompleted(ctx context.Context, orderID string) error {
	// TODO: use grpc.Dial + generated pb client and call MarkOrderCompleted.
	_, _, _ = ctx, orderID, c.target
	return nil
}
