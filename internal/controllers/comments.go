package controllers

import (
	errs "cvwo-backend/internal/errors"
	"cvwo-backend/internal/middleware"
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
	// Validate postId param
	postId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1 // If invalid, just set to default
	}

	// Limit refers to number of records per page
	// Get the "limit" query param and validate it
	limit, err := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10 // If invalid, just set to default
	}

	// Pagination offset: The DB will fetch {limit} number of records starting from the record at this index.
	offset := (page - 1) * limit

	// Get the "sortBy" query param and validate it
	sortBy := ctx.DefaultQuery("sortBy", "new")

	comments, err := controller.service.GetByPostID(uint(postId), limit, offset, sortBy)
	if err != nil {
		errs.HTTPErrorResponse(ctx, err)
		return
	}

	ctx.IndentedJSON(http.StatusOK, comments)
}

// POST /comments
func (controller *CommentController) Create(ctx *gin.Context) {
	// Validate request body
	var requestBody models.NewComment
	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Retrieve the authenticated user from context
	user, err := middleware.GetUserFromContext(ctx)
	if err != nil {
		errs.HTTPErrorResponse(ctx, err)
		return
	}

	// Check that the authorID of the comment corresponds to the currently authenticated user's ID
	if user.ID != requestBody.AuthorID {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Map fields from request body to Comment model
	comment := models.Comment{
		Content:  requestBody.Content,
		PostID:   requestBody.PostID,
		AuthorID: requestBody.AuthorID,
	}

	// Create new comment
	newComment, err := controller.service.Create(&comment)
	if err != nil {
		errs.HTTPErrorResponse(ctx, err)
		return
	}

	ctx.IndentedJSON(http.StatusCreated, newComment)
}

// PATCH /comments/:id
func (controller *CommentController) Update(ctx *gin.Context) {
	// Validate commentID
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	// Retrieve the authenticated user from context
	user, err := middleware.GetUserFromContext(ctx)
	if err != nil {
		errs.HTTPErrorResponse(ctx, err)
		return
	}

	// Fetch the comment to check its authorID
	comment, err := controller.service.GetByID(uint(id))
	if err != nil {
		errs.HTTPErrorResponse(ctx, err)
		return
	}

	// Check that the authorID of the comment corresponds to the currently authenticated user's ID
	if user.ID != comment.AuthorID {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Validate request body
	var requestBody models.CommentUpdate
	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the comment
	updatedComment, err := controller.service.Update(uint(id), requestBody.Content)
	if err != nil {
		errs.HTTPErrorResponse(ctx, err)
		return
	}

	ctx.IndentedJSON(http.StatusOK, updatedComment)
}

// DELETE /comments/:id
func (controller *CommentController) Delete(ctx *gin.Context) {
	// Validate commentID
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	// Retrieve the authenticated user from context
	user, err := middleware.GetUserFromContext(ctx)
	if err != nil {
		errs.HTTPErrorResponse(ctx, err)
		return
	}

	// Fetch the comment to check its authorID
	comment, err := controller.service.GetByID(uint(id))
	if err != nil {
		errs.HTTPErrorResponse(ctx, err)
		return
	}

	// Check that the authorID of the comment corresponds to the currently authenticated user's ID
	if user.ID != comment.AuthorID {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Delete the comment
	if err := controller.service.Delete(uint(id)); err != nil {
		errs.HTTPErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
