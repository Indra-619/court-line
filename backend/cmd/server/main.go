package main

import (
	"log"
	"os"

	"github.com/your-username/book-lapangan/backend/internal/database"
	"github.com/your-username/book-lapangan/backend/internal/routes"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Connect to MongoDB
	if err := database.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Setup router
	router := routes.SetupRouter()

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
