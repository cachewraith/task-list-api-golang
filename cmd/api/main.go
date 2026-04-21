package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/example/todo-api/internal/config"
	"github.com/example/todo-api/internal/handler"
	appmiddleware "github.com/example/todo-api/internal/middleware"
	"github.com/example/todo-api/internal/repository/postgres"
	"github.com/example/todo-api/internal/service"
	"github.com/example/todo-api/pkg/logger"
	"github.com/example/todo-api/pkg/notifier"
	"github.com/example/todo-api/pkg/response"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file if present
	_ = godotenv.Load()

	// Initialize Telegram notifier for ERROR/WARN alerts
	telegramNotifier := notifier.GetTelegramNotifier()
	telegramNotifier.Configure(
		os.Getenv("TELEGRAM_BOT_TOKEN"),
		os.Getenv("TELEGRAM_GROUP_ID"),
		os.Getenv("TELEGRAM_THREAD_ID"),
	)
	if telegramNotifier.IsEnabled() {
		logger.GetLogger().SetNotifier(telegramNotifier)
		log.Println("✓ Telegram alerts enabled for ERROR/WARN logs")
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create database connection pool
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, cfg.Database.DSN())
	if err != nil {
		log.Fatalf("Failed to create connection pool: %v", err)
	}
	defer pool.Close()

	// Verify database connection
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Connected to database successfully")

	// Initialize repository layer
	todoRepo := postgres.NewTodoRepository(pool)
	userRepo := postgres.NewUserRepository(pool)

	// Get JWT secret from env
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-secret-key-change-this-in-production"
		log.Println("Warning: Using default JWT_SECRET. Set JWT_SECRET env var in production!")
	}
	tokenExpiry := 24 * time.Hour

	// Initialize service layer with dependency injection
	todoService := service.NewTodoService(todoRepo)
	authService := service.NewAuthService(userRepo, jwtSecret, tokenExpiry)

	// Initialize handler layer with dependency injection
	todoHandler := handler.NewTodoHandler(todoService)
	authHandler := handler.NewAuthHandler(authService)
	logHandler := handler.NewLogHandler()

	// Setup router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.Timeout(cfg.Server.Timeout))
	r.Use(middleware.SetHeader("Content-Type", "application/json"))

	// Health check (public)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		response.NewSuccessWithMessage("Service is healthy", map[string]string{"status": "up"}).Write(w, http.StatusOK)
	})

	// Auth routes (public)
	r.Route("/api/v1/auth", func(r chi.Router) {
		r.Mount("/", authHandler.Routes())
	})

	// Protected routes - require authentication
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(appmiddleware.Auth(authService))
		r.Mount("/todos", todoHandler.Routes())
		r.Mount("/logs", logHandler.Routes())
	})

	// 404 handler
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		response.NotFound(w, "The requested resource was not found")
	})

	// Create HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  cfg.Server.Timeout,
		WriteTimeout: cfg.Server.Timeout,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %d", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}
