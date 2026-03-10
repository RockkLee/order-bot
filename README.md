# order-bot

> A full-stack text-to-order chatbot demo app that allows: <br>
> - B-side users to manage the order bot's menu and review orders from the admin side. <br>
> - C-side users to place orders by texting with the order bot

<br>

## Features (Use Case Diagram)
```mermaid
flowchart LR

  Admin["Customer Admin (B-side)"]
  Customer["Customer (C-side)"]

  %% direction LR

  subgraph CAS["Customer Admin Service (B-side)"]
    direction TB
    UC_Auth(["Authentication"])
    UC_Signup(["Sign up"])
    UC_Login(["Log in"])
    UC_Logout(["Log out"])
    UC_ManageBot(["Manage Order Bot"])
    UC_ImportMenu(["Manage a menu per bot"])
    UC_RecordEvents(["Record events"])
  end

  subgraph OBS["Order Bot Service (C-side)"]
    direction TB
    UC_OrderItems(["Order items"])
    UC_BotValidation(["Check if the botID & menuID are correct"])
    UC_TextToItems(["Text to order items"])
    UC_DisplayItems(["Display order items"])
    UC_SendEvent(["Send an event to Admin Service"])
  end

  Admin --> UC_Auth
  UC_Auth -.-> UC_Signup
  UC_Auth -.-> UC_Login
  UC_Auth -.-> UC_Logout
  Admin --> UC_ManageBot
  UC_ManageBot -.-> UC_ImportMenu
  UC_ManageBot -.-> UC_RecordEvents

  Customer --> UC_OrderItems

  UC_OrderItems -.-> UC_BotValidation
  UC_OrderItems -.-> UC_TextToItems
  UC_OrderItems -.-> UC_DisplayItems
  UC_OrderItems -.-> UC_SendEvent

  UC_SendEvent -.-> |TODO: <br>For now, read events directly from the DB| UC_RecordEvents
```

<br>

## Modules
### `order-bot-mgmt-svc`
* Go-based management service for B-side (admin) operations.
* Stack:
  * net/http
  * Gin
  * GORM

### `order-bot-svc`
* Python-based C-side order bot service.
* Stack:
  * FastAPI
  * SQLAlchemy
  * LangChain
  * MCP

### `order-bot-frontend`
* Vue 3 + Vite frontend for user interaction.

### `infra-terraform`
* Terraform-based infrastructure definitions to deploy the app in the AWS Cloud.
* Beside the EventBridge schedules, Lambda functions, and VPC-related settings, all the other 
  settings are defined in this module.
* Infrastructure Architecture: 
  * [architecture.md](./infra-terraform/docs/architecture.md)
  * [start_flow.md](./infra-terraform/docs/start_flow.md)

## How to run the app in local environment
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

<br>


## ER-Diagram
* [order-bot-mgmt-svc-erd](./doc/order-bot-mgmt-svc-erd.md)
* [order-bot-svc-erd](./doc/order-bot-svc-erd.md)

<br>

## TODO (Relatively critical features that need to be added)

* `order-bot-mgmt-svc`
  * The service does create access tokens and refresh tokens, but there is no 
  token refresh mechanism for now.
  * The implementation of Dependency Isolation is a bit messy. The related structure still needs 
    to be refactored.
* `order-bot-svc`
  * More MCP tools need to be added to allow the bot to provide more types of response for a better user experience
  * The output model to retrieve the response from the MCP tool still need to be refactored
* `order-bot-frontend`
  * The UI layout is off on mobile browsers
  * The C-side URL can't be copied after clicking the URL generation button on a mobile browser
  * The order data on the dialog page is not rendered yet. It's shown in plain JSON format.
* `infra-terraform`
  * There is still an issue with the cache of AWS CloudFront.
  * Clients may get stale cache files after re-uploading files to linked S3 buckets
