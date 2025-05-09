// internal/model/user.go
package model // <- Ubah package

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
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// Input untuk login
type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
