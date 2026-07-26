package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"blog-api/internal/handler"
	"blog-api/internal/middleware" // твой кастомный middleware
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

	// Создаём экземпляр твоего middleware
	loggingMW := middleware.NewLoggingMiddleware(logger)

	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found, using environment variables directly")
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
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	log.Println("Database connected successfully")

	if err := database.Migrate(db); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	log.Println("Migrations applied successfully")

	jwtManager, err := auth.NewJWTManager(cfg.JWTSecret, cfg.JWTExpiryHours)
	if err != nil {
		log.Fatalf("Failed to create JWT manager: %v", err)
	}

	postRepo := repository.NewPostRepo(db)
	commentRepo := repository.NewCommentRepo(db)
	userRepo := repository.NewUserRepo(db)

	userService := service.NewUserService(userRepo, jwtManager)
	postService := service.NewPostService(postRepo, userRepo)
	commentService := service.NewCommentService(commentRepo, postRepo, userRepo)

	authHandler := handler.NewAuthHandler(userService, jwtManager)
	postHandler := handler.NewPostHandler(postService)
	commentHandler := handler.NewCommentHandler(commentService)

	authMW := middleware.NewAuthMiddleware(jwtManager)

	router := chi.NewRouter()
	//router.Use(loggingMW.Chain())

	// --- ВАЖНО: подключаем твой middleware через экземпляр ---
	// Добавляем цепочку: CORS → RequestID → Logger → Recovery
	router.Use(loggingMW.Chain())

	// Health check
	router.Get("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","service":"blog-api"}`))
	})

	// Единый блок /api
	router.Route("/api", func(r chi.Router) {

		// Публичные маршруты
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
		r.Get("/posts", postHandler.GetAll)
		// другие публичные...

		// Защищённая зона: применяем middleware
		// --- Защищённая группа ---
		// Сначала middleware, потом маршруты
		r.Route("/protected", func(pr chi.Router) {
			pr.Use(authMW.AuthMiddlewareForChi()) // авторизация только тут

			pr.Post("/posts", postHandler.Create)
			pr.Put("/posts/{id}", postHandler.Update)
			pr.Delete("/posts/{id}", postHandler.Delete)

			pr.Post("/posts/{id}/comments", commentHandler.Create)
			pr.Put("/comments/{id}", commentHandler.Update)
			//pr.Delete("/comments/{id}", commentHandler.Delete)
		})
	})

	addr := cfg.ServerHost + ":" + strconv.Itoa(cfg.ServerPort)
	log.Printf("Starting server on %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
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
		DBUser:          getEnv("DB_USER", "postgres"),
		DBPassword:      getEnv("DB_PASSWORD", "password"),
		DBName:          getEnv("DB_NAME", "blog_db"),
		DBSSLMode:       getEnv("DB_SSLMODE", "disable"),
		JWTSecret:       getEnv("JWT_SECRET", "supersecretkey"),
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
