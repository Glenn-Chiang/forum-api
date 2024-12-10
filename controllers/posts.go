package controllers

import (
	"cvwo-backend/models"
	"cvwo-backend/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PostController struct {
	service services.PostService
}

func NewPostController(service services.PostService) *PostController {
	return &PostController{service}
}

// GET /posts
func (controller *PostController) GetAll(ctx *gin.Context) {
	posts, err := controller.service.GetAll()
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch posts"})
		return
	}
	ctx.IndentedJSON(http.StatusOK, posts)
}

// GET /posts/:id
func (controller *PostController) GetByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid post ID"})
	}

	post, err := controller.service.GetByID(uint(id))
	if err != nil {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": "post not found"})
		return
	}
	ctx.IndentedJSON(http.StatusOK, post)
}

// POST /posts
func (controller *PostController) Create(ctx *gin.Context) {
	var post models.Post

	// TODO: Parse and validate post data
	if err := ctx.BindJSON(&post); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid post data"})
		return
	}

	newPost, err := controller.service.Create(&post)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create post"})
	}

	ctx.IndentedJSON(http.StatusCreated, newPost)
}

// DELETE /posts/:id
func (controller *PostController) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid post ID"})		
	}

	if controller.service.Delete(uint(id)); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
