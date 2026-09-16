package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User is the BSON persistence struct for the users collection; the
// corresponding domain struct is entity.User in internal/domain/entity.
type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	GoogleID  string             `bson:"googleId" json:"googleId"`
	Email     string             `bson:"email" json:"email"`
	Name      string             `bson:"name" json:"name"`
	Picture   string             `bson:"picture" json:"picture"`
	Role      string             `bson:"role" json:"role"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}
