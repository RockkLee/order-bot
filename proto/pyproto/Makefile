PROTO_DIR := ../proto
COMMON_PROTO := $(PROTO_DIR)/orderbot/v1/common.proto
MGMT_PROTO_DIR := $(PROTO_DIR)/orderbot/v1/order_bot_mgmt_svc
SVC_PROTO_DIR := $(PROTO_DIR)/orderbot/v1/order_bot_svc
PYTHON ?= python3
PY_PROTO_OUT := .
MGMT_PROTO_FILES := \
	$(COMMON_PROTO) \
	$(MGMT_PROTO_DIR)/order_sync_service.proto
SVC_PROTO_FILES := \
	$(COMMON_PROTO) \
	$(SVC_PROTO_DIR)/order_status_service.proto
PY_INIT_FILES := \
	orderbot/__init__.py \
	orderbot/v1/__init__.py \
	orderbot/v1/order_bot_mgmt_svc/__init__.py \
	orderbot/v1/order_bot_svc/__init__.py

gen-proto-py-mgmt:
	@echo "Generating Python protobuf/grpc stubs for order-bot-mgmt-svc..."
	@$(PYTHON) -m grpc_tools.protoc -I $(PROTO_DIR) \
		--python_out=$(PY_PROTO_OUT) \
		--grpc_python_out=$(PY_PROTO_OUT) \
		$(MGMT_PROTO_FILES)
	@for file in $(PY_INIT_FILES); do \
		mkdir -p "$$(dirname "$$file")"; \
		touch "$$file"; \
	done

gen-proto-py-svc:
	@echo "Generating Python protobuf/grpc stubs for order-bot-svc..."
	@$(PYTHON) -m grpc_tools.protoc -I $(PROTO_DIR) \
		--python_out=$(PY_PROTO_OUT) \
		--grpc_python_out=$(PY_PROTO_OUT) \
		$(SVC_PROTO_FILES)
	@for file in $(PY_INIT_FILES); do \
		mkdir -p "$$(dirname "$$file")"; \
		touch "$$file"; \
	done

gen-proto-py: gen-proto-py-mgmt gen-proto-py-svc

.PHONY: gen-proto-py gen-proto-py-mgmt gen-proto-py-svc
