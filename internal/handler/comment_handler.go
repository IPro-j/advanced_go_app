package handler

import (
	"blog-api/internal/model"
	"blog-api/internal/service"
	"blog-api/pkg/apperr"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type CommentHandler struct {
	commentService *service.CommentService
}

func NewCommentHandler(commentService *service.CommentService) *CommentHandler {
	return &CommentHandler{
		commentService: commentService,
	}
}

// Create обрабатывает создание нового комментария
func (h *CommentHandler) Create(w http.ResponseWriter, r *http.Request) {
	postIDStr := chi.URLParam(r, "id")
	if postIDStr == "" {
		writeError(w, "missing postID in URL", http.StatusBadRequest)
		return
	}
	postID, err := strconv.Atoi(postIDStr)
	if err != nil || postID <= 0 {
		writeError(w, "invalid postID", http.StatusBadRequest)
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
		case apperr.ErrPostNotExists:
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
func (h *CommentHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		writeError(w, "missing comment id", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		writeError(w, "invalid comment id", http.StatusBadRequest)
		return
	}

	comment, err := h.commentService.GetByID(r.Context(), id)
	if err != nil {
		if err == apperr.ErrCommentNotFound {
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
func (h *CommentHandler) GetByPost(w http.ResponseWriter, r *http.Request) {
	postIDStr := chi.URLParam(r, "id")
	if postIDStr == "" {
		writeError(w, "missing postID in URL", http.StatusBadRequest)
		return
	}
	postID, err := strconv.Atoi(postIDStr)
	if err != nil || postID <= 0 {
		writeError(w, "invalid postID", http.StatusBadRequest)
		return
	}

	query := r.URL.Query()
	limit, _ := strconv.Atoi(query.Get("limit"))
	if limit <= 0 {
		limit = 20
	}
	offset, _ := strconv.Atoi(query.Get("offset"))
	if offset < 0 {
		offset = 0
	}

	comments, total, err := h.commentService.GetByPost(r.Context(), postID, limit, offset)
	if err != nil {
		if err == apperr.ErrPostNotExists {
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

// Update — PUT /api/protected/comments/{id}
func (h *CommentHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r.Context())
	if !ok {
		writeError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		writeError(w, "missing comment id", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		writeError(w, "invalid comment id", http.StatusBadRequest)
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
		case apperr.ErrCommentNotFound:
			writeError(w, "comment not found", http.StatusNotFound)
		case apperr.ErrForbidden:
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

// Delete — DELETE /api/protected/comments/{id}
func (h *CommentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r.Context())
	if !ok {
		writeError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		writeError(w, "missing comment id", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		writeError(w, "invalid comment id", http.StatusBadRequest)
		return
	}

	err = h.commentService.Delete(r.Context(), id, userID)
	if err != nil {
		switch err {
		case apperr.ErrCommentNotFound:
			writeError(w, "comment not found", http.StatusNotFound)
		case apperr.ErrForbidden:
			writeError(w, "you can only delete your own comments", http.StatusForbidden)
		default:
			writeError(w, "failed to delete comment", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
