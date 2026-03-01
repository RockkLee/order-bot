-- DDL for order-bot-mgmt-svc entities.

create table order_bot_mgmt.bot
(
    id         text not null
        primary key,
    bot_name   text not null,
    created_at timestamp,
    updated_at timestamp
);

alter table order_bot_mgmt.bot
    owner to melkey;

create table order_bot_mgmt.menu
(
    id         text not null
        primary key,
    bot_id     text not null
        references order_bot_mgmt.bot,
    created_at timestamp,
    updated_at timestamp
);

alter table order_bot_mgmt.menu
    owner to melkey;

create index idx_menu_bot_id
    on order_bot_mgmt.menu (bot_id);

create table order_bot_mgmt.menu_item
(
    id             text not null
        primary key,
    menu_id        text not null
        references order_bot_mgmt.menu,
    menu_item_name text not null,
    price          double precision,
    created_at     timestamp,
    updated_at     timestamp
);

alter table order_bot_mgmt.menu_item
    owner to melkey;

create index idx_menu_item_menu_id
    on order_bot_mgmt.menu_item (menu_id);

create table order_bot_mgmt.users
(
    id            text not null
        primary key,
    email         text not null
        unique,
    password_hash text not null,
    access_token  text not null,
    refresh_token text not null,
    created_at    timestamp,
    updated_at    timestamp
);

alter table order_bot_mgmt.users
    owner to melkey;

create table order_bot_mgmt.user_bot
(
    id         text not null
        primary key,
    user_id    text not null
        references order_bot_mgmt.users,
    bot_id     text not null
        references order_bot_mgmt.bot,
    created_at timestamp,
    updated_at timestamp
);

alter table order_bot_mgmt.user_bot
    owner to melkey;

create index idx_user_bot_user_id
    on order_bot_mgmt.user_bot (user_id);

create index idx_user_bot_bot_id
    on order_bot_mgmt.user_bot (bot_id);

