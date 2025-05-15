// internal/api/error_codes.go
package api

// Daftar kode error aplikasi untuk konsistensi antara backend dan frontend
const (
	// General Errors
	ErrCodeBadRequest       = "BAD_REQUEST"
	ErrCodeInternalServer   = "INTERNAL_SERVER_ERROR"
	ErrCodeUnauthorized     = "UNAUTHORIZED"
	ErrCodeNotFound         = "NOT_FOUND"
	ErrCodeValidationFailed = "VALIDATION_FAILED"

	// Auth Specific Errors
	ErrCodeEmailTaken         = "AUTH_EMAIL_TAKEN"
	ErrCodeUsernameTaken      = "AUTH_USERNAME_TAKEN"
	ErrCodeInvalidCredentials = "AUTH_INVALID_CREDENTIALS"
	ErrCodeUserNotFound       = "AUTH_USER_NOT_FOUND" // Bisa digunakan jika profil tidak ditemukan
	ErrCodeTokenExpired       = "AUTH_TOKEN_EXPIRED"
	ErrCodeTokenInvalid       = "AUTH_TOKEN_INVALID"
	ErrCodeMissingAuthHeader  = "AUTH_MISSING_HEADER"
	ErrCodeInvalidAuthHeader  = "AUTH_INVALID_HEADER"
)
