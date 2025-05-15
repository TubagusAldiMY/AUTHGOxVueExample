// cmd/server/main.go
package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	// Import internal packages
	"go-auth-example/internal/api"
	"go-auth-example/internal/logger" // <-- IMPORT LOGGER
	"go-auth-example/internal/repository"
	"go-auth-example/internal/service"
	"go-auth-example/internal/storage"

	"github.com/joho/godotenv"
)

func main() {
	// Muat .env SEBELUM logger diinisialisasi jika logger bergantung pada env vars
	if err := godotenv.Load(); err != nil {
		// Gunakan log standar Go di sini karena logger kita mungkin belum siap
		// atau logger.Log akan menggunakan default jika env var tidak ada
		// log.Println("No .env file found, reading environment variables")
		// Untuk sekarang, biarkan logger.Log yang menangani pesan ini
	}

	// Logger akan diinisialisasi secara otomatis saat package logger diimpor.
	// Anda bisa menambahkan pesan log pertama di sini jika mau.
	// logger.Log.Info("Application starting...") // Pesan ini sudah ada di init logger

	if os.Getenv("JWT_SECRET_KEY") == "" {
		logger.Log.Fatal("FATAL: JWT_SECRET_KEY environment variable is not set.")
	}
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		logger.Log.Fatal("FATAL: DATABASE_URL environment variable is not set.")
	}

	db, err := storage.ConnectDB()
	if err != nil {
		logger.Log.Fatalf("FATAL: Could not connect to database: %v", err)
	}

	if err := storage.CreateTableIfNotExists(db); err != nil {
		logger.Log.Fatalf("FATAL: Could not create/check tables: %v", err)
	}

	userRepo := repository.NewPostgresUserRepository(db)
	authService := service.NewAuthService(userRepo)
	userService := service.NewUserService(userRepo)
	authHandler := api.NewAuthHandler(authService, userService)
	router := api.SetupRouter(authHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		logger.Log.Infof("Server starting on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Fatalf("Could not listen on %s: %v\n", port, err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := storage.CloseDB(db); err != nil {
		logger.Log.Errorf("Error closing database: %v", err)
	}

	if err := server.Shutdown(ctx); err != nil {
		logger.Log.Fatal("Server forced to shutdown:", err)
	}

	logger.Log.Info("Server exiting")
}
