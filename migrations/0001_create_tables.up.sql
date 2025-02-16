-- Создание таблицы users
CREATE TABLE IF NOT EXISTS users
(
    id         SERIAL PRIMARY KEY,
    username   VARCHAR(255) NOT NULL UNIQUE,
    password   TEXT NOT NULL,   -- HASH
    coins      INTEGER      NOT NULL DEFAULT 1000 CHECK (coins >= 0),
    created_at TIMESTAMP    NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP    NOT NULL DEFAULT NOW()
);

-- Создание таблицы store
CREATE TABLE IF NOT EXISTS store
(
    slug  VARCHAR(255) PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    price INTEGER      NOT NULL DEFAULT 1 CHECK (price >= 0)
);

-- Создание таблицы inventory
CREATE TABLE IF NOT EXISTS inventory
(
    user_id    INTEGER      NOT NULL,
    item_slug  VARCHAR(255) NOT NULL,
    quantity   INTEGER      NOT NULL DEFAULT 1 CHECK (quantity >= 0),
    created_at TIMESTAMP    NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP    NOT NULL DEFAULT NOW(),
    CONSTRAINT unique_inventory UNIQUE (user_id, item_slug),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (item_slug) REFERENCES store (slug) ON DELETE CASCADE
);

-- Создание таблицы transactions
CREATE TABLE IF NOT EXISTS transactions
(
    id          SERIAL PRIMARY KEY,
    sender_id   INTEGER   NOT NULL,
    receiver_id INTEGER   NOT NULL,
    coins       INTEGER   NOT NULL CHECK (coins >= 1),
    created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT check_sender_receiver CHECK (sender_id <> receiver_id),
    FOREIGN KEY (sender_id) REFERENCES users (id) ON DELETE RESTRICT,
    FOREIGN KEY (receiver_id) REFERENCES users (id) ON DELETE RESTRICT
);

-- Дополнительные индексы (если необходимо)
CREATE INDEX IF NOT EXISTS idx_users_username ON users (username);
CREATE INDEX IF NOT EXISTS idx_transactions_sender ON transactions (sender_id);
CREATE INDEX IF NOT EXISTS idx_transactions_receiver ON transactions (receiver_id);
CREATE INDEX IF NOT EXISTS idx_inventory_user ON inventory (user_id);
