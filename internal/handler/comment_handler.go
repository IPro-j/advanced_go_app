package handler

import (
	"blog-api/internal/service"
	"net/http"
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
// POST /api/comments
// Требует аутентификации
func (h *CommentHandler) Create(w http.ResponseWriter, r *http.Request) {
	// TODO: Реализовать создание комментария
	// HINT: Последовательность действий:

	// 1. Проверить метод запроса (должен быть POST)
	//    if r.Method != http.MethodPost {
	//        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	//        return
	//    }

	// 2. Получить ID пользователя из контекста
	//    userID, ok := getUserIDFromContext(r.Context())
	//    if !ok {
	//        writeError(w, "Unauthorized", http.StatusUnauthorized)
	//        return
	//    }

	// 3. Декодировать тело запроса
	//    var req model.CommentCreateRequest
	//    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
	//        writeError(w, "Invalid request body", http.StatusBadRequest)
	//        return
	//    }

	// 4. Создать комментарий через сервис
	//    comment, err := h.commentService.Create(r.Context(), userID, &req)
	//    if err != nil {
	//        switch err {
	//        case service.ErrPostNotExists:
	//            writeError(w, "Post not found", http.StatusNotFound)
	//        default:
	//            writeError(w, "Failed to create comment", http.StatusInternalServerError)
	//        }
	//        return
	//    }

	// 5. Отправить успешный ответ
	//    w.Header().Set("Content-Type", "application/json")
	//    w.WriteHeader(http.StatusCreated)
	//    json.NewEncoder(w).Encode(comment)

	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

// GetByID возвращает комментарий по ID
// GET /api/comments/{id}
// Не требует аутентификации
func (h *CommentHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// TODO: Реализовать получение комментария по ID
	// HINT: Последовательность действий:

	// 1. Проверить метод запроса (должен быть GET)
	//    if r.Method != http.MethodGet {
	//        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	//        return
	//    }

	// 2. Извлечь ID из URL
	//    // Примерный URL: /api/comments/123
	//    idStr := extractIDFromPath(r.URL.Path, "/api/comments/")
	//    id, err := strconv.Atoi(idStr)
	//    if err != nil {
	//        writeError(w, "Invalid comment ID", http.StatusBadRequest)
	//        return
	//    }

	// 3. Получить комментарий через сервис
	//    comment, err := h.commentService.GetByID(r.Context(), id)
	//    if err != nil {
	//        if err == service.ErrCommentNotFound {
	//            writeError(w, "Comment not found", http.StatusNotFound)
	//        } else {
	//            writeError(w, "Failed to get comment", http.StatusInternalServerError)
	//        }
	//        return
	//    }

	// 4. Отправить ответ
	//    w.Header().Set("Content-Type", "application/json")
	//    w.WriteHeader(http.StatusOK)
	//    json.NewEncoder(w).Encode(comment)

	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

// GetByPost возвращает комментарии к посту
// GET /api/posts/{id}/comments?limit=20&offset=0
// Не требует аутентификации
func (h *CommentHandler) GetByPost(w http.ResponseWriter, r *http.Request) {
	// TODO: Реализовать получение комментариев к посту
	// HINT: Последовательность действий:

	// 1. Проверить метод запроса (должен быть GET)
	//    if r.Method != http.MethodGet {
	//        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	//        return
	//    }

	// 2. Извлечь ID поста из URL
	//    // URL вида: /api/posts/123/comments
	//    // Нужно извлечь "123"
	//    path := r.URL.Path
	//    // Пример парсинга:
	//    // - убрать префикс "/api/posts/"
	//    // - взять часть до "/comments"
	//    idStr := extractPostIDFromCommentsPath(path)
	//    postID, err := strconv.Atoi(idStr)
	//    if err != nil {
	//        writeError(w, "Invalid post ID", http.StatusBadRequest)
	//        return
	//    }

	// 3. Извлечь параметры пагинации
	//    query := r.URL.Query()
	//    limit, _ := strconv.Atoi(query.Get("limit"))
	//    if limit <= 0 {
	//        limit = 20 // значение по умолчанию
	//    }
	//    offset, _ := strconv.Atoi(query.Get("offset"))
	//    if offset < 0 {
	//        offset = 0
	//    }

	// 4. Получить комментарии через сервис
	//    comments, total, err := h.commentService.GetByPost(r.Context(), postID, limit, offset)
	//    if err != nil {
	//        if err == service.ErrPostNotExists {
	//            writeError(w, "Post not found", http.StatusNotFound)
	//        } else {
	//            writeError(w, "Failed to get comments", http.StatusInternalServerError)
	//        }
	//        return
	//    }

	// 5. Создать ответ с метаданными
	//    type CommentsResponse struct {
	//        Comments []*model.Comment `json:"comments"`
	//        Total    int             `json:"total"`
	//        Limit    int             `json:"limit"`
	//        Offset   int             `json:"offset"`
	//        PostID   int             `json:"post_id"`
	//    }
	//
	//    resp := CommentsResponse{
	//        Comments: comments,
	//        Total:    total,
	//        Limit:    limit,
	//        Offset:   offset,
	//        PostID:   postID,
	//    }

	// 6. Отправить ответ
	//    w.Header().Set("Content-Type", "application/json")
	//    w.WriteHeader(http.StatusOK)
	//    json.NewEncoder(w).Encode(resp)

	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

// Update обновляет комментарий
// PUT /api/comments/{id}
// Требует аутентификации, может обновить только автор
func (h *CommentHandler) Update(w http.ResponseWriter, r *http.Request) {
	// TODO: Реализовать обновление комментария
	// HINT: Последовательность действий:

	// 1. Проверить метод запроса (должен быть PUT)
	//    if r.Method != http.MethodPut {
	//        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	//        return
	//    }

	// 2. Получить ID пользователя из контекста
	//    userID, ok := getUserIDFromContext(r.Context())
	//    if !ok {
	//        writeError(w, "Unauthorized", http.StatusUnauthorized)
	//        return
	//    }

	// 3. Извлечь ID комментария из URL
	//    idStr := extractIDFromPath(r.URL.Path, "/api/comments/")
	//    id, err := strconv.Atoi(idStr)
	//    if err != nil {
	//        writeError(w, "Invalid comment ID", http.StatusBadRequest)
	//        return
	//    }

	// 4. Декодировать тело запроса
	//    var req model.CommentUpdateRequest
	//    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
	//        writeError(w, "Invalid request body", http.StatusBadRequest)
	//        return
	//    }

	// 5. Обновить комментарий через сервис
	//    comment, err := h.commentService.Update(r.Context(), id, userID, &req)
	//    if err != nil {
	//        switch err {
	//        case service.ErrCommentNotFound:
	//            writeError(w, "Comment not found", http.StatusNotFound)
	//        case service.ErrForbidden:
	//            writeError(w, "You can only update your own comments", http.StatusForbidden)
	//        default:
	//            writeError(w, "Failed to update comment", http.StatusInternalServerError)
	//        }
	//        return
	//    }

	// 6. Отправить обновленный комментарий
	//    w.Header().Set("Content-Type", "application/json")
	//    w.WriteHeader(http.StatusOK)
	//    json.NewEncoder(w).Encode(comment)

	http.Error(w, "Not implemented", http.Status
