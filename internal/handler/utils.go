package handler

import (
	"blog-api/internal/middleware"
	"blog-api/pkg/apperr"
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func writeError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

func getUserIDFromContext(ctx context.Context) (int, bool) {
	val, ok := ctx.Value(middleware.UserIDKey).(int)
	return val, ok
}

func isUsernameError(err error) bool {
	return errors.Is(err, apperr.ErrUsernameTooShort) ||
		errors.Is(err, apperr.ErrUsernameTooLong)
}

func isPasswordError(err error) bool {
	return errors.Is(err, apperr.ErrPasswordTooShort) ||
		errors.Is(err, apperr.ErrPasswordTooLong) ||
		errors.Is(err, apperr.ErrNoDigit) ||
		errors.Is(err, apperr.ErrNoUpper) ||
		errors.Is(err, apperr.ErrNoLower) ||
		errors.Is(err, apperr.ErrNoSpecialChar)
}
