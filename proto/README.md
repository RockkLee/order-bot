# Enterprise proto layout

This repository uses a shared enterprise-style protobuf layout under `proto/orderbot/v1`.

## Structure

- `proto/orderbot/v1/common.proto`
  - Shared request/response and DTO message contracts.
- `proto/orderbot/v1/order_sync_service.proto`
  - `OrderSyncService` contract.
- `proto/orderbot/v1/order_status_service.proto`
  - `OrderStatusService` contract.
- `proto/order_sync.proto`
  - Compatibility umbrella entrypoint importing all contracts for existing tooling.

## Rationale

- Single source of truth for cross-service API contracts.
- Versioned namespace (`orderbot.v1`) for backward-compatible evolution.
- Service definitions split from message definitions to keep contracts modular and reviewable.
