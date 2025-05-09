package auth

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword mengenkripsi password menggunakan bcrypt
func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedBytes), nil
}

// CheckPasswordHash membandingkan password plain text dengan hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil // Jika tidak ada error, password cocok
}
