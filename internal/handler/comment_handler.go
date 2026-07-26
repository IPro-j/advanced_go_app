package handler

import (
	"blog-api/internal/model"
	"blog-api/internal/service"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type CommentHandler struct {
	commentService *service.CommentService
}

func NewCommentHandler(commentService *service.CommentService) *CommentHandler {
	return &CommentHandler{
		commentService: commentService,
	}
}

// writeError отправляет JSON-ответ с ошибкой
//func writeError(w http.ResponseWriter, message string, statusCode int) {
//w.Header().Set("Content-Type", "application/json")
//w.WriteHeader(statusCode)
//	_ = json.NewEncoder(w).Encode(ErrorResponse{Error: message})
//}

// getUserIDFromContext извлекает userID из контекста
//func getUserIDFromContext(ctx context.Context) (int, bool) {
//	val, ok := ctx.Value(middleware.UserIDKey).(int)
//	return val, ok
//}

// extractIDFromPath извлекает ID из пути вида /api/comments/123
func extractIDFromPath(path, prefix string) string {
	if !strings.HasPrefix(path, prefix) {
		return ""
	}

	// Убираем префикс, остаётся "/123" или "123"
	suffix := strings.TrimPrefix(path, prefix)

	// Убираем ведущий слэш, если он есть (TrimPrefix ничего не сделает, если слэша нет)
	suffix = strings.TrimPrefix(suffix, "/")
	fmt.Println(suffix)
	return suffix
}

// extractPostIDFromCommentsPath извлекает postID из пути /api/posts/123/comments
func extractPostIDFromCommentsPath(path string) string {
	// Ожидаемый формат: /api/posts/{id}/comments
	prefix := "/api/protected/posts/"
	if !strings.HasPrefix(path, prefix) {
		return ""
	}
	suffix := strings.TrimPrefix(path, prefix)

	// Ищем позицию "/comments"
	idx := strings.Index(suffix, "/comments")
	if idx == -1 {
		return ""
	}

	postIDStr := suffix[:idx]
	return postIDStr
}

// Create обрабатывает создание нового комментария
// POST /api/comments
func (h *CommentHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 1. Парсим postID из URL (не из контекста!)
	postIDStr := extractPostIDFromCommentsPath(r.URL.Path)
	if postIDStr == "" {
		writeError(w, "invalid URL format: expected /api/posts/{postID}/comments", http.StatusBadRequest)
		return
	}
	postID, err := strconv.Atoi(postIDStr)
	if err != nil || postID <= 0 {
		writeError(w, "invalid post ID", http.StatusBadRequest)
		return
	}

	userID, ok := getUserIDFromContext(r.Context())
	if !ok {
		writeError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req model.CommentCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	comment, err := h.commentService.Create(r.Context(), postID, userID, &req)
	if err != nil {
		switch err {
		case service.ErrPostNotExists:
			writeError(w, "post not found", http.StatusNotFound)
		default:
			writeError(w, "failed to create comment", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(comment)
}

// GetByID возвращает комментарий по ID
// GET /api/comments/{id}
func (h *CommentHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := extractIDFromPath(r.URL.Path, "/api/comments/")
	if idStr == "" {
		writeError(w, "invalid comment ID path", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		writeError(w, "invalid comment ID", http.StatusBadRequest)
		return
	}

	comment, err := h.commentService.GetByID(r.Context(), id)
	if err != nil {
		if err == service.ErrCommentNotFound {
			writeError(w, "comment not found", http.StatusNotFound)
		} else {
			writeError(w, "failed to get comment", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(comment)
}

// GetByPost возвращает комментарии к посту с пагинацией
// GET /api/posts/{id}/comments?limit=20&offset=0
func (h *CommentHandler) GetByPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	postIDStr := extractPostIDFromCommentsPath(r.URL.Path)
	if postIDStr == "" {
		writeError(w, "invalid post ID in path", http.StatusBadRequest)
		return
	}

	postID, err := strconv.Atoi(postIDStr)
	if err != nil || postID <= 0 {
		writeError(w, "invalid post ID", http.StatusBadRequest)
		return
	}

	query := r.URL.Query()
	limit, _ := strconv.Atoi(query.Get("limit"))
	if limit <= 0 {
		limit = 20 // default
	}
	offset, _ := strconv.Atoi(query.Get("offset"))
	if offset < 0 {
		offset = 0
	}

	comments, total, err := h.commentService.GetByPost(r.Context(), postID, limit, offset)
	if err != nil {
		if err == service.ErrPostNotExists {
			writeError(w, "post not found", http.StatusNotFound)
		} else {
			writeError(w, "failed to get comments", http.StatusInternalServerError)
		}
		return
	}

	type CommentsResponse struct {
		Comments []*model.Comment `json:"comments"`
		Total    int              `json:"total"`
		Limit    int              `json:"limit"`
		Offset   int              `json:"offset"`
		PostID   int              `json:"post_id"`
	}

	resp := CommentsResponse{
		Comments: comments,
		Total:    total,
		Limit:    limit,
		Offset:   offset,
		PostID:   postID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

// Update обновляет комментарий
// PUT /api/comments/{id}
func (h *CommentHandler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := getUserIDFromContext(r.Context())
	if !ok {
		writeError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := extractIDFromPath(r.URL.Path, "/api/protected/comments/")
	if idStr == "" {
		writeError(w, "invalid comment ID path", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		writeError(w, "invalid comment ID", http.StatusBadRequest)
		return
	}

	var req model.CommentUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	comment, err := h.commentService.Update(r.Context(), id, userID, &req)
	if err != nil {
		switch err {
		case service.ErrCommentNotFound:
			writeError(w, "comment not found", http.StatusNotFound)
		case service.ErrForbidden:
			writeError(w, "you can only update your own comments", http.StatusForbidden)
		default:
			writeError(w, "failed to update comment", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(comment)
}
