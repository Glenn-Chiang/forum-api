package controllers

import (
	"cvwo-backend/models"
	"cvwo-backend/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CommentController struct {
	service services.CommentService
}

func NewCommentController(service services.CommentService) *CommentController {
	return &CommentController{service}
}

// TODO: Get comments associated with a post

// POST /comments
func (controller *CommentController) Create(ctx *gin.Context) {
	var comment models.Comment

	// TODO: Parse and validate comment data
	if err := ctx.BindJSON(&comment); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid comment data"})
		return
	}

	newComment, err := controller.service.Create(&comment)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create comment"})
	}

	ctx.IndentedJSON(http.StatusCreated, newComment)
}

// DELETE /comments/:id
func (controller *CommentController) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid comment ID"})		
	}

	if controller.service.Delete(uint(id)); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

