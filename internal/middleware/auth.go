package middleware

import (
	errs "cvwo-backend/internal/errors"
	"cvwo-backend/internal/models"
	"cvwo-backend/internal/services"
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
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing authorization header"})
		return
	}

	// Check if token is in valid format: "Bearer mytoken123"
	bearerToken := strings.Split(authHeader, " ")
	if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Malformed token"})
		return
	}

	// Check if token is valid and if so, retrieve the authenticated user
	user, err := authMiddleware.service.ValidateToken(bearerToken[1])
	if err != nil {
		errs.HTTPErrorResponse(ctx, err)
		ctx.Abort()
		return
	}

	// Store the authenticated user in context
	ctx.Set("user", user)

	ctx.Next()
}

// Retrieve the authenticated user from the context; error if not authenticated
func GetUserID(ctx *gin.Context) (uint, error) {
	value, exists := ctx.Get("user")
	if !exists {
		return 0, errs.New(errs.ErrUnauthorized, "Unauthenticated")
	}
	// Check if the "user" value is of the correct structure
	user, ok := value.(*models.User)
	if !ok {
		return 0, errs.New(errs.ErrUnauthorized, "Unauthenticated")
	}
	return user.ID, nil
}

// Retrieve the authenticated user from the context; no error if not authenticated
// userId of 0 indicates unauthenticated
func GetUserIDOrZero(ctx *gin.Context) (uint, error) {
	value, exists := ctx.Get("user")
	if !exists {
		return 0, nil
	}
	// Check if the "user" value is of the correct structure
	user, ok := value.(*models.User)
	if !ok {
		return 0, errs.New(errs.ErrUnauthorized, "Unauthenticated")
	}
	return user.ID, nil
}
