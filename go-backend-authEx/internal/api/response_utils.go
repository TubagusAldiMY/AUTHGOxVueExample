// internal/api/response_utils.go
package api

import "github.com/gin-gonic/gin"

// APIError adalah struktur standar untuk respons error JSON
type APIError struct {
	HTTPStatus int         `json:"-"`                 // Status HTTP, tidak dikirim dalam JSON body
	Code       string      `json:"code"`              // Kode error aplikasi kustom
	Message    string      `json:"message"`           // Pesan error yang lebih user-friendly
	Details    interface{} `json:"details,omitempty"` // Detail tambahan, bisa berupa []ErrorMsg dari validasi, atau string lain
}

// Error mengimplementasikan interface error untuk APIError
func (e *APIError) Error() string {
	return e.Message
}

// NewAPIError membuat instance APIError baru
func NewAPIError(httpStatus int, code string, message string) *APIError {
	return &APIError{
		HTTPStatus: httpStatus,
		Code:       code,
		Message:    message,
	}
}

// RespondWithError mengirimkan respons error JSON yang terstandarisasi
func RespondWithError(c *gin.Context, err *APIError) {
	c.AbortWithStatusJSON(err.HTTPStatus, err)
}

// RespondWithValidationErrors mengirimkan respons error validasi yang terstandarisasi
// Menggunakan ErrorMsg dari validation.go
func RespondWithValidationErrors(c *gin.Context, httpStatus int, details []ErrorMsg) {
	apiErr := &APIError{
		HTTPStatus: httpStatus,
		Code:       ErrCodeValidationFailed,
		Message:    "Input validation failed. Please check the details.",
		Details:    details,
	}
	c.AbortWithStatusJSON(apiErr.HTTPStatus, apiErr)
}
