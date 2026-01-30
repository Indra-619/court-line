package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a user in the system
type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	GoogleID  string             `bson:"googleId" json:"googleId"`
	Email     string             `bson:"email" json:"email"`
	Name      string             `bson:"name" json:"name"`
	Picture   string             `bson:"picture" json:"picture"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}
