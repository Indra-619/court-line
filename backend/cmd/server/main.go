package main

import (
	"log"
	"os"

	"github.com/Indra-619/court-line/backend/internal/database"
	"github.com/Indra-619/court-line/backend/internal/handlers"
	"github.com/Indra-619/court-line/backend/internal/infrastructure"
	"github.com/Indra-619/court-line/backend/internal/middleware"
	"github.com/Indra-619/court-line/backend/internal/routes"
	"github.com/Indra-619/court-line/backend/pkg/config"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Validate required configuration before anything else starts
	config.MustLoad()

	// Connect to MongoDB
	client, err := database.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Wire repositories into handlers
	courtRepo := infrastructure.NewMongoCourtRepository(client, database.DBName)
	bookingRepo := infrastructure.NewMongoBookingRepository(client, database.DBName)
	userRepo := infrastructure.NewMongoUserRepository(client, database.DBName)
	refreshTokenRepo := infrastructure.NewMongoRefreshTokenRepository(client, database.DBName)
	revokedTokenRepo := infrastructure.NewMongoRevokedTokenRepository(client, database.DBName)

	// Back the admin role check with the injected user repository so
	// middleware never reaches for a global database handle.
	middleware.SetUserRoleLookup(userRepo)
	deps := routes.Deps{
		Courts:   handlers.NewCourtHandler(courtRepo),
		Bookings: handlers.NewBookingHandler(bookingRepo, courtRepo),
		Auth:     handlers.NewAuthHandler(userRepo, refreshTokenRepo, revokedTokenRepo),
	}

	// Setup router
	router := routes.SetupRouter(deps)

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
