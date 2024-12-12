package controllers

import (
	"cvwo-backend/internal/models"
	"cvwo-backend/internal/services"
	"errors"
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
		return
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

	// Retrieve the authenticated user from context
	user, exists := ctx.Get("user") 
	// This should not happen as middleware already checks for valid user
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Check that the authorID of the comment corresponds to the currently authenticated user's ID
	userID := user.(*models.User).ID
	if userID != requestBody.AuthorID {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Map fields from request body to Comment model
	comment := models.Comment{
		Content: requestBody.Content,
		PostID: requestBody.PostID,
		AuthorID: requestBody.AuthorID,
	}

	newComment, err := controller.service.Create(&comment)

	// Handle errors
	if err != nil {
		var notFoundErr *services.NotFoundError
		if errors.As(err, &notFoundErr) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": notFoundErr.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.IndentedJSON(http.StatusCreated, newComment)
}

// PATCH /comments/:id
func (controller *CommentController) Update(ctx *gin.Context) {
	// Validate commentID
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid comment ID"})
		return	
	}

	// Retrieve the authenticated user from context
	user, exists := ctx.Get("user") 
	// This should not happen as middleware already checks for valid user
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Check that the authorID of the comment corresponds to the currently authenticated user's ID
	userID := user.(*models.User).ID
	if userID != uint(id) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Validate request body
	var requestBody models.UpdateCommentRequest
	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedComment, err := controller.service.Update(uint(id), requestBody.Content)
	
	// Handle errors
	if err != nil {
		var notFoundErr *services.NotFoundError
		if errors.As(err, &notFoundErr) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": notFoundErr.Error()})
			return
		}
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
		return	
	}

	// Retrieve the authenticated user from context
	user, exists := ctx.Get("user") 
	// This should not happen as middleware already checks for valid user
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Check that the authorID of the comment corresponds to the currently authenticated user's ID
	userID := user.(*models.User).ID
	if userID != uint(id) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err:= controller.service.Delete(uint(id)); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

