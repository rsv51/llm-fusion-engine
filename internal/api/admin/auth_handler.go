package admin

import (
	"crypto/rand"
	"encoding/base64"
	"llm-fusion-engine/internal/database"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthHandler handles authentication requests.
type AuthHandler struct {
	db *gorm.DB
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{db: db}
}

// LoginRequest represents the login request payload.
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents the login response.
type LoginResponse struct {
	Token string      `json:"token"`
	User  UserProfile `json:"user"`
}

// UserProfile represents user profile information.
type UserProfile struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	IsAdmin  bool   `json:"isAdmin"`
}

// Login handles user login and returns a token.
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user by username
	var user database.User
	if err := h.db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate token (simple random token for now)
	token, err := generateToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Return user profile and token
	c.JSON(http.StatusOK, LoginResponse{
		Token: token,
		User: UserProfile{
			ID:       user.ID,
			Username: user.Username,
			IsAdmin:  user.IsAdmin,
		},
	})
}

// generateToken creates a random token for authentication.
func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// GetProfile returns the current user's profile.
func (h *AuthHandler) GetProfile(c *gin.Context) {
	// For now, return a mock profile
	// In production, this should verify the token and return actual user data
	c.JSON(http.StatusOK, UserProfile{
		ID:       1,
		Username: "admin",
		IsAdmin:  true,
	})
}