-- order_bot_schema.sql
-- Generated from the provided SQLAlchemy models (PostgreSQL)

-- Optional: keep everything in a dedicated schema
-- CREATE SCHEMA IF NOT EXISTS order_bot;
-- SET search_path TO order_bot;

-- SQLAlchemy uses "Enum(..., native_enum=False)" so we model it as a TEXT + CHECK constraint
-- (instead of a PostgreSQL ENUM type).
-- Adjust allowed values to match your CartStatus enum definition.
-- Common guess: OPEN / CLOSED
CREATE TABLE IF NOT EXISTS cart (
  id            VARCHAR(36) PRIMARY KEY,
  session_id    VARCHAR(36) NOT NULL UNIQUE,
  status        TEXT NOT NULL DEFAULT 'OPEN',
  total_scaled  INTEGER NOT NULL DEFAULT 0,
  closed_at     TIMESTAMP NULL,

  CONSTRAINT cart_status_chk
    CHECK (status IN ('OPEN', 'CLOSED'))
);

-- Index because SQLAlchemy used index=True on session_id
CREATE INDEX IF NOT EXISTS ix_cart_session_id ON cart (session_id);

CREATE TABLE IF NOT EXISTS cart_item (
  id                 VARCHAR(36) PRIMARY KEY,
  cart_id            VARCHAR(36) NOT NULL,
  menu_item_id       VARCHAR(64) NOT NULL,
  name               VARCHAR(200) NOT NULL,
  quantity           INTEGER NOT NULL,
  unit_price_scaled  INTEGER NOT NULL,
  total_price_scaled INTEGER NOT NULL,

  CONSTRAINT fk_cart_item_cart
    FOREIGN KEY (cart_id) REFERENCES cart (id)
    ON DELETE CASCADE,

  -- Your SQLAlchemy constraint name is "menu_item_id", which is confusing because it's the
  -- *constraint name*, not the column. Kept as-is to match.
  CONSTRAINT menu_item_id
    UNIQUE (cart_id, menu_item_id)
);

CREATE TABLE IF NOT EXISTS orders (
  id           VARCHAR(36) PRIMARY KEY,
  cart_id      VARCHAR(36) NOT NULL,
  bot_id       VARCHAR(64) NOT NULL,
  session_id   VARCHAR(36) NOT NULL,
  total_scaled INTEGER NOT NULL,

  CONSTRAINT fk_orders_cart
    FOREIGN KEY (cart_id) REFERENCES cart (id)
    ON DELETE RESTRICT
);

-- Indexes because SQLAlchemy used index=True on bot_id/session_id
CREATE INDEX IF NOT EXISTS ix_orders_bot_id ON orders (bot_id);
CREATE INDEX IF NOT EXISTS ix_orders_session_id ON orders (session_id);

CREATE TABLE IF NOT EXISTS order_item (
  id                 VARCHAR(36) PRIMARY KEY,
  order_id           VARCHAR(36) NOT NULL,
  menu_item_id       VARCHAR(64) NOT NULL,
  name               VARCHAR(200) NOT NULL,
  quantity           INTEGER NOT NULL,
  unit_price_scaled  INTEGER NOT NULL,
  total_price_scaled INTEGER NOT NULL,

  CONSTRAINT fk_order_item_order
    FOREIGN KEY (order_id) REFERENCES orders (id)
    ON DELETE CASCADE
);

-- Helpful extra indexes (not required by your SQLAlchemy, but usually useful)
CREATE INDEX IF NOT EXISTS ix_cart_item_cart_id   ON cart_item (cart_id);
CREATE INDEX IF NOT EXISTS ix_order_item_order_id ON order_item (order_id);
CREATE INDEX IF NOT EXISTS ix_orders_cart_id      ON orders (cart_id);
