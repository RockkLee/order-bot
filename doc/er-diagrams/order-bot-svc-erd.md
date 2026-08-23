# ER-Diagram: order-bot-svc

```mermaid
erDiagram
    published_menu {
        string id PK
        string bot_id "indexed"
        datetime created_at
        datetime updated_at
    }

    published_menu_item {
        string id PK
        string menu_id UK "indexed"
        string menu_item_name
        float price
        datetime created_at
        datetime updated_at
    }

    cart {
        string id PK
        string session_id UK "indexed"
        string status "cart_status; default OPEN"
        int total_scaled "default 0"
        datetime closed_at "nullable"
        datetime created_at
        datetime updated_at
    }

    cart_item {
        string id PK
        string cart_id FK
        string menu_item_id "unique with cart_id"
        string name
        int quantity
        int unit_price_scaled
        int total_price_scaled
        datetime created_at
        datetime updated_at
    }

    orders {
        string id PK
        string cart_id FK
        string bot_id "indexed"
        string session_id "indexed"
        int total_scaled
        datetime created_at
        datetime updated_at
    }

    order_item {
        string id PK
        string order_id FK
        string menu_item_id
        string name
        int quantity
        int unit_price_scaled
        int total_price_scaled
        datetime created_at
        datetime updated_at
    }

    cart ||--o{ cart_item : contains
    cart ||--o{ orders : creates
    orders ||--o{ order_item : contains

```

`created_at` and `updated_at` are inherited from `BaseModel`. The diagram includes only relationships enforced by `ForeignKey` declarations; menu-related identifier columns are not foreign keys in the current model.
