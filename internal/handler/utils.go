package handler

import (
	"blog-api/internal/middleware"
	"context"
	"encoding/json"
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

func getPostIDFromContext(ctx context.Context) (int, bool) {
	val, ok := ctx.Value(middleware.PostIDKey).(int)
	return val, ok
}
