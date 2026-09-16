package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/Indra-619/court-line/backend/internal/handlers"
	"github.com/Indra-619/court-line/backend/internal/middleware"
)

// Deps carries the handler structs wired by main with their
// repositories. Routes only reference these; nothing reaches for a
// global database handle.
type Deps struct {
	Courts   *handlers.CourtHandler
	Bookings *handlers.BookingHandler
	Auth     *handlers.AuthHandler
}

// SetupRouter configures all routes
func SetupRouter(deps Deps) *gin.Engine {
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
		auth.GET("/google", deps.Auth.GoogleLogin)
		auth.GET("/google/callback", deps.Auth.GoogleCallback)
		auth.POST("/exchange", deps.Auth.ExchangeToken)
	}

	// Protected auth routes
	authProtected := r.Group("/auth")
	authProtected.Use(middleware.AuthMiddleware())
	{
		authProtected.GET("/me", deps.Auth.GetCurrentUser)
		authProtected.POST("/logout", deps.Auth.Logout)
	}

	// API routes
	api := r.Group("/api")
	{
		// Courts - public read, protected write
		api.GET("/courts", deps.Courts.GetCourts)
		api.GET("/courts/:id", deps.Courts.GetCourtByID)
		api.GET("/courts/:id/bookings", deps.Bookings.GetBookingsByCourtID)

		// Protected court routes (admin only)
		courtsProtected := api.Group("/courts")
		courtsProtected.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
		{
			courtsProtected.POST("", deps.Courts.CreateCourt)
			courtsProtected.PUT("/:id", deps.Courts.UpdateCourt)
			courtsProtected.DELETE("/:id", deps.Courts.DeleteCourt)
		}

		// Bookings - all protected
		bookings := api.Group("/bookings")
		bookings.Use(middleware.AuthMiddleware())
		{
			bookings.POST("", deps.Bookings.CreateBooking)
			bookings.GET("", deps.Bookings.GetBookings)
		}
	}

	return r
}
