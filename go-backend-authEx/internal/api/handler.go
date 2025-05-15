// internal/api/handler.go
package api

import (
	"log" // Kita akan ganti ini dengan logger terstruktur nanti
	"net/http"

	"go-auth-example/internal/model"
	"go-auth-example/internal/service"

	"github.com/gin-gonic/gin"
)

// AuthHandler struct (tetap sama)
type AuthHandler struct {
	authService service.AuthService
	userService service.UserService
}

// NewAuthHandler constructor (tetap sama)
func NewAuthHandler(auth service.AuthService, user service.UserService) *AuthHandler {
	return &AuthHandler{
		authService: auth,
		userService: user,
	}
}

// RegisterHandler menangani registrasi user baru
func (h *AuthHandler) RegisterHandler(c *gin.Context) {
	var input model.RegisterInput

	// Gunakan helper validasi baru
	validationErrors := ValidateAndBind(c, &input)
	if validationErrors != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": validationErrors})
		return
	}

	// Panggil service registrasi
	user, err := h.authService.Register(input)
	if err != nil {
		// Tangani error dari service
		if err.Error() == "email already registered" || err.Error() == "username already exists" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else {
			log.Printf("Unhandled registration error: %v", err) // Akan diganti dengan logger terstruktur
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully", "user": user})
}

// LoginHandler menangani login user
func (h *AuthHandler) LoginHandler(c *gin.Context) {
	var input model.LoginInput

	// Gunakan helper validasi baru
	validationErrors := ValidateAndBind(c, &input)
	if validationErrors != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": validationErrors})
		return
	}

	// Panggil service login
	token, err := h.authService.Login(input)
	if err != nil {
		if err.Error() == "invalid email or password" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		} else {
			log.Printf("Unhandled login error: %v", err) // Akan diganti dengan logger terstruktur
			c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred during login"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// ProfileHandler (tetap sama untuk saat ini)
func (h *AuthHandler) ProfileHandler(c *gin.Context) {
	userID, err := getUserIDFromContext(c) // Asumsi fungsi ini ada di middleware.go atau di-refactor
	if err != nil {
		log.Printf("Error getting userID from context in ProfileHandler: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not identify user"})
		return
	}

	user, err := h.userService.GetUserProfile(userID)
	if err != nil {
		if err.Error() == "user associated with token not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			log.Printf("Unhandled profile fetch error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user profile"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Welcome to your profile!", "user": user})
}

// (Pastikan fungsi getUserIDFromContext ada dan berfungsi, biasanya dari middleware.go)
// Jika belum ada, atau untuk sementara, Anda bisa copy dari contoh sebelumnya atau dari file yang ada.
// func getUserIDFromContext(c *gin.Context) (int, error) { ... }
