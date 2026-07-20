package middleware

import (
	"blog-api/pkg/auth"
	"context"
	"net/http"
)

// contextKey is a custom type for context keys to avoid collisions
type contextKey string

const (
	// UserIDKey is the key for storing user ID in context
	UserIDKey contextKey = "userID"
	// UserEmailKey is the key for storing user email in context
	UserEmailKey contextKey = "userEmail"
	// UserNameKey is the key for storing username in context
	UserNameKey contextKey = "username"
)

// AuthMiddleware provides JWT authentication
type AuthMiddleware struct {
	jwtManager *auth.JWTManager
}

// NewAuthMiddleware creates a new auth middleware instance
func NewAuthMiddleware(jwtManager *auth.JWTManager) *AuthMiddleware {
	return &AuthMiddleware{
		jwtManager: jwtManager,
	}
}

// RequireAuth is a middleware that requires valid JWT token
func (m *AuthMiddleware) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Реализовать проверку JWT токена
		// Шаги:
		// 1. Извлечь токен из заголовка Authorization (Bearer токен)
		// 2. Валидировать токен через jwtManager
		// 3. Обработать ошибки валидации (истек, невалидный и т.д.)
		// 4. Добавить данные пользователя в контекст (UserIDKey, UserEmailKey, UserNameKey)
		// 5. Передать управление следующему handler

		// Временная заглушка - удалить после реализации
		http.Error(w, "Authentication not implemented", http.StatusNotImplemented)
	}
}

// OptionalAuth is a middleware that extracts JWT token if present, but doesn't require it
func (m *AuthMiddleware) OptionalAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Реализовать опциональную проверку JWT токена
		// Шаги:
		// 1. Попытаться извлечь токен из заголовка
		// 2. Если токен есть - валидировать его
		// 3. Если токен валидный - добавить данные в контекст
		// 4. Если токена нет или он невалидный - продолжить как анонимный
		// 5. В любом случае передать управление следующему handler

		// Временная реализация
		next(w, r)
	}
}

// extractToken извлекает JWT токен из заголовка Authorization
func extractToken(r *http.Request) string {
	// TODO: Извлечь JWT токен из заголовка Authorization
	// Формат: "Bearer <token>"

	return ""
}

// GetUserIDFromContext извлекает ID пользователя из контекста
func GetUserIDFromContext(ctx context.Context) (int, bool) {
	// TODO: Извлечь userID из контекста (ключ UserIDKey)

	return 0, false
}

// GetUserEmailFromContext извлекает email пользователя из контекста
func GetUserEmailFromContext(ctx context.Context) (string, bool) {
	// TODO: Извлечь email из контекста (ключ UserEmailKey)

	return "", false
}

// GetUsernameFromContext извлекает username из контекста
func GetUsernameFromContext(ctx context.Context) (string, bool) {
	// TODO: Извлечь username из контекста (ключ UserNameKey)

	return "", false
}

// writeJSONError отправляет ошибку в формате JSON
func writeJSONError(w http.ResponseWriter, message string, statusCode int) {
	// TODO: Отправить ошибку в формате JSON
	// Создать структуру ErrorResponse и отправить как JSON

	// Временная реализация
	http.Error(w, message, statusCode)
}

// Вспомогательные функции для упрощения использования middleware

// Chain позволяет объединить несколько middleware в цепочку
func Chain(handler http.HandlerFunc, middlewares ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	// TODO: Реализовать объединение middleware в цепочку
	// Применить их в правильном порядке

	return handler
}
