package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"blog-api/internal/handler"
	"blog-api/internal/middleware"
	"blog-api/internal/repository"
	"blog-api/internal/service"
	"blog-api/pkg/auth"
	"blog-api/pkg/database"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"

	"github.com/joho/godotenv"
)

func main() {
	logger := log.New(os.Stdout, "[API] ", log.LstdFlags|log.Lshortfile)

	// Создаём экземпляр  middleware
	loggingMW := middleware.NewLoggingMiddleware(logger)

	if err := godotenv.Load(); err != nil {
		logger.Printf("Warning: .env file not found, using environment variables directly")
	}

	cfg := loadConfig()

	dbConfig := database.Config{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		DBName:   cfg.DBName,
		SSLMode:  cfg.DBSSLMode,
	}
	db, err := database.NewPostgresDB(dbConfig)
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	logger.Println("Database connected successfully")

	if err := database.Migrate(db); err != nil {
		logger.Fatalf("Migration failed: %v", err)
	}
	logger.Println("Migrations applied successfully")

	jwtManager, err := auth.NewJWTManager(cfg.JWTSecret, cfg.JWTExpiryHours)
	if err != nil {
		logger.Fatalf("Failed to create JWT manager: %v", err)
	}

	postRepo := repository.NewPostRepo(db)
	commentRepo := repository.NewCommentRepo(db)
	userRepo := repository.NewUserRepo(db)

	userService := service.NewUserService(userRepo, jwtManager)
	postService := service.NewPostService(postRepo, userRepo)
	commentService := service.NewCommentService(commentRepo, postRepo)

	authHandler := handler.NewAuthHandler(userService, jwtManager)
	postHandler := handler.NewPostHandler(postService)
	commentHandler := handler.NewCommentHandler(commentService)

	authMW := middleware.NewAuthMiddleware(jwtManager)

	router := chi.NewRouter()

	// Добавляем цепочку: CORS → RequestID → L ogger → Recovery
	router.Use(loggingMW.Chain())

	// Health check

	// Единый блок /api
	router.Route("/api", func(r chi.Router) {

		//проверка состояния сервиса (GET /api/health)
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status":"ok","service":"blog-api"}`))
		})

		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
		r.Get("/posts", postHandler.GetAll)
		r.Get("/posts/{id}", postHandler.GetByID)
		r.Get("/posts/{id}/comments", commentHandler.GetByPost)

		r.Group(func(r chi.Router) {
			r.Use(authMW.AuthMiddlewareForChi()) // авторизация только тут
			r.Post("/posts", postHandler.Create)
			r.Patch("/posts/{id}", postHandler.Update)
			r.Delete("/posts/{id}", postHandler.Delete)
			r.Post("/posts/{id}/comments", commentHandler.Create)
			r.Put("/comments/{id}", commentHandler.Update)
			r.Delete("/comments/{id}", commentHandler.Delete)
		})

	})

	addr := cfg.ServerHost + ":" + strconv.Itoa(cfg.ServerPort)

	readTimeout := parseDurationWithDefault(os.Getenv("HTTP_READ_TIMEOUT"), 15*time.Second)
	writeTimeout := parseDurationWithDefault(os.Getenv("HTTP_WRITE_TIMEOUT"), 15*time.Second)
	idleTimeout := parseDurationWithDefault(os.Getenv("HTTP_IDLE_TIMEOUT"), 60*time.Second)
	shutdownGracePeriod := parseDurationWithDefault(os.Getenv("HTTP_SHUTDOWN_GRACE_PERIOD"), 30*time.Second)

	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}

	logger.Printf("Starting server on %s (read=%v, write=%v, idle=%v, shutdown_grace=%v)",
		addr, readTimeout, writeTimeout, idleTimeout, shutdownGracePeriod)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Server failed to start: %v", err)
		}
	}()

	<-stop
	logger.Println("Shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownGracePeriod)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Printf("Graceful shutdown failed: %v", err)
	} else {
		logger.Println("Server stopped gracefully")
	}
}

type Config struct {
	ServerHost      string
	ServerPort      int
	DBHost          string
	DBPort          int
	DBUser          string
	DBPassword      string
	DBName          string
	DBSSLMode       string
	JWTSecret       string
	JWTExpiryHours  int
	CacheTTLMinutes int
}

func loadConfig() *Config {
	return &Config{
		ServerHost:      getEnv("SERVER_HOST", "0.0.0.0"),
		ServerPort:      getEnvAsInt("SERVER_PORT", 8080),
		DBHost:          getEnv("DB_HOST", "localhost"),
		DBPort:          getEnvAsInt("DB_PORT", 5432),
		DBUser:          getEnv("DB_USER", "blouser"),
		DBPassword:      getEnv("DB_PASSWORD", "blogpassword"),
		DBName:          getEnv("DB_NAME", "blogdb"),
		DBSSLMode:       getEnv("DB_SSLMODE", "disable"),
		JWTSecret:       getEnv("JWT_SECRET", ""),
		JWTExpiryHours:  getEnvAsInt("JWT_EXPIRY_HOURS", 24),
		CacheTTLMinutes: getEnvAsInt("CACHE_TTL_MINUTES", 60),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

// parseDurationWithDefault парсит строку вида "15s", "1m", "2h" и т.п.
// Если строка невалидна — возвращает значение по умолчанию.
func parseDurationWithDefault(value string, defaultVal time.Duration) time.Duration {
	if value == "" {
		return defaultVal
	}
	d, err := time.ParseDuration(value)
	if err != nil {
		log.Printf("Invalid duration for %s, using default %v", value, defaultVal)
		return defaultVal
	}
	return d
}
