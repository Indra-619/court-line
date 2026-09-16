package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// RefreshToken is the BSON persistence struct for the refresh_tokens
// collection. TokenHash holds the SHA-256 hex digest of the raw token;
// the raw value is never stored. The domain struct is
// entity.RefreshToken in internal/domain/entity.
type RefreshToken struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    string             `bson:"userId" json:"userId"`
	TokenHash string             `bson:"tokenHash" json:"tokenHash"`
	ExpiresAt time.Time          `bson:"expiresAt" json:"expiresAt"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	Revoked   bool               `bson:"revoked" json:"revoked"`
}

// RevokedToken is the BSON persistence struct for the revoked_tokens
// collection: blacklisted JWT identifiers kept until the token they
// belong to expires (TTL index on expiresAt). The domain struct is
// entity.RevokedToken in internal/domain/entity.
type RevokedToken struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	JTI       string             `bson:"jti" json:"jti"`
	ExpiresAt time.Time          `bson:"expiresAt" json:"expiresAt"`
}
