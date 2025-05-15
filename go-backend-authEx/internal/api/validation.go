// internal/api/validation.go
package api

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// ErrorMsg struct untuk pesan error validasi yang lebih user-friendly
type ErrorMsg struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// getErrorMsg mengembalikan pesan error yang lebih deskriptif berdasarkan tag validasi
func getErrorMsg(fe validator.FieldError) string {
	// Mengubah nama field JSON menjadi lebih ramah (misal: "Email" bukan "email")
	// Anda bisa membuat ini lebih canggih jika perlu, misal dengan mengambil tag `json:"..."`
	// fieldName := strings.Title(fe.Field()) // Ini hanya contoh sederhana

	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return fmt.Sprintf("Should be at least %s characters long", fe.Param())
	case "max":
		return fmt.Sprintf("Should be at most %s characters long", fe.Param())
	case "alphanum":
		return "Should only contain alphanumeric characters"
	// Tambahkan case lain sesuai kebutuhan tag validasi Anda
	default:
		return "Invalid value" // Pesan default
	}
}

// ValidateAndBind mem-bind JSON dan menjalankan validasi.
// Mengembalikan slice dari ErrorMsg jika ada error validasi.
func ValidateAndBind(c *gin.Context, input interface{}) []ErrorMsg {
	var errors []ErrorMsg

	// Bind JSON ke struct input
	if err := c.ShouldBindJSON(input); err != nil {
		// Handle error binding JSON yang fundamental (misal, JSON tidak valid)
		// Ini berbeda dari error validasi field.
		// Anda bisa mengembalikan satu ErrorMsg umum atau log error ini.
		// Untuk kesederhanaan, kita bisa tambahkan satu error umum.
		errors = append(errors, ErrorMsg{Field: "request_body", Message: "Invalid request body: " + err.Error()})
		return errors
	}

	// Lakukan validasi pada struct input
	err := validate.Struct(input)
	if err != nil {
		// Jika ada error validasi dari validator
		for _, fe := range err.(validator.ValidationErrors) {
			errors = append(errors, ErrorMsg{Field: strings.ToLower(fe.Field()), Message: getErrorMsg(fe)})
		}
		return errors
	}
	return nil // Tidak ada error
}
