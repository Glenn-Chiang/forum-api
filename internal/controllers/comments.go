package controllers

import (
	"cvwo-backend/internal/models"
	"cvwo-backend/internal/services"
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

// GET /posts/:id/comments
func (controller *CommentController) GetByPostID(ctx *gin.Context) {
	postId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid post ID"})
	}

	comments, err := controller.service.GetByPostID(uint(postId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch comments"})
		return
	}

	ctx.IndentedJSON(http.StatusOK, comments)
}

// POST /comments
func (controller *CommentController) Create(ctx *gin.Context) {
	// Validate request body
	var requestBody models.CreateCommentRequest
	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Map fields from request body to Comment model
	comment := models.Comment{
		Content: requestBody.Content,
		PostID: requestBody.PostID,
		AuthorID: requestBody.AuthorID,
	}

	newComment, err := controller.service.Create(&comment)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	ctx.IndentedJSON(http.StatusCreated, newComment)
}

// PATCH /comments/:id
func (controller *CommentController) Update(ctx *gin.Context) {
	// Validate commentID
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid comment ID"})		
	}

	// Validate request body
	var requestBody models.UpdateCommentRequest
	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedComment, err := controller.service.Update(uint(id), requestBody.Content)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.IndentedJSON(http.StatusOK, updatedComment)
}

// DELETE /comments/:id
func (controller *CommentController) Delete(ctx *gin.Context) {
	// Validate commentID
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

