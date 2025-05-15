// internal/api/handler.go
package api

import (
	"net/http"

	"go-auth-example/internal/logger" // Dari Tugas 1.3
	"go-auth-example/internal/model"
	"go-auth-example/internal/service" // <-- PASTIKAN IMPORT INI ADA DAN DIGUNAKAN

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus" // Dari Tugas 1.3
)

// AuthHandler struct untuk menampung dependencies handler
// PASTIKAN STRUCT INI ADA DAN DIEKSPOR (Huruf Awal Kapital)
type AuthHandler struct {
	authService service.AuthService // Menggunakan service.AuthService
	userService service.UserService // Menggunakan service.UserService
}

// NewAuthHandler constructor untuk AuthHandler
// PASTIKAN FUNGSI INI ADA DAN DIEKSPOR
func NewAuthHandler(auth service.AuthService, user service.UserService) *AuthHandler {
	return &AuthHandler{
		authService: auth,
		userService: user,
	}
}

// RegisterHandler menangani registrasi user baru
// PASTIKAN METHOD INI TERKAIT DENGAN POINTER KE AuthHandler (*AuthHandler)
func (h *AuthHandler) RegisterHandler(c *gin.Context) {
	var input model.RegisterInput
	// requestID := c.GetString("requestID") // Jika Anda memiliki middleware untuk request ID

	logFields := logrus.Fields{
		"handler": "RegisterHandler",
		// "request_id": requestID, // Contoh field tambahan
	}

	validationErrors := ValidateAndBind(c, &input) // Dari Tugas 1.1
	if validationErrors != nil {
		logger.Log.WithFields(logFields).Warnf("Validation failed for registration: %v", validationErrors)
		RespondWithValidationErrors(c, http.StatusBadRequest, validationErrors) // Dari Tugas 1.2
		return
	}
	logFields["email"] = input.Email
	logFields["username"] = input.Username

	user, err := h.authService.Register(input)
	if err != nil {
		switch err.Error() {
		case "email already registered":
			logger.Log.WithFields(logFields).Info("Registration attempt with existing email.")
			RespondWithError(c, NewAPIError(http.StatusConflict, ErrCodeEmailTaken, "The email address is already in use."))
		case "username already exists":
			logger.Log.WithFields(logFields).Info("Registration attempt with existing username.")
			RespondWithError(c, NewAPIError(http.StatusConflict, ErrCodeUsernameTaken, "The username is already taken."))
		default:
			logger.Log.WithFields(logFields).Errorf("Unhandled registration error: %v", err)
			RespondWithError(c, NewAPIError(http.StatusInternalServerError, ErrCodeInternalServer, "Failed to register user. Please try again later."))
		}
		return
	}

	logFields["user_id"] = user.ID
	logger.Log.WithFields(logFields).Info("User registered successfully")
	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully", "user": user})
}

// LoginHandler menangani login user
// PASTIKAN METHOD INI TERKAIT DENGAN POINTER KE AuthHandler (*AuthHandler)
func (h *AuthHandler) LoginHandler(c *gin.Context) {
	var input model.LoginInput
	// requestID := c.GetString("requestID")
	logFields := logrus.Fields{
		"handler": "LoginHandler",
		// "request_id": requestID,
		"email": input.Email,
	}

	validationErrors := ValidateAndBind(c, &input)
	if validationErrors != nil {
		logFields["email"] = input.Email
		logger.Log.WithFields(logFields).Warnf("Validation failed for login: %v", validationErrors)
		RespondWithValidationErrors(c, http.StatusBadRequest, validationErrors)
		return
	}
	logFields["email"] = input.Email

	token, err := h.authService.Login(input)
	if err != nil {
		switch err.Error() {
		case "invalid email or password":
			logger.Log.WithFields(logFields).Warn("Invalid login attempt.")
			RespondWithError(c, NewAPIError(http.StatusUnauthorized, ErrCodeInvalidCredentials, "Invalid email or password."))
		default:
			logger.Log.WithFields(logFields).Errorf("Unhandled login error: %v", err)
			RespondWithError(c, NewAPIError(http.StatusInternalServerError, ErrCodeInternalServer, "An error occurred during login. Please try again later."))
		}
		return
	}
	logger.Log.WithFields(logFields).Info("User logged in successfully")
	c.JSON(http.StatusOK, gin.H{"token": token})
}

// ProfileHandler contoh handler untuk route yang dilindungi
// PASTIKAN METHOD INI TERKAIT DENGAN POINTER KE AuthHandler (*AuthHandler)
func (h *AuthHandler) ProfileHandler(c *gin.Context) {
	// requestID := c.GetString("requestID")
	userID, errCtx := getUserIDFromContext(c)
	logFields := logrus.Fields{
		"handler": "ProfileHandler",
		// "request_id": requestID,
		"user_id": userID,
	}

	if errCtx != nil {
		logger.Log.WithFields(logFields).Errorf("Error getting userID from context in ProfileHandler: %v", errCtx)
		RespondWithError(c, NewAPIError(http.StatusInternalServerError, ErrCodeInternalServer, "Could not identify user."))
		return
	}

	user, err := h.userService.GetUserProfile(userID)
	if err != nil {
		logFields["user_id_queried"] = userID
		switch err.Error() {
		case "user associated with token not found":
			logger.Log.WithFields(logFields).Warn("User profile not found for user ID from token.")
			RespondWithError(c, NewAPIError(http.StatusNotFound, ErrCodeUserNotFound, "User profile not found."))
		default:
			logger.Log.WithFields(logFields).Errorf("Unhandled profile fetch error: %v", err)
			RespondWithError(c, NewAPIError(http.StatusInternalServerError, ErrCodeInternalServer, "Failed to fetch user profile. Please try again later."))
		}
		return
	}
	logger.Log.WithFields(logFields).Info("User profile fetched successfully")
	c.JSON(http.StatusOK, gin.H{"message": "Welcome to your profile!", "user": user})
}
