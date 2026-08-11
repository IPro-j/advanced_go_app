package handler

import (
	"blog-api/internal/model"
	"blog-api/internal/service"
	"blog-api/pkg/apperr"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type PostHandler struct {
	postService *service.PostService
}

const (
	defaultLimit = 20
	maxLimit     = 100
)

func NewPostHandler(postService *service.PostService) *PostHandler {
	return &PostHandler{
		postService: postService,
	}
}

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
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		writeError(w, "missing post id", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		writeError(w, "invalid post id", http.StatusBadRequest)
		return
	}

	post, err := h.postService.GetByID(r.Context(), id)
	if err != nil {
		if err == apperr.ErrPostNotFound {
			writeError(w, "post not found", http.StatusNotFound)
		} else {
			log.Printf("GetByID error: %v", err)
			writeError(w, "failed to get post", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(post)
}

// GetAll возвращает список постов с пагинацией
/*func (h *PostHandler) GetAll(w http.ResponseWriter, r *http.Request) {
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
}*/

// GetAll возвращает список постов с пагинацией
func (h *PostHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query()
	limitStr := query.Get("limit")
	offsetStr := query.Get("offset")

	// 1. Валидация limit
	limit := defaultLimit
	if limitStr != "" {
		l, err := strconv.Atoi(limitStr)
		if err != nil || l <= 0 {
			writeError(w, "invalid limit: must be a positive integer", http.StatusBadRequest)
			return
		}
		limit = l
	}

	// 2. Нормализация limit (после валидации!)
	if limit > maxLimit {
		limit = maxLimit
	}

	// 3. Валидация offset
	offset := 0
	if offsetStr != "" {
		o, err := strconv.Atoi(offsetStr)
		if err != nil || o < 0 {
			writeError(w, "invalid offset: must be a non-negative integer", http.StatusBadRequest)
			return
		}
		offset = o
	}

	posts, total, err := h.postService.GetAll(r.Context(), limit, offset)
	if err != nil {
		log.Printf("GetAll posts error: %v", err)
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
		Limit:  limit, // теперь это нормализованное значение
		Offset: offset,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

// Update обновляет пост
func (h *PostHandler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		writeError(w, "method not allowed. Use PATCH for partial update.", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := getUserIDFromContext(r.Context())
	if !ok {
		writeError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		writeError(w, "missing post id", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		writeError(w, "invalid post id", http.StatusBadRequest)
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
		case apperr.ErrPostNotFound:
			writeError(w, "post not found", http.StatusNotFound)
		case apperr.ErrForbidden:
			writeError(w, "you can only update your own posts", http.StatusForbidden)
		default:
			log.Printf("Update post error: %v", err)
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
	userID, ok := getUserIDFromContext(r.Context())
	if !ok {
		writeError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		writeError(w, "missing post id", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		writeError(w, "invalid post id", http.StatusBadRequest)
		return
	}

	err = h.postService.Delete(r.Context(), id, userID)
	if err != nil {
		switch err {
		case apperr.ErrPostNotFound:
			writeError(w, "post not found", http.StatusNotFound)
		case apperr.ErrForbidden:
			writeError(w, "you can only delete your own posts", http.StatusForbidden)
		default:
			log.Printf("Delete post error: %v", err)
			writeError(w, "failed to delete post", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetByAuthor возвращает посты конкретного автора
func (h *PostHandler) GetByAuthor(w http.ResponseWriter, r *http.Request) {
	authorIDStr := chi.URLParam(r, "author_id") // нужно, чтобы в роуте было {author_id}
	if authorIDStr == "" {
		writeError(w, "missing author id", http.StatusBadRequest)
		return
	}

	authorID, err := strconv.Atoi(authorIDStr)
	if err != nil || authorID <= 0 {
		writeError(w, "invalid author id", http.StatusBadRequest)
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

	/*posts, total, err := h.postService.GetByAuthor(r.Context(), authorID, limit, offset)
	if err != nil {
		writeError(w, "failed to get author posts", http.StatusInternalServerError)
		return
	}*/

	type PostsResponse struct {
		Limit    int `json:"limit"`
		Offset   int `json:"offset"`
		AuthorID int `json:"author_id"`
	}

	resp := PostsResponse{

		Limit:    limit,
		Offset:   offset,
		AuthorID: authorID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}
