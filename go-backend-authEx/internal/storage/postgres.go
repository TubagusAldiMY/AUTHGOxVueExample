// internal/storage/postgres.go
package storage // <- Ubah package

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib" // Import driver pgx
)

var DB *sql.DB // Variabel global untuk koneksi database (sementara)

// ConnectDB menginisialisasi koneksi ke database PostgreSQL
func ConnectDB() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	var err error
	DB, err = sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	// Cek koneksi
	err = DB.Ping()
	if err != nil {
		DB.Close() // Tutup koneksi jika ping gagal
		log.Fatalf("Unable to ping database: %v\n", err)
	}

	fmt.Println("Successfully connected to database!")

	// (Opsional) Anda bisa membuat tabel di sini jika belum ada
	// createTableIfNotExists()
}

// (Opsional) Fungsi untuk membuat tabel jika belum ada
func CreateTableIfNotExists() {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(50) UNIQUE NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL,
		password_hash VARCHAR(255) NOT NULL,
		created_at TIMESTAMPTZ DEFAULT NOW()
	);`

	_, err := DB.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Unable to create users table: %v\n", err)
	}
	fmt.Println("Users table checked/created successfully.")
}

// CloseDB menutup koneksi database
func CloseDB() {
	if DB != nil {
		DB.Close()
		fmt.Println("Database connection closed.")
	}
}
