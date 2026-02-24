package handlers

import (
	"github.com/bac-unified/api/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	user, err := h.authService.Register(req.Username, req.Email, req.Password)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	token, _ := h.authService.GenerateToken(user.ID)
	c.JSON(201, gin.H{"user": user, "token": token})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Simplified - would verify password
	userID := uuid.New()
	token, _ := h.authService.GenerateToken(userID)
	c.JSON(200, gin.H{"token": token})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	c.JSON(200, gin.H{"token": "new_token"})
}

func (h *AuthHandler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.JSON(401, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}
		// Simplified - would validate JWT
		c.Set("user_id", uuid.New())
		c.Next()
	}
}
