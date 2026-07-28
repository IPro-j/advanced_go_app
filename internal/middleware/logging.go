package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

// contextKey — кастомный тип для ключей контекста, чтобы избежать коллизий
type contextKey string

const (
	RequestIDKey contextKey = "requestID"
)

// LoggingMiddleware предоставляет утилиты: логирование, CORS, recovery, request ID и т.д.
type LoggingMiddleware struct {
	logger *log.Logger
}

// NewLoggingMiddleware создаёт новый экземпляр middleware
func NewLoggingMiddleware(logger *log.Logger) *LoggingMiddleware {
	return &LoggingMiddleware{
		logger: logger,
	}
}

// Logger логирует все HTTP-запросы с временем выполнения и статусом
func (m *LoggingMiddleware) Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Получаем RequestID из контекста (если он уже был добавлен, например, в RequestID middleware)
		reqID, _ := GetRequestIDFromContext(r.Context())
		if reqID == "" {
			reqID = uuid.New().String()
		}

		rw := newResponseWriter(w)
		next(rw, r.WithContext(context.WithValue(r.Context(), RequestIDKey, reqID)))

		duration := time.Since(start)
		ip := getClientIP(r)

		m.logger.Printf("[%s] %s %s | IP: %s | Status: %d | Duration: %v",
			reqID,
			r.Method,
			r.URL.Path,
			ip,
			rw.statusCode,
			duration,
		)
	}
}

func (m *LoggingMiddleware) Chain() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		// Порядок важен: CORS → RequestID → Logger → Recovery
		h := http.HandlerFunc(m.CORS(next.ServeHTTP))
		h = http.HandlerFunc(m.RequestID(h.ServeHTTP))
		h = http.HandlerFunc(m.Logger(h.ServeHTTP))
		h = http.HandlerFunc(m.Recovery(h.ServeHTTP))

		return h
	}
}

// Recovery перехватывает паники и возвращает 500, логируя стек
func (m *LoggingMiddleware) Recovery(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if v := recover(); v != nil {
				reqID, _ := GetRequestIDFromContext(r.Context())

				stack := string(debug.Stack())
				msg := fmt.Sprintf("panic recovered: %v\nstack:\n%s", r, stack)

				m.logger.Printf("[%s] PANIC: %s", reqID, msg)

				// Безопасный ответ клиенту: не отдаём стек, только общую ошибку
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}()

		next(w, r)
	}
}

// CORS добавляет CORS-заголовки и обрабатывает preflight (OPTIONS)
func (m *LoggingMiddleware) CORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		// В продакшене лучше проверять origin по списку разрешённых доменов
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Max-Age", "3600")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next(w, r)
	}
}

// RequestID генерирует уникальный ID запроса, кладёт в контекст и добавляет в заголовок ответа
func (m *LoggingMiddleware) RequestID(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqID := uuid.New().String()

		ctx := context.WithValue(r.Context(), RequestIDKey, reqID)
		r = r.WithContext(ctx)

		w.Header().Set("X-Request-ID", reqID)

		next(w, r)
	}
}

// RateLimiter — простая реализация rate limiting по IP (в памяти, без persistence)
func (m *LoggingMiddleware) RateLimiter(maxRequests int, window time.Duration) func(http.HandlerFunc) http.HandlerFunc {
	type limiterState struct {
		count     int
		lastReset time.Time
	}

	// Хранилище: IP -> состояние
	store := make(map[string]*limiterState)
	mu := &sync.Mutex{} // нужно добавить import "sync"

	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			mu.Lock()
			defer mu.Unlock()

			now := time.Now()
			ip := getClientIP(r)
			state, ok := store[ip]

			if !ok {
				state = &limiterState{
					count:     1,
					lastReset: now,
				}
				store[ip] = state
				next(w, r)
				return
			}

			// Если окно времени истекло — сбрасываем счётчик
			if now.Sub(state.lastReset) >= window {
				state.count = 1
				state.lastReset = now
				next(w, r)
				return
			}

			// Окно ещё не истекло: проверяем лимит
			if state.count >= maxRequests {
				m.logger.Printf("Rate limit exceeded for IP: %s", ip)
				http.Error(w, "Too many requests", http.StatusTooManyRequests)
				return
			}

			state.count++
			next(w, r)
		}
	}
}

// ContentTypeJSON принудительно ставит Content-Type: application/json для всех ответов
func (m *LoggingMiddleware) ContentTypeJSON(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	}
}

// getClientIP извлекает реальный IP клиента с учётом прокси
func getClientIP(r *http.Request) string {
	// X-Forwarded-For: client, proxy1, proxy2
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		parts := strings.Split(forwarded, ",")
		if len(parts) > 0 {
			return strings.TrimSpace(parts[0])
		}
	}

	xRealIP := r.Header.Get("X-Real-IP")
	if xRealIP != "" {
		return xRealIP
	}

	// Fallback: RemoteAddr (может быть в формате ip:port)
	addr := r.RemoteAddr
	if idx := strings.LastIndex(addr, ":"); idx != -1 {
		addr = addr[:idx]
	}
	return addr
}

// responseWriter — обёртка для перехвата статус-кода
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	written    bool
}

func (rw *responseWriter) WriteHeader(code int) {
	if !rw.written {
		rw.statusCode = code
		rw.ResponseWriter.WriteHeader(code)
		rw.written = true
	}
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.written {
		rw.WriteHeader(http.StatusOK)
	}
	return rw.ResponseWriter.Write(b)
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
		written:        false,
	}
}

// GetRequestIDFromContext извлекает RequestID из контекста
func GetRequestIDFromContext(ctx context.Context) (string, bool) {
	val, ok := ctx.Value(RequestIDKey).(string)
	return val, ok
}
