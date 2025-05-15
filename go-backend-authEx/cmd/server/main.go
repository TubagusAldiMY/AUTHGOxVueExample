// cmd/server/main.go
package main

import (
	"context"
	//"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-auth-example/internal/api"
	"go-auth-example/internal/repository"
	"go-auth-example/internal/service"
	"go-auth-example/internal/storage"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading environment variables")
	}

	if os.Getenv("JWT_SECRET_KEY") == "" {
		log.Fatal("FATAL: JWT_SECRET_KEY environment variable is not set.")
	}
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("FATAL: DATABASE_URL environment variable is not set.")
	}

	// --- Inisialisasi Koneksi Database ---
	// Panggil ConnectDB dan simpan instance DB
	db, err := storage.ConnectDB()
	if err != nil {
		log.Fatalf("FATAL: Could not connect to database: %v", err)
	}
	// Defer CloseDB dengan instance db yang valid, dipanggil sebelum main selesai
	// Namun, kita akan memindahkannya ke graceful shutdown untuk penutupan yang lebih terkontrol
	// defer storage.CloseDB(db) // Pindahkan ini ke bawah

	// --- (Opsional) Migrasi/Setup Tabel ---
	// Panggil CreateTableIfNotExists dengan instance db
	if err := storage.CreateTableIfNotExists(db); err != nil {
		log.Fatalf("FATAL: Could not create/check tables: %v", err)
	}

	// --- Dependency Injection ---
	// 1. Buat Repository (inject db)
	// Gunakan instance db yang telah kita dapatkan
	userRepo := repository.NewPostgresUserRepository(db)

	// 2. Buat Services (inject repository)
	authService := service.NewAuthService(userRepo)
	userService := service.NewUserService(userRepo)

	// 3. Buat Handlers (inject services)
	authHandler := api.NewAuthHandler(authService, userService)

	// 4. Setup Router (inject handlers)
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
		log.Printf("Server starting on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on %s: %v\n", port, err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Tutup koneksi database sebelum server mati
	// Panggil CloseDB dengan instance db
	if err := storage.CloseDB(db); err != nil {
		log.Printf("Error closing database: %v", err) // Log error tapi jangan fatal agar shutdown server tetap berjalan
	}

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
