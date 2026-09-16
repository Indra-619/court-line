package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Indra-619/court-line/backend/internal/database"
)

// userRoleLookup resolves the role for a user ID. It is a package-level
// variable so unit tests can stub out the MongoDB implementation.
var userRoleLookup = func(ctx context.Context, userID primitive.ObjectID) (string, error) {
	collection := database.Client.Database("booklapangan").Collection("users")
	var result struct {
		Role string `bson:"role"`
	}
	err := collection.FindOne(ctx, bson.M{"_id": userID}).Decode(&result)
	if err != nil {
		return "", err
	}
	return result.Role, nil
}

// AdminMiddleware restricts access to users with the "admin" role.
// It expects AuthMiddleware to have run first and set "userID" in the context.
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			c.Abort()
			return
		}

		userObjID, ok := userID.(primitive.ObjectID)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			c.Abort()
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
		defer cancel()

		role, err := userRoleLookup(ctx, userObjID)
		if err != nil || role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			c.Abort()
			return
		}

		c.Next()
	}
}
