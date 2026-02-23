package main

import (
	"log"
	"os"

	"github.com/your-username/book-lapangan/backend/internal/database"
	"github.com/your-username/book-lapangan/backend/internal/infrastructure/repository"
	"github.com/your-username/book-lapangan/backend/internal/presentation/api"
	"github.com/your-username/book-lapangan/backend/internal/routes"
	"github.com/your-username/book-lapangan/backend/internal/usecase"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 1. Connect to MongoDB
	if err := database.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 2. Initialize Repositories
	db := database.Client.Database("booklapangan")
	courtRepo := repository.NewCourtMongoRepository(db)
	bookingRepo := repository.NewBookingMongoRepository(db)

	// 3. Initialize Use Cases
	courtUseCase := usecase.NewCourtUseCase(courtRepo)
	bookingUseCase := usecase.NewBookingUseCase(bookingRepo, courtRepo)

	// 4. Initialize Handlers
	courtHandler := api.NewCourtHandler(courtUseCase)
	bookingHandler := api.NewBookingHandler(bookingUseCase)

	// 5. Setup router with dependencies
	router := routes.SetupRouter(routes.RouterConfig{
		CourtHandler:   courtHandler,
		BookingHandler: bookingHandler,
	})

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
