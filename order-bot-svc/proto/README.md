# Proto ownership for `order-bot-svc`

This project **consumes** shared gRPC contracts from the repository root:

- `../../proto/order_sync.proto`

We intentionally keep one canonical proto source at repo root to avoid drift between
services. Generated Python stubs for this service should be output into this project
(e.g. under `src/infra/grpcpb/`) but the proto definition itself stays shared.
