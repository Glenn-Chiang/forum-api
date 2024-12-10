package controllers

import (
	"cvwo-backend/models"
	"cvwo-backend/services"
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
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch users"})
		return
	}
	ctx.IndentedJSON(http.StatusOK, users)
}

// GET /users/:id
func (controller *UserController) GetByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
	}

	user, err := controller.service.GetByID(uint(id))
	if err != nil {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": "user not found"})
		return
	}
	ctx.IndentedJSON(http.StatusOK, user)
}

// POST /users
func (controller *UserController) Create(ctx *gin.Context) {
	var user models.User

	// TODO: Parse and validate user data
	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user data"})
		return
	}

	newUser, err := controller.service.Create(&user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
	}

	ctx.IndentedJSON(http.StatusCreated, newUser)
}

