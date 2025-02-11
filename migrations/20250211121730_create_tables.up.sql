-- Создание таблицы пользователей (сотрудников)
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    coin_balance INTEGER NOT NULL DEFAULT 1000 CHECK (coin_balance >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Создание таблицы товаров мерча (фиксированный справочник)
CREATE TABLE merch_items (
    name VARCHAR(50) PRIMARY KEY,
    price INTEGER NOT NULL CHECK (price > 0)
);

-- Инициализация таблицы merch_items (справочник мерча)
INSERT INTO merch_items (name, price) VALUES
    ('t-shirt', 80),
    ('cup', 20),
    ('book', 50),
    ('pen', 10),
    ('powerbank', 200),
    ('hoody', 300),
    ('umbrella', 200),
    ('socks', 10),
    ('wallet', 50),
    ('pink-hoody', 500);

-- Создание таблицы инвентаря, где фиксируются купленные товары
CREATE TABLE inventory (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    merch_name VARCHAR(50) NOT NULL REFERENCES merch_items(name),
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    UNIQUE (user_id, merch_name)
);

-- Создание таблицы для записей перевода монет между пользователями
CREATE TABLE coin_transfers (
    id SERIAL PRIMARY KEY,
    from_user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    to_user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    amount INTEGER NOT NULL CHECK (amount > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Создание таблицы для записей покупки мерча (монеты списываются при покупке)
CREATE TABLE merch_purchases (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    merch_name VARCHAR(50) NOT NULL REFERENCES merch_items(name),
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    total_price INTEGER NOT NULL CHECK (total_price >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
