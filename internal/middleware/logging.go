package middleware

import (
	"log"
	"net/http"
	"time"
)

// LoggingMiddleware provides request logging, CORS, recovery and other utility middleware
type LoggingMiddleware struct {
	logger *log.Logger
}

// NewLoggingMiddleware creates a new logging middleware instance
func NewLoggingMiddleware(logger *log.Logger) *LoggingMiddleware {
	return &LoggingMiddleware{
		logger: logger,
	}
}

// Logger logs all HTTP requests
func (m *LoggingMiddleware) Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Реализовать логирование запросов
		// Шаги:
		// 1. Засечь время начала запроса
		// 2. Создать wrapper для ResponseWriter чтобы захватить статус код
		// 3. Вызвать следующий handler с wrapped writer
		// 4. После выполнения залогировать: метод, путь, IP, статус, время выполнения

		// Временная реализация
		next(w, r)
	}
}

// Recovery восстанавливается после паник
func (m *LoggingMiddleware) Recovery(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Реализовать восстановление после паник
		// Шаги:
		// 1. Использовать defer с recover() для перехвата паник
		// 2. При панике залогировать ошибку
		// 3. Опционально: добавить stack trace
		// 4. Вернуть клиенту 500 Internal Server Error
		// 5. Вызвать следующий handler

		// Временная реализация
		next(w, r)
	}
}

// CORS добавляет CORS заголовки
func (m *LoggingMiddleware) CORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Реализовать CORS заголовки
		// Шаги:
		// 1. Добавить необходимые CORS заголовки (Origin, Methods, Headers, Max-Age)
		// 2. Обработать preflight запросы (OPTIONS метод) - вернуть 204
		// 3. Для остальных методов вызвать следующий handler

		// Временная реализация
		next(w, r)
	}
}

// RequestID добавляет уникальный ID к каждому запросу
func (m *LoggingMiddleware) RequestID(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Реализовать генерацию Request ID
		// Шаги:
		// 1. Сгенерировать уникальный ID (UUID или timestamp+random)
		// 2. Добавить ID в контекст запроса для использования в логах
		// 3. Добавить ID в заголовок ответа X-Request-ID
		// 4. Залогировать запрос с Request ID
		// 5. Вызвать следующий handler

		// Временная реализация
		next(w, r)
	}
}

// RateLimiter ограничивает количество запросов от одного клиента
func (m *LoggingMiddleware) RateLimiter(maxRequests int, window time.Duration) func(http.HandlerFunc) http.HandlerFunc {
	// TODO: Реализовать rate limiting (продвинутое задание)
	// Шаги:
	// 1. Создать хранилище для отслеживания запросов по IP адресам
	// 2. Использовать mutex для безопасного доступа к хранилищу
	// 3. Для каждого запроса:
	//    - Получить IP клиента
	//    - Проверить количество запросов в текущем окне времени
	//    - Если превышен лимит - вернуть 429 Too Many Requests
	//    - Иначе увеличить счетчик и пропустить запрос

	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {

			// Временная реализация
			next(w, r)
		}
	}
}

// ContentTypeJSON устанавливает Content-Type: application/json для всех ответов
func (m *LoggingMiddleware) ContentTypeJSON(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Установить Content-Type: application/json для всех ответов

		// Временная реализация
		next(w, r)
	}
}

// getClientIP извлекает IP адрес клиента
func getClientIP(r *http.Request) string {
	// TODO: Извлечь реальный IP адрес клиента
	// Проверить заголовки: X-Forwarded-For, X-Real-IP, затем RemoteAddr
	// Учесть что X-Forwarded-For может содержать несколько IP

	return r.RemoteAddr
}

// responseWriter обертка для захвата статус кода
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	written    bool
}

// WriteHeader сохраняет статус код
func (rw *responseWriter) WriteHeader(code int) {
	if !rw.written {
		rw.statusCode = code
		rw.ResponseWriter.WriteHeader(code)
		rw.written = true
	}
}

// Write вызывает WriteHeader если еще не был вызван
func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.written {
		rw.WriteHeader(http.StatusOK)
	}
	return rw.ResponseWriter.Write(b)
}

// newResponseWriter создает новую обертку
func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
		written:        false,
	}
}
