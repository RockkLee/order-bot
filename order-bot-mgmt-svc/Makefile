# Simple Makefile for a Go project

PROTO_DIR := ../proto
PROTO_V1_DIR := $(PROTO_DIR)/orderbot/v1
GO_PROTO_OUT := .
PROTO_FILES := \
	$(PROTO_V1_DIR)/common.proto \
	$(PROTO_V1_DIR)/order_sync_service.proto \
	$(PROTO_V1_DIR)/order_status_service.proto

# Build the application
all: build test

build:
	@echo "Building..."
	@go build -o main cmd/api/main.go

# Generate Go gRPC stubs from shared proto contracts
gen-proto-go:
	@echo "Generating Go protobuf/grpc stubs..."
	@protoc -I $(PROTO_DIR) \
		--go_out=$(GO_PROTO_OUT) \
		--go-grpc_out=$(GO_PROTO_OUT) \
		$(PROTO_FILES)

# Run the application
run:
	@go run cmd/api/main.go

# Create DB container
docker-run:
	@if docker compose up --build 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose up --build; \
	fi

# Shutdown DB container
docker-down:
	@if docker compose down 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose down; \
	fi

# Test the application
test:
	@echo "Testing..."
	@go test ./... -v

# Integrations Tests for the application
itest:
	@echo "Running integration tests..."
	@go test ./internal/database -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

# Live Reload
watch:
	@if command -v air > /dev/null; then \
		air; \
		echo "Watching...";\
	else \
		read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
		if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
			go install github.com/air-verse/air@latest; \
			air; \
			echo "Watching...";\
		else \
			echo "You chose not to install air. Exiting..."; \
			exit 1; \
		fi; \
	fi

.PHONY: all build run test clean watch docker-run docker-down itest gen-proto-go
