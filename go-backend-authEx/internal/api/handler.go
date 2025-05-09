// internal/api/handler.go
package api

import (
	"log"
	"net/http"

	"go-auth-example/internal/model"   // <- Import model
	"go-auth-example/internal/service" // <- Import service

	"github.com/gin-gonic/gin"
)

// AuthHandler struct untuk menampung dependencies handler
type AuthHandler struct {
	authService service.AuthService // Dependensi ke AuthService
	userService service.UserService // Dependensi ke UserService
}

// NewAuthHandler constructor untuk AuthHandler
func NewAuthHandler(auth service.AuthService, user service.UserService) *AuthHandler {
	return &AuthHandler{
		authService: auth,
		userService: user,
	}
}

// RegisterHandler menangani registrasi user baru
func (h *AuthHandler) RegisterHandler(c *gin.Context) {
	var input model.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	// Panggil service registrasi
	// Panggil service registrasi
	user, err := h.authService.Register(input)
	if err != nil {
		// Tangani error dari service
		// Cek error spesifik bisnis jika perlu (misal, email conflict)
		if err.Error() == "email already registered" || err.Error() == "username already exists" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else {
			// Error generik untuk masalah lain
			log.Printf("Unhandled registration error: %v", err) // Log error asli
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully", "user": user})
}

// LoginHandler menangani login user
func (h *AuthHandler) LoginHandler(c *gin.Context) {
	var input model.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	// Panggil service login
	token, err := h.authService.Login(input)
	if err != nil {
		// Tangani error dari service
		if err.Error() == "invalid email or password" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		} else {
			log.Printf("Unhandled login error: %v", err) // Log error asli
			c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred during login"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// ProfileHandler contoh handler untuk route yang dilindungi
func (h *AuthHandler) ProfileHandler(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		log.Printf("Error getting userID from context in ProfileHandler: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not identify user"})
		return
	}

	// Panggil service untuk get profile
	user, err := h.userService.GetUserProfile(userID)
	if err != nil {
		if err.Error() == "user associated with token not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			log.Printf("Unhandled profile fetch error: %v", err) // Log error asli
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user profile"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Welcome to your profile!", "user": user})
}
