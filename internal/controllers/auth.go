package controllers

import (
	"cvwo-backend/internal/models"
	"cvwo-backend/internal/services"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	service *services.AuthService
}

func NewAuthController(service *services.AuthService) *AuthController {
	return &AuthController{service}
}

// GET /login
func (controller *AuthController) Login(ctx *gin.Context) {
	var authInput models.AuthInput

	if err := ctx.ShouldBindJSON(&authInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := controller.service.Authenticate(&authInput)
	if err != nil {
		var authError *services.UnauthorizedError
		if errors.As(err, &authError) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": authError.Error()})
			return
		}
		// Otherwise return server error
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

// Middleware to check if the request is authenticated
func (controller *AuthController) CheckAuth(ctx *gin.Context) {
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
	user, err := controller.service.ValidateToken(bearerToken[1]); 
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
