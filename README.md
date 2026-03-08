# order-bot

> A full-stack demo platform for building and operating conversational order bots: manage bots and menus on the admin side, then place and track customer orders through an API-first order flow.

## Features

### `order-bot-mgmt-svc`
* Go-based management service for B-side (admin) operations.
* Handles core management entities such as users, bots, menus, and menu items.
* Supports bot ownership/sharing relationships through `user_bot` mappings.
* Designed to power admin workflows like authentication, bot configuration, and menu management.

### `order-bot-svc`
* FastAPI-based C-side order bot service.
* Session-aware cart lifecycle (`create/update/remove/show cart`, `checkout`).
* Persists carts, orders, and order items in PostgreSQL.
* Includes test coverage under `tests/` and supports containerized startup.

### `order-bot-frontend`
* Vue 3 + Vite frontend for user interaction.
* Provides a client application layer for order bot flows.
* Includes unit test and E2E test setup (Vitest + Playwright).

### `ddl`
* Contains SQL DDL scripts for both schemas:
  * `order_bot` (runtime ordering domain)
  * `order_bot_mgmt` (management/admin domain)
* Used to initialize database tables after schema creation.

### `doc`
* Architecture and modeling documentation:
  * Use case diagram
  * ER diagrams for both services
* Useful for understanding service boundaries and data relationships.

### `infra-terraform`
* Terraform-based infrastructure definitions for deployment environments.
* Module-oriented structure for provisioning cloud resources.

## How to install

1. **Build service images from each service Dockerfile** (or let Compose build automatically):
   * `order-bot-mgmt-svc/Dockerfile`
   * `order-bot-svc/Dockerfile`
2. **Run Docker Compose from the repository root**:
   ```bash
   docker compose up --build -d
   ```
3. **Create two PostgreSQL schemas**:
   * `order_bot`
   * `order_bot_mgmt`

   Example:
   ```sql
   CREATE SCHEMA IF NOT EXISTS order_bot;
   CREATE SCHEMA IF NOT EXISTS order_bot_mgmt;
   ```
4. **Run DDL files to create tables**:
   * `ddl/order_bot_ddl.sql`
   * `ddl/order_bot_mgmt_ddl.sql`

   Example:
   ```bash
   psql -h <host> -U <user> -d <database> -f ddl/order_bot_ddl.sql
   psql -h <host> -U <user> -d <database> -f ddl/order_bot_mgmt_ddl.sql
   ```

## TODO

* Add a complete local `.env.example` with all service/database variables documented.
* Replace placeholder/demo auth handling with production-ready token lifecycle and secret management.
* Integrate a real LLM provider flow for intent parsing and ordering orchestration.
* Implement event delivery from `order-bot-svc` to `order-bot-mgmt-svc` (currently noted as TODO in docs).
* Expand API and integration test coverage across services.
* Add CI pipeline gates for lint/test/build across frontend, Python, and Go services.
