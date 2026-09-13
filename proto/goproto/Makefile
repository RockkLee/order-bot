PROTO_DIR := ../proto
COMMON_PROTO := $(PROTO_DIR)/orderbot/v1/common.proto
MGMT_PROTO_DIR := $(PROTO_DIR)/orderbot/v1/order_bot_mgmt_svc
SVC_PROTO_DIR := $(PROTO_DIR)/orderbot/v1/order_bot_svc
GO_PROTO_OUT := .
GO_PROTO_MODULE := github.com/RockkLee/order-bot/goproto
MGMT_PROTO_FILES := \
	$(COMMON_PROTO) \
	$(MGMT_PROTO_DIR)/order_sync_service.proto
SVC_PROTO_FILES := \
	$(COMMON_PROTO) \
	$(SVC_PROTO_DIR)/order_status_service.proto

gen-datamodels-from-protobuf-mgmt:
	@echo "Run gen-datamodels-from-protobuf-mgmt..."
	@protoc -I $(PROTO_DIR) \
		--go_out=$(GO_PROTO_OUT) \
		--go_opt=module=$(GO_PROTO_MODULE) \
		$(MGMT_PROTO_FILES)

gen-datamodels-from-protobuf-svc:
	@echo "Run gen-datamodels-from-protobuf-svc..."
	@protoc -I $(PROTO_DIR) \
		--go_out=$(GO_PROTO_OUT) \
		--go_opt=module=$(GO_PROTO_MODULE) \
		$(SVC_PROTO_FILES)

gen-rpccode-from-protobuf-mgmt:
	@echo "Run gen-rpccode-from-protobuf-mgmt..."
	@protoc -I $(PROTO_DIR) \
		--go-grpc_out=$(GO_PROTO_OUT) \
		--go-grpc_opt=module=$(GO_PROTO_MODULE) \
		$(MGMT_PROTO_FILES)

gen-rpccode-from-protobuf-svc:
	@echo "Run gen-rpccode-from-protobuf-svc..."
	@protoc -I $(PROTO_DIR) \
		--go-grpc_out=$(GO_PROTO_OUT) \
		--go-grpc_opt=module=$(GO_PROTO_MODULE) \
		$(SVC_PROTO_FILES)

gen-datamodels-from-protobuf: gen-datamodels-from-protobuf-mgmt gen-datamodels-from-protobuf-svc

gen-rpccode-from-protobuf: gen-rpccode-from-protobuf-mgmt gen-rpccode-from-protobuf-svc

gen-proto-go: gen-datamodels-from-protobuf gen-rpccode-from-protobuf

.PHONY: \
	gen-proto-go \
	gen-datamodels-from-protobuf \
	gen-rpccode-from-protobuf \
	gen-datamodels-from-protobuf-mgmt \
	gen-datamodels-from-protobuf-svc \
	gen-rpccode-from-protobuf-mgmt \
	gen-rpccode-from-protobuf-svc
