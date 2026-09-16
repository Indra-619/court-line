package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func Connect() error {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		return fmt.Errorf("MONGODB_URI environment variable not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	// Ping the database
	if err := client.Ping(ctx, nil); err != nil {
		return err
	}

	// Ensure a compound index on bookings for the double-booking lookup
	// (courtId + date). CreateOne is idempotent; failures are logged but
	// do not block startup (standalone dev Mongo).
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "courtId", Value: 1}, {Key: "date", Value: 1}},
		Options: options.Index().SetName("courtId_date"),
	}
	_, err = client.Database("booklapangan").Collection("bookings").Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		fmt.Printf("Warning: failed to create bookings index: %v\n", err)
	}

	Client = client
	fmt.Println("Connected to MongoDB!")
	return nil
}
