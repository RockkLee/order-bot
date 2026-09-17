# Enterprise proto layout

This repository uses a shared enterprise-style protobuf layout under `proto/orderbot/v1`.

The protobuf module root is `proto/`. Imports such as `orderbot/v1/order_bot_mgmt_svc/common.proto`
are resolved relative to that directory. Tooling should compile with `-I proto`
or use [buf.yaml](/home/hiddenlotus/Desktop/Code/order-bot/proto/buf.yaml).

Generated code destinations are configured by the consumer's generator command,
not by an entrypoint wrapper proto. For Go, each generated
`.proto` still needs its own `go_package` because that option is file-local.
This repository uses shared language-specific generated libraries while keeping
the source `.proto` files grouped by owning app.
With plain `protoc`, import-only wrapper protos do not generate code for their
imports, so consumers still need to list the concrete service proto plus any
shared dependency protos they want emitted.

## Structure

- `proto/orderbot/v1/order_bot_mgmt_svc/common.proto`
  - Shared request/response and DTO message contracts owned by `order-bot-mgmt-svc`.
- `proto/orderbot/v1/order_bot_mgmt_svc/order_sync_service.proto`
  - `OrderSyncService` contract.
- `proto/orderbot/v1/order_bot_svc/common.proto`
  - Shared request/response and DTO message contracts owned by `order-bot-svc`.
- `proto/orderbot/v1/order_bot_svc/order_status_service.proto`
  - `OrderStatusService` contract.
- `proto/orderbot/v1/order_bot_mgmt_svc.proto`
  - `order-bot-mgmt-svc` ownership wrapper importing only `order_sync_service.proto`.
- `proto/orderbot/v1/order_bot_svc.proto`
  - `order-bot-svc` ownership wrapper importing only `order_status_service.proto`.

## Generated Libraries

- Shared Go protobuf/gRPC code is generated into `goproto`.
- The canonical Go import paths are rooted at `github.com/RockkLee/order-bot/goproto/...`.
- Shared Python protobuf/gRPC code is generated into `pyproto`.
- App services should consume `goproto` or `pyproto` rather than owning shared codegen themselves.

## Rationale

- Single source of truth for cross-service API contracts.
- Versioned namespace (`orderbot.v1`) for backward-compatible evolution.
- Service definitions split from message definitions to keep contracts modular and reviewable.
- Per-project wrapper protos document ownership, while each consumer `Makefile`
  explicitly generates the shared message proto plus the one service proto it needs.
- App ownership in the proto path stays separate from the canonical generated Go import path.
