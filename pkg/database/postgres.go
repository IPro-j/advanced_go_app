package database

import (
	"database/sql"
	"fmt"
	"time"

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

// GetDSN формирует строку подключения к PostgreSQL.
func GetDSN(cfg Config) string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)
}

// NewPostgresDB создает новое подключение к PostgreSQL
func NewPostgresDB(cfg Config) (*sql.DB, error) {
	dsn := GetDSN(cfg)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection: %w", err)
	}

	// Проверка соединения сразу после открытия
	if err := CheckConnection(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Настройка пула соединений
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}

// Migrate выполняет миграции базы данных
func Migrate(db *sql.DB) error {
	queries := []string{
		// Таблица users
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(255) NOT NULL UNIQUE,
			email VARCHAR(255) NOT NULL UNIQUE,
			password VARCHAR(255) NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
		)`,

		// Таблица posts
		`CREATE TABLE IF NOT EXISTS posts (
			id SERIAL PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			content TEXT NOT NULL,
			author_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
		)`,

		// Таблица comments
		`CREATE TABLE IF NOT EXISTS comments (
			id SERIAL PRIMARY KEY,
			content TEXT NOT NULL,
			post_id INT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
			author_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
		)`,

		// Индексы для ускорения выборки
		`CREATE INDEX IF NOT EXISTS idx_posts_author_id ON posts(author_id)`,
		`CREATE INDEX IF NOT EXISTS idx_comments_post_id ON comments(post_id)`,
		`CREATE INDEX IF NOT EXISTS idx_comments_author_id ON comments(author_id)`,
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	for i, query := range queries {
		_, err := tx.Exec(query)
		if err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("migration failed at step %d: %w", i, err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit migration: %w", err)
	}

	return nil
}

// CheckConnection проверяет соединение с базой данных
func CheckConnection(db *sql.DB) error {
	if err := db.Ping(); err != nil {
		return fmt.Errorf("ping failed: %w", err)
	}
	return nil
}

// Close закрывает соединение с базой данных
func Close(db *sql.DB) error {
	if db == nil {
		return nil // Безопасный вызов, если db nil
	}
	err := db.Close()
	if err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}
	return nil
}

// TestConnection выполняет тестовый запрос к БД (SELECT 1) для проверки работоспособности.
func TestConnection(db *sql.DB) error {
	var result int
	row := db.QueryRow("SELECT 1")
	err := row.Scan(&result)
	if err != nil {
		return fmt.Errorf("test query failed: %w", err)
	}
	if result != 1 {
		return fmt.Errorf("unexpected test query result: %d, expected 1", result)
	}
	return nil
}
