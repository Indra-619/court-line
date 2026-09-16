package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Indra-619/court-line/backend/internal/domain/repository"
	"github.com/Indra-619/court-line/backend/pkg/config"
)

// tokenBlacklist reports whether a JWT identifier (jti) has been
// revoked via logout. It is a package-level variable so unit tests can
// stub it and SetTokenBlacklist can bind it to the injected revoked
// token repository at startup. The default treats every token as live
// so middleware without wiring behaves exactly as before.
var tokenBlacklist = func(ctx context.Context, jti string) (bool, error) {
	return false, nil
}

// SetTokenBlacklist backs the jti revocation check with a
// RevokedTokenRepository. main calls it once during wiring.
func SetTokenBlacklist(revoked repository.RevokedTokenRepository) {
	tokenBlacklist = revoked.Exists
}

// AuthMiddleware validates JWT tokens
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Check Bearer format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Parse and validate token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(config.JWTSecret()), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		userIDStr, ok := claims["userId"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
			c.Abort()
			return
		}

		userObjID, err := primitive.ObjectIDFromHex(userIDStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID format"})
			c.Abort()
			return
		}

		// Reject tokens whose jti was blacklisted by a logout. When the
		// blacklist cannot be reached the request fails closed: it is
		// safer to refuse a valid token than to honour a revoked one.
		if jti, ok := claims["jti"].(string); ok && jti != "" {
			revoked, err := tokenBlacklist(c.Request.Context(), jti)
			if err != nil {
				log.Printf("auth middleware: token blacklist lookup failed: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Authentication service unavailable"})
				c.Abort()
				return
			}
			if revoked {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token revoked"})
				c.Abort()
				return
			}
			c.Set("jti", jti)
			if exp, ok := claims["exp"].(float64); ok {
				c.Set("tokenExp", time.Unix(int64(exp), 0))
			}
		}

		// Set user ID in context
		c.Set("userID", userObjID)
		c.Next()
	}
}

// OptionalAuthMiddleware tries to authenticate but doesn't require it
func OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Next()
			return
		}

		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(config.JWTSecret()), nil
		})

		if err != nil || !token.Valid {
			c.Next()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.Next()
			return
		}

		userIDStr, ok := claims["userId"].(string)
		if !ok {
			c.Next()
			return
		}

		userObjID, err := primitive.ObjectIDFromHex(userIDStr)
		if err != nil {
			c.Next()
			return
		}

		c.Set("userID", userObjID)
		c.Next()
	}
}
