package config

import "testing"

func TestLoadIncludesGrpcConfig(t *testing.T) {
	t.Setenv("ADDRESS", "127.0.0.1")
	t.Setenv("PORT", "8080")
	t.Setenv("GIN_MODE", "debug")
	t.Setenv("GRPC_ADDRESS", "127.0.0.1")
	t.Setenv("GRPC_PORT", "9090")
	t.Setenv("ORDER_BOT_GRPC_ADDRESS", "127.0.0.2")
	t.Setenv("ORDER_BOT_GRPC_PORT", "9091")
	t.Setenv("BLUEPRINT_DB_DATABASE", "blueprint")
	t.Setenv("BLUEPRINT_DB_PASSWORD", "password1234")
	t.Setenv("BLUEPRINT_DB_USERNAME", "melkey")
	t.Setenv("BLUEPRINT_DB_PORT", "5432")
	t.Setenv("BLUEPRINT_DB_HOST", "localhost")
	t.Setenv("BLUEPRINT_DB_SCHEMA", "order_bot_mgmt")
	t.Setenv("BLUEPRINT_DB_ORDER_BOT_SCHEMA", "order_bot")

	cfg := Load()

	if cfg.Grpc.Address != "127.0.0.1" {
		t.Fatalf("cfg.Grpc.Address = %q, want %q", cfg.Grpc.Address, "127.0.0.1")
	}
	if cfg.Grpc.Port != 9090 {
		t.Fatalf("cfg.Grpc.Port = %d, want %d", cfg.Grpc.Port, 9090)
	}
	if cfg.OrderBotGrpc.Address != "127.0.0.2" {
		t.Fatalf("cfg.OrderBotGrpc.Address = %q, want %q", cfg.OrderBotGrpc.Address, "127.0.0.2")
	}
	if cfg.OrderBotGrpc.Port != 9091 {
		t.Fatalf("cfg.OrderBotGrpc.Port = %d, want %d", cfg.OrderBotGrpc.Port, 9091)
	}
}
