package controllers

import (
	errs "cvwo-backend/internal/errors"
	"cvwo-backend/internal/models"
	"cvwo-backend/internal/services"
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
		errs.HTTPErrorResponse(ctx, err)
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token, "user": user})
}

