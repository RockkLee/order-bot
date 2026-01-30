-- DDL for order-bot-mgmt-svc entities.

CREATE TABLE IF NOT EXISTS bot (
    id TEXT PRIMARY KEY,
    bot_name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS menu (
    id TEXT PRIMARY KEY,
    bot_id TEXT NOT NULL REFERENCES bot(id)
);

CREATE TABLE IF NOT EXISTS menu_item (
    id TEXT PRIMARY KEY,
    menu_id TEXT NOT NULL REFERENCES menu(id),
    menu_item_name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    access_token TEXT NOT NULL,
    refresh_token TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS user_bot (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id),
    bot_id TEXT NOT NULL REFERENCES bot(id)
);

CREATE INDEX IF NOT EXISTS idx_menu_bot_id ON menu (bot_id);
CREATE INDEX IF NOT EXISTS idx_menu_item_menu_id ON menu_item (menu_id);
CREATE INDEX IF NOT EXISTS idx_user_bot_user_id ON user_bot (user_id);
CREATE INDEX IF NOT EXISTS idx_user_bot_bot_id ON user_bot (bot_id);
