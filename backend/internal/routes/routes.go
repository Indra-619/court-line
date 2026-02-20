package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/your-username/book-lapangan/backend/internal/handlers"
	"github.com/your-username/book-lapangan/backend/internal/middleware"
	"github.com/your-username/book-lapangan/backend/internal/presentation/api"
)

type RouterConfig struct {
	CourtHandler   *api.CourtHandler
	BookingHandler *api.BookingHandler
}

// SetupRouter configures all routes
func SetupRouter(cfg RouterConfig) *gin.Engine {
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

	// Auth routes (public) - Still using legacy handlers for now
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
	apiGroup := r.Group("/api")
	{
		// Courts - public read, protected write
		apiGroup.GET("/courts", cfg.CourtHandler.GetAllCourts)
		apiGroup.GET("/courts/:id", cfg.CourtHandler.GetCourtByID)
		apiGroup.GET("/courts/:id/bookings", cfg.BookingHandler.GetBookingsByCourt)

		// Protected court routes
		courtsProtected := apiGroup.Group("/courts")
		courtsProtected.Use(middleware.AuthMiddleware())
		{
			courtsProtected.POST("", cfg.CourtHandler.CreateCourt)
			courtsProtected.PUT("/:id", cfg.CourtHandler.UpdateCourt)
			courtsProtected.DELETE("/:id", cfg.CourtHandler.DeleteCourt)
		}

		// Bookings - all protected
		bookings := apiGroup.Group("/bookings")
		bookings.Use(middleware.AuthMiddleware())
		{
			bookings.POST("", cfg.BookingHandler.CreateBooking)
			bookings.GET("", cfg.BookingHandler.GetMyBookings)
		}
	}

	return r
}
