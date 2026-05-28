# Proto ownership for `order-bot-mgmt-svc`

This project **consumes** shared gRPC contracts from the repository root:

- `../../proto/order_sync.proto`

We intentionally keep one canonical proto source at repo root to avoid drift between
services. Generated Go stubs for this service should be output into this project
(e.g. under `internal/infra/grpcpb/`) while proto definitions remain shared.
