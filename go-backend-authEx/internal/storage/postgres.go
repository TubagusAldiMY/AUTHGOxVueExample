// internal/storage/postgres.go
package storage

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib" // Import driver pgx
)

// ConnectDB menginisialisasi dan mengembalikan koneksi ke database PostgreSQL
func ConnectDB() (*sql.DB, error) { // Sudah benar, mengembalikan *sql.DB dan error
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is required")
	}

	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		return nil, fmt.Errorf("unable to open database connection: %w", err)
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	fmt.Println("Successfully connected to database!")
	return db, nil
}

// CreateTableIfNotExists membuat tabel users jika belum ada
func CreateTableIfNotExists(db *sql.DB) error { // Sudah benar, menerima *sql.DB
	createTableSQL := `
    CREATE TABLE IF NOT EXISTS users (
       id SERIAL PRIMARY KEY,
       username VARCHAR(50) UNIQUE NOT NULL,
       email VARCHAR(255) UNIQUE NOT NULL,
       password_hash VARCHAR(255) NOT NULL,
       created_at TIMESTAMPTZ DEFAULT NOW()
    );`

	_, err := db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("unable to create users table: %w", err)
	}
	fmt.Println("Users table checked/created successfully.")
	return nil
}

// CloseDB menutup koneksi database
func CloseDB(db *sql.DB) error { // Sudah benar, menerima *sql.DB
	if db != nil {
		err := db.Close()
		if err != nil {
			return fmt.Errorf("error closing database connection: %w", err)
		}
		fmt.Println("Database connection closed.")
	}
	return nil
}
