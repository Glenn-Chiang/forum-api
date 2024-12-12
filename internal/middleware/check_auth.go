package middleware

import (
	"cvwo-backend/internal/services"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	service *services.AuthService
}

func NewAuthMiddleware(service *services.AuthService) *AuthMiddleware {
	return &AuthMiddleware{service}
}

// Middleware to check if the request is authenticated
func (authMiddleware *AuthMiddleware) CheckAuth(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")

	// Check if Authorization header is missing
	if authHeader == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
		return
	}

	// Check if token is in valid format: "Bearer mytoken123"
	bearerToken := strings.Split(authHeader, " ")
	if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	// Check if token is valid and if so, retrieve the authenticated user
	user, err := authMiddleware.service.ValidateToken(bearerToken[1]); 
	if err != nil {
		var authError *services.UnauthorizedError
		if errors.As(err, &authError) {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": authError.Error()})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// Store the authenticated user
	ctx.Set("current_user", user)

	// Pass control to handlers
	ctx.Next()
}

