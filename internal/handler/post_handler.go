package handler

import (
	"blog-api/internal/service"
	"net/http"
)

type PostHandler struct {
	postService *service.PostService
}

func NewPostHandler(postService *service.PostService) *PostHandler {
	return &PostHandler{
		postService: postService,
	}
}

// Create обрабатывает создание нового поста
// POST /api/posts
// Требует аутентификации
func (h *PostHandler) Create(w http.ResponseWriter, r *http.Request) {
	// TODO: Реализовать создание поста
	// Шаги:
	// 1. Проверить метод запроса (должен быть POST)
	// 2. Получить userID из контекста (установлен middleware)
	// 3. Декодировать JSON тело в PostCreateRequest
	// 4. Создать пост через postService.Create
	// 5. Вернуть созданный пост как JSON (201 Created)

	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

// GetByID возвращает пост по ID
// GET /api/posts/{id}
// Не требует аутентификации
func (h *PostHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// TODO: Реализовать получение поста по ID
	// Шаги:
	// 1. Проверить метод запроса (должен быть GET)
	// 2. Извлечь ID из URL пути
	// 3. Получить пост через postService.GetByID
	// 4. Обработать ошибки (ErrPostNotFound -> 404)
	// 5. Вернуть пост как JSON (200 OK)

	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

// GetAll возвращает список постов с пагинацией
// GET /api/posts?limit=10&offset=0
// Не требует аутентификации
func (h *PostHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	// TODO: Реализовать получение списка постов
	// Шаги:
	// 1. Проверить метод запроса (должен быть GET)
	// 2. Извлечь параметры пагинации из query string
	// 3. Получить посты через postService.GetAll
	// 4. Создать ответ с метаданными пагинации
	// 5. Вернуть список постов как JSON (200 OK)

	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

// Update обновляет пост
// PUT /api/posts/{id}
// Требует аутентификации, может обновить только автор
func (h *PostHandler) Update(w http.ResponseWriter, r *http.Request) {
	// TODO: Реализовать обновление поста
	// Шаги:
	// 1. Проверить метод запроса (должен быть PUT)
	// 2. Получить userID из контекста
	// 3. Извлечь ID поста из URL
	// 4. Декодировать JSON тело в PostUpdateRequest
	// 5. Обновить через postService.Update
	// 6. Обработать ошибки (404 для не найден, 403 для чужого поста)
	// 7. Вернуть обновленный пост как JSON (200 OK)

	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

// Delete удаляет пост
// DELETE /api/posts/{id}
// Требует аутентификации, может удалить только автор
func (h *PostHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// TODO: Реализовать удаление поста
	// Шаги:
	// 1. Проверить метод запроса (должен быть DELETE)
	// 2. Получить userID из контекста
	// 3. Извлечь ID поста из URL
	// 4. Удалить через postService.Delete
	// 5. Обработать ошибки (404 для не найден, 403 для чужого поста)
	// 6. Вернуть 204 No Content при успехе

	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

// GetByAuthor возвращает посты конкретного автора
// GET /api/posts/author/{authorID}?limit=10&offset=0
// Не требует аутентификации
func (h *PostHandler) GetByAuthor(w http.ResponseWriter, r *http.Request) {
	// TODO: Реализовать получение постов автора
	// Шаги:
	// 1. Проверить метод запроса (должен быть GET)
	// 2. Извлечь ID автора из URL
	// 3. Извлечь параметры пагинации из query string
	// 4. Получить посты через postService.GetByAuthor
	// 5. Создать ответ с метаданными и списком постов
	// 6. Вернуть как JSON (200 OK)

	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

// extractIDFromPath извлекает ID из пути URL
func extractIDFromPath(path, prefix string) string {
	// TODO: Реализовать извлечение ID из пути
	// Пример: path = "/api/posts/123", prefix = "/api/posts/"
	// Должен вернуть "123"

	return ""
}
