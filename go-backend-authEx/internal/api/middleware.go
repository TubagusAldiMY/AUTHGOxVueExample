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

// AuthMiddleware (bagian error handlingnya)
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// Gunakan helper error baru
			RespondWithError(c, NewAPIError(http.StatusUnauthorized, ErrCodeMissingAuthHeader, "Authorization header is required."))
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			// Gunakan helper error baru
			RespondWithError(c, NewAPIError(http.StatusUnauthorized, ErrCodeInvalidAuthHeader, "Authorization header format must be Bearer {token}."))
			return
		}

		tokenString := parts[1]
		token, err := auth.ValidateToken(tokenString) // Dari internal/auth
		if err != nil {
			// Error dari ValidateToken mungkin perlu dipetakan ke kode error kita
			// Contoh sederhana:
			errMsg := "Invalid or expired token."
			errCode := ErrCodeTokenInvalid
			if strings.Contains(err.Error(), "expired") { // Ini cara deteksi sederhana, bisa lebih baik
				errCode = ErrCodeTokenExpired
			}
			RespondWithError(c, NewAPIError(http.StatusUnauthorized, errCode, errMsg))
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid { // Sebenarnya !token.Valid sudah dicek oleh jwt.Parse
			RespondWithError(c, NewAPIError(http.StatusUnauthorized, ErrCodeTokenInvalid, "Invalid token claims."))
			return
		}

		userIDFloat, ok := claims["sub"].(float64)
		if !ok {
			RespondWithError(c, NewAPIError(http.StatusUnauthorized, ErrCodeTokenInvalid, "Invalid user ID format in token."))
			return
		}
		userID := int(userIDFloat)

		c.Set("userID", userID)
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
