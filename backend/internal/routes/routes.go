package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/your-username/book-lapangan/backend/internal/handlers"
	"github.com/your-username/book-lapangan/backend/internal/middleware"
)

// SetupRouter configures all routes
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// CORS configuration
	config := cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://frontend:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}
	r.Use(cors.New(config))

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK"})
	})

	// Auth routes (public)
	auth := r.Group("/auth")
	{
		auth.GET("/google", handlers.GoogleLogin)
		auth.GET("/google/callback", handlers.GoogleCallback)
	}

	// Protected auth routes
	authProtected := r.Group("/auth")
	authProtected.Use(middleware.AuthMiddleware())
	{
		authProtected.GET("/me", handlers.GetCurrentUser)
		authProtected.POST("/logout", handlers.Logout)
	}

	// API routes
	api := r.Group("/api")
	{
		// Courts - public read, protected write
		api.GET("/courts", handlers.GetCourts)
		api.GET("/courts/:id", handlers.GetCourtByID)
		api.GET("/courts/:id/bookings", handlers.GetBookingsByCourtID)

		// Protected court routes
		courtsProtected := api.Group("/courts")
		courtsProtected.Use(middleware.AuthMiddleware())
		{
			courtsProtected.POST("", handlers.CreateCourt)
			courtsProtected.PUT("/:id", handlers.UpdateCourt)
			courtsProtected.DELETE("/:id", handlers.DeleteCourt)
		}

		// Bookings - all protected
		bookings := api.Group("/bookings")
		bookings.Use(middleware.AuthMiddleware())
		{
			bookings.POST("", handlers.CreateBooking)
			bookings.GET("", handlers.GetBookings)
		}
	}

	return r
}
