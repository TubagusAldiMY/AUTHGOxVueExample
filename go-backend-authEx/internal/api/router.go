// internal/api/router.go
package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

// SetupRouter mengkonfigurasi dan mengembalikan instance Gin Engine
func SetupRouter(authHandler *AuthHandler) *gin.Engine {
	router := gin.Default()

	// --- Konfigurasi CORS ---
	router.Use(cors.New(cors.Config{
		// AllowOrigins adalah daftar origin yang diizinkan.
		// Untuk development, Anda bisa menggunakan alamat Vite dev server.
		// Untuk production, ganti dengan domain frontend Anda.
		// Menggunakan "*" akan mengizinkan semua origin (kurang aman untuk production).
		AllowOrigins:     []string{"http://localhost:5173"}, // Alamat default Vite dev server
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true, // Jika Anda perlu mengirim cookie atau header Authorization
		MaxAge:           12 * time.Hour,
	}))

	// Rute Publik
	router.POST("/register", authHandler.RegisterHandler)
	router.POST("/login", authHandler.LoginHandler)

	// Rute Terproteksi
	authorized := router.Group("/api")
	authorized.Use(AuthMiddleware())
	{
		authorized.GET("/profile", authHandler.ProfileHandler)
	}

	return router
}
