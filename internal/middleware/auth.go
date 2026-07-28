package middleware

import (
	"blog-api/pkg/auth"
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

const (
	PostIDKey contextKey = "postID"
	// UserIDKey — ключ для хранения ID пользователя в контексте
	UserIDKey contextKey = "userID"
	// UserEmailKey — ключ для хранения email пользователя в контексте
	UserEmailKey contextKey = "userEmail"
	// UserNameKey — ключ для хранения username пользователя в контексте
	UserNameKey contextKey = "username"
)

// ErrorResponse — структура для единообразного возврата ошибок в JSON
type ErrorResponse struct {
	Error string `json:"error"`
}

// AuthMiddleware предоставляет middleware для JWT-аутентификации
type AuthMiddleware struct {
	jwtManager *auth.JWTManager
}

// NewAuthMiddleware создаёт новый экземпляр middleware аутентификации
func NewAuthMiddleware(jwtManager *auth.JWTManager) *AuthMiddleware {
	return &AuthMiddleware{
		jwtManager: jwtManager,
	}
}

// RequireAuth — middleware, требующий валидный JWT-токен.
// Если токена нет или он невалиден — возвращает 401.
func (m *AuthMiddleware) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := extractToken(r)
		if token == "" {
			writeJSONError(w, "missing authorization token", http.StatusUnauthorized)
			return
		}

		claims, err := m.jwtManager.ValidateToken(token)
		if err != nil {
			// Можно различать типы ошибок (истек, неверный формат и т.д.)
			writeJSONError(w, "invalid or expired token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, UserEmailKey, claims.Email)
		ctx = context.WithValue(ctx, UserNameKey, claims.Username)

		next(w, r.WithContext(ctx))
	}
}

// Если токен валиден — добавляет данные пользователя в контекст.
// Если токена нет или он невалиден — продолжает обработку как анонимный запрос.
func (m *AuthMiddleware) OptionalAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := extractToken(r)
		var ctx context.Context = r.Context()

		if token != "" {
			claims, err := m.jwtManager.ValidateToken(token)
			if err == nil {
				ctx = context.WithValue(ctx, UserIDKey, claims.UserID)
				ctx = context.WithValue(ctx, UserEmailKey, claims.Email)
				ctx = context.WithValue(ctx, UserNameKey, claims.Username)
			}
			// Если токен невалиден, просто продолжаем без добавления данных в контекст
		}

		next(w, r.WithContext(ctx))
	}
}

// extractToken извлекает JWT-токен из заголовка Authorization.
// Ожидаемый формат: "Bearer <token>"
func extractToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}

// GetUserIDFromContext извлекает ID пользователя из контекста
func GetUserIDFromContext(ctx context.Context) (int, bool) {
	val, ok := ctx.Value(UserIDKey).(int)
	return val, ok
}

// GetUserEmailFromContext извлекает email пользователя из контекста
func GetUserEmailFromContext(ctx context.Context) (string, bool) {
	val, ok := ctx.Value(UserEmailKey).(string)
	return val, ok
}

// GetUsernameFromContext извлекает username из контекста
func GetUsernameFromContext(ctx context.Context) (string, bool) {
	val, ok := ctx.Value(UserNameKey).(string)
	return val, ok
}

// writeJSONError отправляет ошибку в формате JSON
func writeJSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

// Chain позволяет объединить несколько middleware в цепочку.
func Chain(handler http.HandlerFunc, middlewares ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	wrapped := handler
	for i := len(middlewares) - 1; i >= 0; i-- {
		wrapped = middlewares[i](wrapped)
	}
	return wrapped
}

// AuthMiddlewareForChi возвращает middleware, совместимый с chi.Router.Use.
func (m *AuthMiddleware) AuthMiddlewareForChi() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Превращаем next (Handler) в HandlerFunc, чтобы можно было передать в RequireAuth
			nextHandlerFunc := http.HandlerFunc(next.ServeHTTP)
			wrapped := m.RequireAuth(nextHandlerFunc)
			wrapped.ServeHTTP(w, r)
		})
	}
}
