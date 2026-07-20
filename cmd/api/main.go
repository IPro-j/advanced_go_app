package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	// Загружаем конфигурацию из .env файла
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found, using environment variables")
	}

	// TODO: Загрузить конфигурацию из переменных окружения
	// cfg := loadConfig()

	// TODO: Подключиться к базе данных
	// - Создать database.Config из параметров конфигурации
	// - Вызвать database.NewPostgresDB
	// - Обработать ошибки подключения
	// - Не забыть defer db.Close()

	// TODO: Выполнить миграции базы данных
	// - Вызвать database.Migrate(db)

	// TODO: Инициализировать JWT менеджер
	// - Создать jwtManager через auth.NewJWTManager

	// TODO: Создать слои приложения
	// 1. Репозитории (передать db)
	// 2. Сервисы (передать репозитории и jwtManager)
	// 3. Хендлеры (передать сервисы)
	// 4. Middleware (передать необходимые зависимости)

	// Настраиваем маршруты
	router := chi.NewRouter()

	// TODO: Настроить middleware
	// - Добавить глобальные middleware (logging, recovery, CORS)

	// TODO: Настроить маршруты
	// Публичные эндпоинты:
	// - POST /api/register
	// - POST /api/login
	// - GET /api/posts
	// - GET /api/posts/{id}
	// - GET /api/posts/{id}/comments
	//
	// Защищенные эндпоинты (требуют JWT):
	// - POST /api/posts
	// - PUT /api/posts/{id}
	// - DELETE /api/posts/{id}
	// - POST /api/posts/{id}/comments

	// Health check эндпоинт
	router.Get("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","service":"blog-api"}`))
	})

	// TODO: Запустить HTTP сервер
	// - Сформировать адрес из конфигурации
	// - Вывести информацию о запуске
	// - Запустить сервер и обработать ошибки
}

// Config представляет конфигурацию приложения
type Config struct {
	// Server
	ServerHost string
	ServerPort int

	// Database
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	// JWT
	JWTSecret      string
	JWTExpiryHours int

	// Cache
	CacheTTLMinutes int
}

// loadConfig загружает конфигурацию из переменных окружения
func loadConfig() *Config {
	// TODO: Реализовать загрузку всех параметров конфигурации
	// Использовать вспомогательные функции getEnv и getEnvAsInt
	// Установить разумные значения по умолчанию

	return nil // Заменить на правильную реализацию
}

// getEnv получает значение переменной окружения или возвращает значение по умолчанию
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt получает значение переменной окружения как int или возвращает значение по умолчанию
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}
