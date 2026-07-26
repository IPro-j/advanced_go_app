package handler

import (
	"blog-api/internal/model"
	"blog-api/internal/service"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type PostHandler struct {
	postService *service.PostService
}

func NewPostHandler(postService *service.PostService) *PostHandler {
	return &PostHandler{
		postService: postService,
	}
}

// writeError отправляет JSON-ответ с ошибкой
//func writeError(w http.ResponseWriter, message string, statusCode int) {
//	w.Header().Set("Content-Type", "application/json")
//w.WriteHeader(statusCode)
//_ = json.NewEncoder(w).Encode(ErrorResponse{Error: message})
//}

// getUserIDFromContext извлекает userID из контекста
//func getUserIDFromContext(ctx context.Context) (int, bool) {
//	val, ok := ctx.Value(middleware.UserIDKey).(int)
//	return val, ok
//}

// extractIDFromPath извлекает ID из пути URL
//func extractIDFromPath(path, prefix string) string {
//	if !strings.HasPrefix(path, prefix) {
//		return ""
//	}
//	suffix := strings.TrimPrefix(path, prefix)
//	if strings.HasPrefix(suffix, "/") {
//		suffix = strings.TrimPrefix(suffix, "/")
//	}
// Если дальше есть ещё слеш (например, /api/posts/123/comments), обрезаем после первого числа
//	idx := strings.Index(suffix, "/")
//	if idx != -1 {
//		suffix = suffix[:idx]
//	}
//	return suffix
//}

// Create обрабатывает создание нового поста
func (h *PostHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := getUserIDFromContext(r.Context())
	if !ok {
		writeError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req model.PostCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	post, err := h.postService.Create(r.Context(), userID, &req)
	if err != nil {
		log.Printf("Create post failed: %v", err)
		// Если в сервисе есть специфичные ошибки — можно добавить switch
		writeError(w, "failed to create post", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(post)
}

// GetByID возвращает пост по ID
func (h *PostHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := extractIDFromPath(r.URL.Path, "/api/posts/")
	if idStr == "" {
		writeError(w, "invalid post ID path", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		writeError(w, "invalid post ID", http.StatusBadRequest)
		return
	}

	post, err := h.postService.GetByID(r.Context(), id)
	if err != nil {
		if err == service.ErrPostNotFound {
			writeError(w, "post not found", http.StatusNotFound)
		} else {
			writeError(w, "failed to get post", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(post)
}

// GetAll возвращает список постов с пагинацией
func (h *PostHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query()
	limit, _ := strconv.Atoi(query.Get("limit"))
	if limit <= 0 {
		limit = 10
	}
	offset, _ := strconv.Atoi(query.Get("offset"))
	if offset < 0 {
		offset = 0
	}

	posts, total, err := h.postService.GetAll(r.Context(), limit, offset)
	if err != nil {
		writeError(w, "failed to get posts", http.StatusInternalServerError)
		return
	}

	type PostsResponse struct {
		Posts  []*model.Post `json:"posts"`
		Total  int           `json:"total"`
		Limit  int           `json:"limit"`
		Offset int           `json:"offset"`
	}

	resp := PostsResponse{
		Posts:  posts,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

// Update обновляет пост
func (h *PostHandler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := getUserIDFromContext(r.Context())
	if !ok {
		writeError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := extractIDFromPath(r.URL.Path, "/api/protected/posts/")
	if idStr == "" {
		writeError(w, "invalid post ID path", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		writeError(w, "invalid post ID", http.StatusBadRequest)
		return
	}

	var req model.PostUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	post, err := h.postService.Update(r.Context(), id, userID, &req)
	if err != nil {
		switch err {
		case service.ErrPostNotFound:
			writeError(w, "post not found", http.StatusNotFound)
		case service.ErrForbidden:
			writeError(w, "you can only update your own posts", http.StatusForbidden)
		default:
			writeError(w, "failed to update post", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(post)
}

// Delete удаляет пост
func (h *PostHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := getUserIDFromContext(r.Context())
	if !ok {
		writeError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := extractIDFromPath(r.URL.Path, "/api/protected/posts/")
	if idStr == "" {
		writeError(w, "invalid post ID path", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		writeError(w, "invalid post ID", http.StatusBadRequest)
		return
	}

	err = h.postService.Delete(r.Context(), id, userID)
	if err != nil {
		switch err {
		case service.ErrPostNotFound:
			writeError(w, "post not found", http.StatusNotFound)
		case service.ErrForbidden:
			writeError(w, "you can only delete your own posts", http.StatusForbidden)
		default:
			writeError(w, "failed to delete post", http.StatusInternalServerError)
		}
		return
	}

	// Успешное удаление — 204 No Content, тело не возвращаем
	w.WriteHeader(http.StatusNoContent)
}

// GetByAuthor возвращает посты конкретного автора
func (h *PostHandler) GetByAuthor(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Путь вида: /api/posts/author/123
	idStr := extractIDFromPath(r.URL.Path, "/api/posts/author/")
	if idStr == "" {
		writeError(w, "invalid author ID path", http.StatusBadRequest)
		return
	}

	authorID, err := strconv.Atoi(idStr)
	if err != nil || authorID <= 0 {
		writeError(w, "invalid author ID", http.StatusBadRequest)
		return
	}

	query := r.URL.Query()
	limit, _ := strconv.Atoi(query.Get("limit"))
	if limit <= 0 {
		limit = 10
	}
	offset, _ := strconv.Atoi(query.Get("offset"))
	if offset < 0 {
		offset = 0
	}

	posts, total, err := h.postService.GetByAuthor(r.Context(), authorID, limit, offset)
	if err != nil {
		// В зависимости от реализации сервиса может быть ErrUserNotFound и т.п.
		writeError(w, "failed to get author posts", http.StatusInternalServerError)
		return
	}

	type PostsResponse struct {
		Posts    []*model.Post `json:"posts"`
		Total    int           `json:"total"`
		Limit    int           `json:"limit"`
		Offset   int           `json:"offset"`
		AuthorID int           `json:"author_id"`
	}

	resp := PostsResponse{
		Posts:    posts,
		Total:    total,
		Limit:    limit,
		Offset:   offset,
		AuthorID: authorID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}
