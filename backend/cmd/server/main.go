package main

import (
	"log"
	"net/http"
	"os"

    "github.com/your-username/book-lapangan/backend/internal/database"
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

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
