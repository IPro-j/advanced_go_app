package handler

import (
	"blog-api/internal/model"
	"blog-api/internal/service"
	"blog-api/pkg/apperr"
	"blog-api/pkg/auth"
	"strings"

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

	token, err := h.userService.Register(r.Context(), &req)

	if err != nil {

		if isUsernameError(err) {
			log.Printf("username invalid, username=%q, err=%v", req.Username, err)
			writeError(w, "username length is invalid", http.StatusBadRequest)
			return
		}

		if isPasswordError(err) {
			log.Printf("password invalid, email=%q, err=%v", req.Email, err)
			writeError(w, "password does not meet requirements", http.StatusBadRequest)
			return
		}

		if errors.Is(err, apperr.ErrInvalidEmail) {
			log.Printf("invalid email, email=%q, err=%v", req.Email, err)
			writeError(w, "email is invalid", http.StatusBadRequest)
			return
		}

		if errors.Is(err, apperr.ErrUserAlreadyExists) || errors.Is(err, apperr.ErrEmailAlreadyExists) {
			log.Printf("user/email already exists, username=%q email=%q, err=%v", req.Username, req.Email, err)
			writeError(w, "user/email already exists", http.StatusConflict)
			return
		}

		// Все остальные ошибки (БД, токены, непредвиденные)
		log.Printf("internal error, username=%q, email=%q, path=%s, err=%v",
			req.Username, req.Email, r.URL.Path, err)
		writeError(w, "registration failed due to an internal error", http.StatusInternalServerError)
		return

	}

	log.Printf("user registered, email=%q, username=%q", req.Email, req.Username)

	accessToken, accessExp := token.Token, token.ExpiresAt

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

	token, err := h.userService.Login(r.Context(), &req)
	if err != nil {

		if errors.Is(err, apperr.ErrInvalidCredentials) || errors.Is(err, apperr.ErrUserNotFound) {
			log.Printf("login failed, email=%q, err=%v", req.Email, err)
			writeError(w, "login failed,", http.StatusUnauthorized)
			return
		}

		log.Printf("login failed, email=%q, err=%v", req.Email, err)
		writeError(w, "login failed", http.StatusInternalServerError)
		return
	}

	log.Printf("user login is successful, email=%q", req.Email)

	accessToken, accessExp := token.Token, token.ExpiresAt

	resp := AuthTokenResponse{
		AccessToken: accessToken,
		ExpiresAt:   accessExp.Unix(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}
