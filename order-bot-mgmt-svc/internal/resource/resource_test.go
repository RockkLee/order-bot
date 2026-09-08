package resource

import (
	"testing"

	"order-bot-mgmt-svc/internal/config"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestNewStoresInfrastructureHandles(t *testing.T) {
	conn, err := grpc.NewClient("passthrough:///unused", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("grpc.NewClient() error = %v", err)
	}

	res := New(nil, nil, GRPCConn{OrderBot: conn})
	if res.DB != nil {
		t.Fatalf("expected DB to be nil")
	}
	if res.OrderBotDB != nil {
		t.Fatalf("expected OrderBotDB to be nil")
	}
	if res.GRPCConn.OrderBot != conn {
		t.Fatalf("expected GRPCConn.OrderBot to match input conn")
	}
	if err := res.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
}

func TestNewOrderBotGRPCConn(t *testing.T) {
	conn, err := NewOrderBotGRPCConn(config.Grpc{
		Address: "127.0.0.1",
		Port:    9090,
	})
	if err != nil {
		t.Fatalf("NewOrderBotGRPCConn() error = %v", err)
	}
	if conn == nil {
		t.Fatalf("expected conn to be non-nil")
	}
	if err := conn.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
}
