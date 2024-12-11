package controllers

import (
	"cvwo-backend/internal/models"
	"cvwo-backend/internal/services"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	service services.UserService
}

func NewUserController(service services.UserService) *UserController {
	return &UserController{service}
}

// GET /users
func (controller *UserController) GetAll(ctx *gin.Context) {
	users, err := controller.service.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.IndentedJSON(http.StatusOK, users)
}

// GET /users/:id
func (controller *UserController) GetByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	user, err := controller.service.GetByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.IndentedJSON(http.StatusOK, user)
}

// POST /users
func (controller *UserController) Create(ctx *gin.Context) {
	var user models.AuthInput

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newUser, err := controller.service.Create(&user)
	if err != nil {
		// Check if error is due to username already in use
		var conflictErr *services.AlreadyInUseError
		if errors.As(err, &conflictErr) {
			ctx.JSON(http.StatusConflict, gin.H{"error": conflictErr.Error()})
			return
		}
		// Otherwise return server error
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.IndentedJSON(http.StatusCreated, newUser)
}
