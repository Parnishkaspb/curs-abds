-- Создание БД
CREATE DATABASE IF NOT EXISTS abds;

-- Создание пользователя
CREATE USER IF NOT EXISTS user
IDENTIFIED WITH sha256_password BY 'SOME_STRONG_PASSWORD';

-- Права
GRANT ALL ON abds.* TO user;