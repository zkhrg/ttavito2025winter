-- Активируем расширение для работы с UUID
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Создаём ENUM для пола
DO $$
BEGIN
   IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'gender_enum') THEN
      CREATE TYPE gender_enum AS ENUM ('male', 'female');
   END IF;
END
$$;

-- Создаём таблицу user
CREATE TABLE IF NOT EXISTS users (
   id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
   first_name VARCHAR(100) NOT NULL,
   last_name VARCHAR(100) NOT NULL,
   birth_date DATE NOT NULL,
   gender gender_enum NOT NULL
);