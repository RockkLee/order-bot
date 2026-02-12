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
