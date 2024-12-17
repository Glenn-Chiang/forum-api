package controllers

import (
	"cvwo-backend/internal/models"
	"cvwo-backend/internal/services"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	service *services.AuthService
}

func NewAuthController(service *services.AuthService) *AuthController {
	return &AuthController{service}
}

// POST /login
func (controller *AuthController) Login(ctx *gin.Context) {
	var authInput models.AuthInput

	if err := ctx.ShouldBindJSON(&authInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, token, err := controller.service.Authenticate(&authInput)
	if err != nil {
		var authError *services.UnauthorizedError
		if errors.As(err, &authError) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": authError.Error()})
			return
		}
		// Otherwise return server error
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token, "user": user})
}
