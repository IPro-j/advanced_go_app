package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// Config содержит параметры подключения к PostgreSQL
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// NewPostgresDB создает новое подключение к PostgreSQL
func NewPostgresDB(cfg Config) (*sql.DB, error) {
	// TODO: Реализовать подключение к PostgreSQL
	// Шаги:
	// 1. Сформировать строку подключения (DSN) из параметров конфигурации
	// 2. Открыть соединение с БД используя sql.Open("postgres", dsn)
	// 3. Проверить соединение методом Ping()
	// 4. Настроить пул соединений (SetMaxOpenConns, SetMaxIdleConns)
	// 5. Вернуть подключение или ошибку

	return nil, fmt.Errorf("not implemented")
}

// Migrate выполняет миграции базы данных
func Migrate(db *sql.DB) error {
	// TODO: Реализовать применение миграций
	// Шаги:
	// 1. Создать таблицу users если не существует
	// 2. Создать таблицу posts если не существует
	// 3. Создать таблицу comments если не существует
	// 4. Создать необходимые индексы
	// 5. Вернуть ошибку если что-то пошло не так
	//
	// Структура таблиц:
	// - users: id, username, email, password_hash, created_at, updated_at
	// - posts: id, title, content, author_id, created_at, updated_at
	// - comments: id, content, post_id, author_id, created_at, updated_at

	queries := []string{
		// TODO: Добавить SQL запросы для создания таблиц
		// Пример:
		// `CREATE TABLE IF NOT EXISTS users (...)`,
		// `CREATE TABLE IF NOT EXISTS posts (...)`,
		// `CREATE TABLE IF NOT EXISTS comments (...)`,
		// `CREATE INDEX IF NOT EXISTS ...`,
	}

	// TODO: Выполнить каждый запрос в транзакции
	_ = queries // Удалить после реализации

	return fmt.Errorf("not implemented")
}

// CheckConnection проверяет соединение с базой данных
func CheckConnection(db *sql.DB) error {
	// TODO: Реализовать проверку соединения
	// Использовать db.Ping() для проверки

	return fmt.Errorf("not implemented")
}

// GetDSN формирует строку подключения к PostgreSQL
func GetDSN(cfg Config) string {
	// TODO: Сформировать DSN строку
	// Формат: "host=%s port=%d user=%s password=%s dbname=%s sslmode=%s"

	return ""
}

// Close закрывает соединение с базой данных
func Close(db *sql.DB) error {
	// TODO: Корректно закрыть соединение

	return fmt.Errorf("not implemented")
}

// TestConnection выполняет тестовый запрос к БД (опциональное задание)
func TestConnection(db *sql.DB) error {
	// TODO: Выполнить простой запрос для проверки работы БД
	// Например: SELECT 1

	return fmt.Errorf("not implemented")
}
