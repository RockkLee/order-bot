PROTO_DIR := ../proto
PROTO_V1_DIR := $(PROTO_DIR)/orderbot/v1
PY_PROTO_OUT := src/infra/grpcpb
PROTO_FILES := \
	$(PROTO_V1_DIR)/common.proto \
	$(PROTO_V1_DIR)/order_sync_service.proto \
	$(PROTO_V1_DIR)/order_status_service.proto

gen-proto-py:
	@echo "Generating Python protobuf/grpc stubs..."
	@mkdir -p $(PY_PROTO_OUT)
	@python -m grpc_tools.protoc -I $(PROTO_DIR) \
		--python_out=$(PY_PROTO_OUT) \
		--grpc_python_out=$(PY_PROTO_OUT) \
		$(PROTO_FILES)
	@touch $(PY_PROTO_OUT)/__init__.py

.PHONY: gen-proto-py
