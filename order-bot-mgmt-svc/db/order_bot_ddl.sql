-- DDL for order-bot-svc entities.

create table order_bot.published_menu
(
    id         text not null
        primary key,
    bot_id     text not null,
    created_at timestamp,
    updated_at timestamp
);

alter table order_bot.published_menu
    owner to melkey;

create table order_bot.published_menu_item
(
    id             text not null
        primary key,
    menu_id        text not null
        references order_bot.published_menu,
    menu_item_name text not null,
    price          double precision,
    created_at     timestamp,
    updated_at     timestamp
);

alter table order_bot.published_menu_item
    owner to melkey;

create index idx_menu_item_menu_id
    on order_bot.published_menu_item (menu_id);

create table order_bot.cart
(
    id           varchar(36) not null
        primary key,
    session_id   varchar(36) not null,
    status       varchar(6)  not null,
    total_scaled integer     not null,
    closed_at    timestamp,
    created_at   timestamp,
    updated_at   timestamp
);

alter table order_bot.cart
    owner to melkey;

create unique index ix_cart_session_id
    on order_bot.cart (session_id);

create table order_bot.cart_item
(
    id                 varchar(36)  not null
        primary key,
    cart_id            varchar(36)  not null
        references order_bot.cart,
    menu_item_id       varchar(64)  not null,
    name               varchar(200) not null,
    quantity           integer      not null,
    unit_price_scaled  integer      not null,
    total_price_scaled integer      not null,
    created_at         timestamp,
    updated_at         timestamp,
    constraint menu_item_id
        unique (cart_id, menu_item_id)
);

alter table order_bot.cart_item
    owner to melkey;

create index ix_cart_item_cart_id
    on order_bot.cart_item (cart_id);

create table order_bot.orders
(
    id           varchar(36) not null
        primary key,
    cart_id      varchar(36) not null
        references order_bot.cart,
    session_id   varchar(36) not null,
    total_scaled integer     not null,
    bot_id       varchar(36) not null,
    created_at   timestamp,
    updated_at   timestamp
);

alter table order_bot.orders
    owner to melkey;

create index ix_orders_session_id
    on order_bot.orders (session_id);

create index ix_orders_cart_id
    on order_bot.orders (cart_id);

create table order_bot.order_item
(
    id                 varchar(36)  not null
        primary key,
    order_id           varchar(36)  not null
        references order_bot.orders,
    menu_item_id       varchar(64)  not null,
    name               varchar(200) not null,
    quantity           integer      not null,
    unit_price_scaled  integer      not null,
    total_price_scaled integer      not null,
    created_at         timestamp,
    updated_at         timestamp
);

alter table order_bot.order_item
    owner to melkey;

create index ix_order_item_order_id
    on order_bot.order_item (order_id);

