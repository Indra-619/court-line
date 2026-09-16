package handlers

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/Indra-619/court-line/backend/internal/domain/entity"
	"github.com/Indra-619/court-line/backend/internal/domain/repository"
	"github.com/Indra-619/court-line/backend/pkg/config"
)

var googleOauthConfig *oauth2.Config

const oauthStateCookie = "oauth_state"

// generateState returns a cryptographically random 32-byte hex string
// used as the OAuth state parameter.
func generateState() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

// validateState compares the state received from the provider with the
// expected value using a constant-time comparison. Both values must be
// non-empty.
func validateState(received, expected string) bool {
	if received == "" || expected == "" {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(received), []byte(expected)) == 1
}

func init() {
	googleOauthConfig = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}

// AuthHandler serves the /auth endpoints on top of injected
// repositories; it never touches the database driver directly.
type AuthHandler struct {
	users         repository.UserRepository
	refreshTokens repository.RefreshTokenRepository
	revokedTokens repository.RevokedTokenRepository
}

// NewAuthHandler builds an AuthHandler backed by the given
// repositories. refreshTokens stores refresh-token hashes;
// revokedTokens is the jti blacklist used by logout.
func NewAuthHandler(
	users repository.UserRepository,
	refreshTokens repository.RefreshTokenRepository,
	revokedTokens repository.RevokedTokenRepository,
) *AuthHandler {
	return &AuthHandler{
		users:         users,
		refreshTokens: refreshTokens,
		revokedTokens: revokedTokens,
	}
}

// GoogleUserInfo represents the user info from Google
type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

// GoogleLogin redirects to Google OAuth
func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	state, err := generateState()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate state"})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(oauthStateCookie, state, 600, "/", "", false, true)
	url := googleOauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// GoogleCallback handles the OAuth callback
func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	expectedState, err := c.Cookie(oauthStateCookie)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or missing OAuth state"})
		return
	}
	receivedState := c.Query("state")
	if !validateState(receivedState, expectedState) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or missing OAuth state"})
		return
	}
	// Clear the state cookie now that it has been validated
	c.SetCookie(oauthStateCookie, "", -1, "/", "", false, true)

	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Code not found"})
		return
	}

	token, err := googleOauthConfig.Exchange(ctx, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
		return
	}

	// Get user info from Google
	client := googleOauthConfig.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}
	defer resp.Body.Close()

	var googleUser GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode user info"})
		return
	}

	// Upsert user: refresh profile fields, or insert with the default role
	if err := h.users.UpsertGoogleUser(ctx, &entity.User{
		GoogleID: googleUser.ID,
		Email:    googleUser.Email,
		Name:     googleUser.Name,
		Picture:  googleUser.Picture,
		Role:     "user",
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save user"})
		return
	}

	// Get user from database
	user, err := h.users.FindByGoogleID(ctx, googleUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	// Issue a single-use exchange code and redirect to the frontend
	exchangeCode := createExchangeCode(user.ID)
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:3000"
	}
	c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s/auth/callback?code=%s", frontendURL, exchangeCode))
}

// ExchangeTokenInput represents the request body for the token exchange endpoint
type ExchangeTokenInput struct {
	Code string `json:"code" binding:"required"`
}

// ExchangeToken trades a one-time exchange code for a JWT
func (h *AuthHandler) ExchangeToken(c *gin.Context) {
	var input ExchangeTokenInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired code"})
		return
	}

	userID, ok := consumeExchangeCode(input.Code)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired code"})
		return
	}

	jwtToken, err := generateJWT(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	refreshToken, err := h.issueRefreshToken(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": gin.H{"token": jwtToken, "refreshToken": refreshToken}})
}

// refreshTokenTTL is how long an issued refresh token stays valid.
const refreshTokenTTL = 30 * 24 * time.Hour

// issueRefreshToken generates a fresh refresh token, persists only its
// SHA-256 hash, and returns the raw token for the client.
func (h *AuthHandler) issueRefreshToken(ctx context.Context, userID string) (string, error) {
	raw, err := generateRefreshToken()
	if err != nil {
		return "", err
	}
	now := time.Now()
	record := &entity.RefreshToken{
		UserID:    userID,
		TokenHash: sha256Hex(raw),
		ExpiresAt: now.Add(refreshTokenTTL),
		CreatedAt: now,
	}
	if err := h.refreshTokens.Create(ctx, record); err != nil {
		return "", err
	}
	return raw, nil
}

// GetCurrentUser returns the current authenticated user
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userObjID := userID.(primitive.ObjectID)

	user, err := h.users.FindByID(ctx, userObjID.Hex())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

// Logout handles user logout
func (h *AuthHandler) Logout(c *gin.Context) {
	// For JWT-based auth, logout is handled client-side by removing the token
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// generateJWT generates a JWT token for the user. Each token carries a
// unique jti so it can be blacklisted individually on logout.
func generateJWT(userID string) (string, error) {
	jti, err := generateJTI()
	if err != nil {
		return "", err
	}
	claims := jwt.MapClaims{
		"userId": userID,
		"jti":    jti,
		"exp":    time.Now().Add(time.Hour * 24).Unix(), // 24 hours
		"iat":    time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.JWTSecret()))
}
