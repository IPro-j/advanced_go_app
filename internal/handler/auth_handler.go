package handler

import (
	"blog-api/internal/service"
	"context"
	"net/http"
)

type AuthHandler struct {
	userService *service.UserService
}

func NewAuthHandler(userService *service.UserService) *AuthHandler {
	return &AuthHandler{
		userService: userService,
	}
}

// Register обрабатывает запрос на регистрацию нового пользователя
// POST /api/register
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	// TODO: Реализовать обработку регистрации
	// Шаги:
	// 1. Проверить метод запроса (должен быть POST)
	// 2. Декодировать JSON тело в UserCreateRequest
	// 3. Вызвать userService.Register
	// 4. Обработать ошибки (ErrUserAlreadyExists -> 409 Conflict)
	// 5. Вернуть JSON ответ с токеном (201 Created)

	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

// Login обрабатывает запрос на вход пользователя
// POST /api/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// TODO: Реализовать обработку входа
	// Шаги:
	// 1. Проверить метод запроса (должен быть POST)
	// 2. Декодировать JSON тело в UserLoginRequest
	// 3. Вызвать userService.Login
	// 4. Обработать ошибки (ErrInvalidCredentials -> 401 Unauthorized)
	// 5. Вернуть JSON ответ с токеном (200 OK)

	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

// GetProfile возвращает профиль текущего пользователя (опционально)
// Этот метод не используется в эталонной реализации
func (h *AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	// TODO: Опционально - реализовать получение профиля
	// Этот эндпоинт не обязателен для базовой реализации

	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

// writeError отправляет JSON ответ с ошибкой
func writeError(w http.ResponseWriter, message string, statusCode int) {
	// TODO: Реализовать отправку ошибки в формате JSON
	// Создать структуру ErrorResponse и отправить как JSON

	http.Error(w, message, statusCode)
}

// getUserIDFromContext извлекает ID пользователя из контекста
func getUserIDFromContext(ctx context.Context) (int, bool) {
	// TODO: Извлечь userID из контекста
	// Ключ устанавливается в auth middleware

	return 0, false
}
