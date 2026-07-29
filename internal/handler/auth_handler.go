package handler

import (
	//"blog-api/internal/middleware"
	"blog-api/internal/model"
	"blog-api/internal/service"
	"blog-api/pkg/apperr"
	"blog-api/pkg/auth"
	"strings"

	//"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

// AuthTokenResponse — ответ с токеном
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

	// нормализация данных
	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)

	//Простая валидация
	if req.Username == "" || req.Email == "" || req.Password == "" {
		writeError(w, "username, email and password are required", http.StatusBadRequest)
		return
	}

	req.Email = strings.ToLower(req.Email)

	user, err := h.userService.Register(r.Context(), &req)

	if errors.Is(err, apperr.ErrInvalidUsername) {
		writeError(w, "username is invalid", http.StatusBadRequest)
		return
	}
	if errors.Is(err, apperr.ErrInvalidEmail) {
		writeError(w, "email is invalid", http.StatusBadRequest)
		return
	}
	if errors.Is(err, apperr.ErrWeakPassword) {
		writeError(w, "password is too weak", http.StatusBadRequest)
		return
	}

	if err != nil {
		// Проверяем конкретную ошибку дубликата
		if errors.Is(err, apperr.ErrUserAlreadyExists) {
			writeError(w, "user already exists", http.StatusConflict)
			return

		}

		if errors.Is(err, apperr.ErrEmailAlreadyExists) {
			writeError(w, "email already exists", http.StatusConflict)
			return

		}

		log.Printf("registration failed: %v", err)
		writeError(w, "registration failed", http.StatusInternalServerError)
		return
	}

	accessToken, accessExp, err := h.jwtManager.GenerateToken(user.User.ID, user.User.Email, user.User.Username)
	if err != nil {
		writeError(w, "failed to generate access token", http.StatusInternalServerError)
		return
	}

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

	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)

	if req.Email == "" || req.Password == "" {
		writeError(w, "email and password are required", http.StatusBadRequest)
		return
	}

	user, err := h.userService.Login(r.Context(), &req)
	if err != nil {
		//ошибка общая для email/password
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
