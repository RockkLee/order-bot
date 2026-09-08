package grpcclient

import (
	"context"
	"fmt"

	orderbotv1pb "github.com/RockkLee/order-bot/goproto/orderbot/v1"
	orderbotsvcpb "github.com/RockkLee/order-bot/goproto/orderbot/v1/order_bot_svc"
	"google.golang.org/grpc"
)

func MarkOrderCompleted(ctx context.Context, conn grpc.ClientConnInterface, orderID string) error {
	client := orderbotsvcpb.NewOrderStatusServiceClient(conn)
	resp, err := client.MarkOrderCompleted(ctx, &orderbotv1pb.MarkOrderCompletedRequest{OrderId: orderID})
	if err != nil {
		return fmt.Errorf("grpcclient.MarkOrderCompleted rpc: %w", err)
	}
	if !resp.Ok {
		return fmt.Errorf("grpcclient.MarkOrderCompleted: response ok=false")
	}
	return nil
}
