package auth

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"go-auth-example/internal/model"
	"os"
	"time"
)

var jwtSecretKey = []byte(os.Getenv("JWT_SECRET_KEY"))

// GenerateJWT membuat token JWT baru untuk user
func GenerateJWT(user model.User) (string, error) {
	// Set standard claims
	claims := jwt.MapClaims{
		"sub": user.ID,         // Subject (biasanya ID user)
		"iss": "your-app-name", // Issuer (nama aplikasi Anda)
		// "aud": "your-audience", // Audience (siapa yang boleh menggunakan token ini) - Opsional
		"exp": time.Now().Add(time.Hour * 1).Unix(), // Expiration time (1 jam dari sekarang)
		"iat": time.Now().Unix(),                    // Issued at
		"nbf": time.Now().Unix(),                    // Not before
		// Custom claims
		"user_id":  user.ID,
		"username": user.Username,
		"email":    user.Email,
	}

	// Buat token dengan claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Tandatangani token dengan secret key
	tokenString, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// validateToken memvalidasi token JWT dari header Authorization
func ValidateToken(encodedToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		// Validasi signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return token, nil
}
