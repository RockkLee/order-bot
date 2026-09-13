package config

import "testing"

func TestLoadIncludesGrpcConfig(t *testing.T) {
	t.Setenv("ADDRESS", "127.0.0.1")
	t.Setenv("PORT", "8080")
	t.Setenv("GIN_MODE", "debug")
	t.Setenv("GRPC_ADDRESS", "127.0.0.1")
	t.Setenv("GRPC_PORT", "9090")
	t.Setenv("ORDER_BOT_GRPC_CLIENT_ADDRESS", "127.0.0.2")
	t.Setenv("ORDER_BOT_GRPC_CLIENT_PORT", "9091")
	t.Setenv("BLUEPRINT_DB_DATABASE", "blueprint")
	t.Setenv("BLUEPRINT_DB_PASSWORD", "password1234")
	t.Setenv("BLUEPRINT_DB_USERNAME", "melkey")
	t.Setenv("BLUEPRINT_DB_PORT", "5432")
	t.Setenv("BLUEPRINT_DB_HOST", "localhost")
	t.Setenv("BLUEPRINT_DB_SCHEMA", "order_bot_mgmt")
	t.Setenv("BLUEPRINT_DB_ORDER_BOT_SCHEMA", "order_bot")

	cfg := Load()

	if cfg.GrpcServer.Address != "127.0.0.1" {
		t.Fatalf("cfg.GrpcServer.Address = %q, want %q", cfg.GrpcServer.Address, "127.0.0.1")
	}
	if cfg.GrpcServer.Port != 9090 {
		t.Fatalf("cfg.GrpcServer.Port = %d, want %d", cfg.GrpcServer.Port, 9090)
	}
	if cfg.OrderBotGrpcClient.Address != "127.0.0.2" {
		t.Fatalf("cfg.OrderBotGrpcClient.Address = %q, want %q", cfg.OrderBotGrpcClient.Address, "127.0.0.2")
	}
	if cfg.OrderBotGrpcClient.Port != 9091 {
		t.Fatalf("cfg.OrderBotGrpcClient.Port = %d, want %d", cfg.OrderBotGrpcClient.Port, 9091)
	}
}
