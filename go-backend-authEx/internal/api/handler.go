// internal/api/handler.go
package api

import (
	"log" // Akan kita ganti dengan logger terstruktur
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

	validationErrors := ValidateAndBind(c, &input) // Dari Tugas 1.1
	if validationErrors != nil {
		// Gunakan helper baru untuk error validasi
		RespondWithValidationErrors(c, http.StatusBadRequest, validationErrors)
		return
	}

	user, err := h.authService.Register(input)
	if err != nil {
		// Tangani error dari service dengan struktur error baru
		// Kita akan membuat error lebih spesifik dari service nanti,
		// untuk sekarang kita petakan berdasarkan string error.
		switch err.Error() { // Ini masih sederhana, idealnya service mengembalikan error yang lebih terstruktur
		case "email already registered":
			RespondWithError(c, NewAPIError(http.StatusConflict, ErrCodeEmailTaken, "The email address is already in use."))
		case "username already exists":
			RespondWithError(c, NewAPIError(http.StatusConflict, ErrCodeUsernameTaken, "The username is already taken."))
		default:
			log.Printf("Unhandled registration error: %v", err) // Akan diganti
			RespondWithError(c, NewAPIError(http.StatusInternalServerError, ErrCodeInternalServer, "Failed to register user. Please try again later."))
		}
		return
	}

	// Untuk sukses, kita bisa tetap menggunakan c.JSON atau buat helper juga
	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully", "user": user})
}

// LoginHandler menangani login user
func (h *AuthHandler) LoginHandler(c *gin.Context) {
	var input model.LoginInput

	validationErrors := ValidateAndBind(c, &input) // Dari Tugas 1.1
	if validationErrors != nil {
		RespondWithValidationErrors(c, http.StatusBadRequest, validationErrors)
		return
	}

	token, err := h.authService.Login(input)
	if err != nil {
		switch err.Error() { // Sama seperti di atas, idealnya error dari service lebih terstruktur
		case "invalid email or password":
			RespondWithError(c, NewAPIError(http.StatusUnauthorized, ErrCodeInvalidCredentials, "Invalid email or password."))
		default:
			log.Printf("Unhandled login error: %v", err) // Akan diganti
			RespondWithError(c, NewAPIError(http.StatusInternalServerError, ErrCodeInternalServer, "An error occurred during login. Please try again later."))
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// ProfileHandler contoh handler untuk route yang dilindungi
func (h *AuthHandler) ProfileHandler(c *gin.Context) {
	userID, errCtx := getUserIDFromContext(c)
	if errCtx != nil {
		log.Printf("Error getting userID from context in ProfileHandler: %v", errCtx) // Akan diganti
		// Ini error dari middleware, jadi bisa langsung pakai error code dari middleware
		RespondWithError(c, NewAPIError(http.StatusInternalServerError, ErrCodeInternalServer, "Could not identify user."))
		return
	}

	user, err := h.userService.GetUserProfile(userID)
	if err != nil {
		switch err.Error() { // Sama seperti di atas
		case "user associated with token not found":
			RespondWithError(c, NewAPIError(http.StatusNotFound, ErrCodeUserNotFound, "User profile not found."))
		default:
			log.Printf("Unhandled profile fetch error: %v", err) // Akan diganti
			RespondWithError(c, NewAPIError(http.StatusInternalServerError, ErrCodeInternalServer, "Failed to fetch user profile. Please try again later."))
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Welcome to your profile!", "user": user})
}
