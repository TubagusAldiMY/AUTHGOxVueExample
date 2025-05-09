// cmd/server/main.go
package main

import (
	"context" // <- Import context
	"log"
	"net/http" // <- Import net/http
	"os"
	"os/signal" // <- Import os/signal
	"syscall"   // <- Import syscall
	"time"      // <- Import time

	// Import internal packages dengan path lengkap
	"go-auth-example/internal/api"
	"go-auth-example/internal/repository"
	"go-auth-example/internal/service"
	"go-auth-example/internal/storage"

	"github.com/joho/godotenv"
)

func main() {
	// Muat .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading environment variables")
	}

	// --- Validasi Konfigurasi Awal ---
	if os.Getenv("JWT_SECRET_KEY") == "" {
		log.Fatal("FATAL: JWT_SECRET_KEY environment variable is not set.")
	}
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("FATAL: DATABASE_URL environment variable is not set.")
	}

	// --- Inisialisasi Koneksi Database ---
	storage.ConnectDB() // Gunakan fungsi dari package storage
	// Pindahkan defer CloseDB ke bagian graceful shutdown

	// --- (Opsional) Migrasi/Setup Tabel ---
	storage.CreateTableIfNotExists() // Gunakan fungsi dari package storage

	// --- Dependency Injection ---
	// 1. Buat Repository (inject DB)
	// Gunakan variabel DB dari storage (untuk sementara)
	userRepo := repository.NewPostgresUserRepository(storage.DB)

	// 2. Buat Services (inject repository)
	authService := service.NewAuthService(userRepo)
	userService := service.NewUserService(userRepo)

	// 3. Buat Handlers (inject services)
	authHandler := api.NewAuthHandler(authService, userService)

	// 4. Setup Router (inject handlers)
	router := api.SetupRouter(authHandler)

	// --- Konfigurasi Server HTTP ---
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  10 * time.Second, // Tambahkan timeout
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// --- Jalankan Server dalam Goroutine ---
	go func() {
		log.Printf("Server starting on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on %s: %v\n", port, err)
		}
	}()

	// --- Graceful Shutdown ---
	quit := make(chan os.Signal, 1)
	// Menunggu sinyal interrupt (Ctrl+C) atau sinyal terminate
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // Blok sampai sinyal diterima
	log.Println("Shutting down server...")

	// Beri waktu (misal 5 detik) untuk request yang sedang berjalan selesai
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Tutup koneksi database sebelum server mati
	storage.CloseDB()

	// Matikan server HTTP
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
