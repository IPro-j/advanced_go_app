-- Миграция для создания начальной схемы базы данных
-- TODO: Реализуйте создание таблиц для блог-платформы

-- Таблица пользователей
-- TODO: Создайте таблицу users со следующими полями:
-- - id (serial, primary key)
-- - username (varchar(50), unique, not null)
-- - email (varchar(255), unique, not null)
-- - password (varchar(255), not null) - для хешированного пароля
-- - created_at (timestamp, not null)
-- - updated_at (timestamp, not null)

-- Пример структуры:
-- CREATE TABLE IF NOT EXISTS users (
--     ...ваши поля здесь...
-- );

-- Таблица постов
-- TODO: Создайте таблицу posts со следующими полями:
-- - id (serial, primary key)
-- - title (varchar(200), not null)
-- - content (text, not null)
-- - author_id (integer, foreign key на users.id)
-- - created_at (timestamp, not null)
-- - updated_at (timestamp, not null)

-- Таблица комментариев
-- TODO: Создайте таблицу comments со следующими полями:
-- - id (serial, primary key)
-- - content (text, not null)
-- - post_id (integer, foreign key на posts.id)
-- - author_id (integer, foreign key на users.id)
-- - created_at (timestamp, not null)

-- Индексы
-- TODO: Создайте индексы для оптимизации запросов:
-- - Индекс на posts.author_id для быстрого поиска постов пользователя
-- - Индекс на comments.post_id для быстрого поиска комментариев к посту
-- - Индекс на posts.created_at для сортировки по дате

-- Подсказки:
-- 1. Используйте IF NOT EXISTS для избежания ошибок при повторном запуске
-- 2. Для foreign key используйте ON DELETE CASCADE для автоматического удаления связанных записей
-- 3. Для timestamp полей можно использовать DEFAULT CURRENT_TIMESTAMP
-- 4. Не забудьте про ограничения (constraints) для валидации данных
