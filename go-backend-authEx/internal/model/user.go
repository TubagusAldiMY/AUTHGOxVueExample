// internal/model/user.go
package model // <- Pastikan package tetap model

import "time"

// User struct merepresentasikan data pengguna
type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Jangan kirim hash password ke client
	CreatedAt    time.Time `json:"created_at"`
}

// Input untuk registrasi
type RegisterInput struct {
	Username string `json:"username" validate:"required,alphanum,min=3,max=30"` // Contoh: alfanumerik, 3-30 karakter
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=72"` // Contoh: minimal 8 karakter, maksimal 72 (batas bcrypt)
	// Anda bisa menambahkan validasi password yang lebih kompleks nanti jika perlu
	// seperti `containsany=!@#$%^&*()`, atau membuat custom validator.
}

// Input untuk login
type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"` // Untuk login, biasanya hanya 'required' sudah cukup
}
