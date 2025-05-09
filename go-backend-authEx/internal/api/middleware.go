// internal/api/middleware.go
package api // <- Ubah package

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"go-auth-example/internal/auth" // <- Import auth package

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5" // <- Pindahkan import jwt ke sini jika getUserIDFromContext membutuhkannya
)

// AuthMiddleware adalah middleware Gin untuk otentikasi JWT
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			return
		}

		tokenString := parts[1]
		// Gunakan ValidateToken dari package auth
		token, err := auth.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		// Ambil User ID (pastikan sub adalah float64 lalu konversi ke int)
		userIDFloat, ok := claims["sub"].(float64)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
			return
		}
		userID := int(userIDFloat)

		c.Set("userID", userID) // Set user ID di context
		c.Next()
	}
}

// getUserIDFromContext helper untuk mendapatkan User ID dari context Gin
func getUserIDFromContext(c *gin.Context) (int, error) {
	idInterface, exists := c.Get("userID")
	if !exists {
		return 0, errors.New("user ID not found in context")
	}

	userID, ok := idInterface.(int)
	if !ok {
		// Coba konversi dari string jika perlu (tergantung bagaimana disimpan)
		idStr, okStr := idInterface.(string)
		if okStr {
			var errConv error
			userID, errConv = strconv.Atoi(idStr)
			if errConv != nil {
				return 0, errors.New("invalid user ID type in context (string conversion failed)")
			}
			return userID, nil
		}
		return 0, errors.New("invalid user ID type in context")
	}

	return userID, nil
}
