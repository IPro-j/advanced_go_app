package handler

import (
	"blog-api/internal/middleware"
	"blog-api/internal/model"
	"blog-api/internal/service"
	"blog-api/pkg/auth"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	//	"github.com/gorilla/mux" // если используешь mux; если нет — можно убрать
)

// ErrorResponse — структура для единообразного формата ошибок
type ErrorResponse struct {
	Error string `json:"error"`
}

// AuthTokenResponse — ответ с токенами (access + refresh)
type AuthTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresAt   int64  `json:"expires_at"` // Unix timestamp
}

type AuthHandler struct {
	userService *service.UserService
	jwtManager  *auth.JWTManager
}

func NewAuthHandler(userService *service.UserService, jwtManager *auth.JWTManager) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		jwtManager:  jwtManager,
	}
}

// Register обрабатывает запрос на регистрацию нового пользователя
// POST /api/register
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req model.UserCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	// Простая валидация (можно вынести в отдельный валидатор)
	if req.Username == "" || req.Email == "" || req.Password == "" {
		writeError(w, "username, email and password are required", http.StatusBadRequest)
		return
	}

	user, err := h.userService.Register(r.Context(), &req)
	if err != nil {
		// Проверяем конкретную ошибку дубликата
		if errors.Is(err, service.ErrUserAlreadyExists) {
			writeError(w, "user already exists", http.StatusConflict)
			return
		}

		// Для остальных ошибок лучше логировать детали (но не отдавать клиенту)
		log.Printf("registration failed: %v", err)
		writeError(w, "registration failed", http.StatusInternalServerError)
		return
	}

	accessToken, accessExp, err := h.jwtManager.GenerateToken(user.User.ID, user.User.Email, user.User.Username)
	if err != nil {
		writeError(w, "failed to generate access token", http.StatusInternalServerError)
		return
	}

	//refreshToken, err := h.jwtManager.GenerateToken(user.User.ID, user.User.Email, user.User.Username)
	//if err != nil {
	//	writeError(w, "failed to generate refresh token", http.StatusInternalServerError)
	//	return
	//}

	resp := AuthTokenResponse{
		AccessToken: accessToken,
		ExpiresAt:   accessExp.Unix(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)
}

// Login обрабатывает запрос на вход пользователя
// POST /api/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req model.UserLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		writeError(w, "email and password are required", http.StatusBadRequest)
		return
	}

	user, err := h.userService.Login(r.Context(), &req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			writeError(w, "invalid credentials", http.StatusUnauthorized)
			return
		}
		writeError(w, "login failed", http.StatusInternalServerError)
		return
	}

	accessToken, exp, err := h.jwtManager.GenerateToken(user.User.ID, user.User.Email, user.User.Username)
	if err != nil {
		writeError(w, "failed to generate access token", http.StatusInternalServerError)
		return
	}

	resp := AuthTokenResponse{
		AccessToken: accessToken,
		ExpiresAt:   exp.Unix(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

// GetProfile возвращает профиль текущего пользователя
// GET /api/profile (пример маршрута)
func (h *AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	// Предполагается, что AuthMiddleware уже положил userID в контекст
	userID, ok := getUserIDFromContext(r.Context())
	if !ok {
		writeError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.userService.GetByID(r.Context(), userID)
	if err != nil {
		// Если пользователь вдруг не найден (например, удалён)
		writeError(w, "user not found", http.StatusNotFound)
		return
	}

	profile := map[string]interface{}{
		"id":         user.ID,
		"username":   user.Username,
		"email":      user.Email,
		"created_at": user.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(profile)
}

// writeError отправляет JSON-ответ с ошибкой
func writeError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

// getUserIDFromContext извлекает userID из контекста
func getUserIDFromContext(ctx context.Context) (int, bool) {
	val, ok := ctx.Value(middleware.UserIDKey).(int)
	return val, ok
}

// getUserIDFromContext извлекает userID из контекста
func getPostIDFromContext(ctx context.Context) (int, bool) {
	val, ok := ctx.Value(middleware.PostIDKey).(int)
	return val, ok
}
