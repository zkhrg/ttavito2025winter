-- Активируем расширение для работы с UUID
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Создаём таблицу products с уникальным ограничением на product_name
CREATE TABLE IF NOT EXISTS products (
   product_name VARCHAR(100) PRIMARY KEY NOT NULL UNIQUE,  -- Добавлено уникальное ограничение
   price INT NOT NULL CHECK (price >= 0)
);

-- Создаём таблицу users
CREATE TABLE IF NOT EXISTS users (
   username VARCHAR(100) PRIMARY KEY NOT NULL UNIQUE,
   user_password VARCHAR(255) NOT NULL, -- длина может быть больше, если хеш пароля длиннее
   balance INT NOT NULL DEFAULT 1000 CHECK (balance >= 0) -- баланс пользователя
);

INSERT INTO users (username, user_password) VALUES
   ('test', 'test')
ON CONFLICT (username) DO NOTHING; 

-- Заполняем таблицу данными
INSERT INTO products (product_name, price) VALUES
   ('t-shirt', 80),
   ('cup', 20),
   ('book', 50),
   ('pen', 10),
   ('powerbank', 200),
   ('hoody', 300),
   ('umbrella', 200),
   ('socks', 10),
   ('wallet', 50),
   ('pink-hoody', 500)
ON CONFLICT (product_name) DO NOTHING; 

-- Создаём таблицу transfers (переводы)
CREATE TABLE IF NOT EXISTS transfers (
   id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
   sender_username VARCHAR(100) NOT NULL, -- имя отправителя
   receiver_username VARCHAR(100) NOT NULL, -- имя получателя
   amount INT NOT NULL CHECK (amount > 0), -- сумма перевода
   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   FOREIGN KEY (sender_username) REFERENCES users (username), -- связь с таблицей пользователей
   FOREIGN KEY (receiver_username) REFERENCES users (username) -- связь с таблицей пользователей
);

-- Создаём таблицу purchases (покупки)
CREATE TABLE IF NOT EXISTS purchases (
   id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
   username VARCHAR(100) NOT NULL, -- имя покупателя
   product_name VARCHAR(100) NOT NULL, -- ID продукта
   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   FOREIGN KEY (username) REFERENCES users (username), -- связь с таблицей пользователей
   FOREIGN KEY (product_name) REFERENCES products (product_name) -- связь с таблицей продуктов
);

-- Индексы для ускорения поиска
CREATE INDEX IF NOT EXISTS idx_transfers_sender ON transfers(sender_username);
CREATE INDEX IF NOT EXISTS idx_transfers_receiver ON transfers(receiver_username);
CREATE INDEX IF NOT EXISTS idx_purchases_username ON purchases(username);
CREATE INDEX IF NOT EXISTS idx_purchases_product ON purchases(product_name);
