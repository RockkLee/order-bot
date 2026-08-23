# ER Diagram: order-bot-mgmt-svc

```mermaid
erDiagram
    direction LR
    users {
        string id PK
        string email
        string password_hash
        string access_token
        string refresh_token
        datetime created_at
        datetime updated_at
    }

    bot {
        string id PK
        string bot_name
        datetime created_at
        datetime updated_at
    }

    user_bot {
        string id PK
        string user_id
        string bot_id
        datetime created_at
        datetime updated_at
    }

    menu {
        string id PK
        string bot_id
        datetime created_at
        datetime updated_at
    }

    menu_item {
        string id PK
        string menu_id
        string menu_item_name
        float price
        datetime created_at
        datetime updated_at
    }

    recipe {
        string menu_item_id PK, FK
        string sku PK, FK
        int required_quantity
    }

    inventory {
        string sku PK
        string name
        string base_unit
    }

    store_inventory {
        string store_id PK, FK
        string sku PK, FK
        int on_hand_quantity
        int reserved_quantity
        datetime updated_at
    }

    store {
        string id PK
        string name
        string location
    }

    orders {
        string id PK
        string bot_id
        string cart_id
        string session_id
        int total_scaled
        datetime created_at
        datetime updated_at
    }

    order_item {
        string id PK
        string order_id
        string menu_item_id
        string name
        int quantity
        int unit_price_scaled
        int total_price_scaled
        datetime created_at
        datetime updated_at
    }

    users ||--o{ user_bot : ""
    bot ||--o{ user_bot : ""
    bot ||--o{ menu : ""
    menu ||--o{ menu_item : ""
    menu_item ||--o{ recipe : ""
    inventory ||--o{ recipe : ""
    inventory ||--o{ store_inventory : ""
    store ||--o{ store_inventory : ""
    bot ||--o{ orders : ""
    orders ||--o{ order_item : ""
    menu_item ||--o{ order_item : ""
```
