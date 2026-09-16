package handlers

import (
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Indra-619/court-line/backend/internal/domain/repository"
)

// parseHexID validates a path/query identifier as a 24-character hex
// ObjectID and maps parse failures to repository.ErrInvalidID so
// handlers answer with the same 400s as before.
func parseHexID(id string) (primitive.ObjectID, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return primitive.NilObjectID, repository.ErrInvalidID
	}
	return objID, nil
}
