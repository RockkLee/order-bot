-- DDL for order-bot-svc entities.

CREATE TABLE IF NOT EXISTS published_menu (
    id TEXT PRIMARY KEY,
    bot_id TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS published_menu_item (
    id TEXT PRIMARY KEY,
    menu_id TEXT NOT NULL REFERENCES published_menu(id),
    menu_item_name TEXT NOT NULL,
    price DOUBLE PRECISION
);
CREATE INDEX IF NOT EXISTS idx_menu_item_menu_id ON published_menu_item (menu_id);

CREATE TABLE IF NOT EXISTS orders (
    id TEXT PRIMARY KEY,
    cart_id TEXT NOT NULL,
    session_id TEXT NOT NULL,
    total_scaled INTEGER NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_orders_session_id ON orders (session_id);

CREATE TABLE IF NOT EXISTS order_item (
    id TEXT PRIMARY KEY,
    order_id TEXT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    menu_item_id TEXT NOT NULL,
    name TEXT NOT NULL,
    quantity INTEGER NOT NULL,
    unit_price_scaled INTEGER NOT NULL,
    total_price_scaled INTEGER NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_order_item_order_id ON order_item (order_id);
