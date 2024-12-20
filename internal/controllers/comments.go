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
	sortBy := ctx.DefaultQuery("sort", "new")

	// Retrieve the authenticated user from context
	user, err := middleware.GetUserOrNil(ctx)
	if err != nil {
		errs.HTTPErrorResponse(ctx, err)
		return
	}

	// Get the list of comments
	comments, totalCount, err := controller.service.GetByPostID(uint(postId), limit, offset, sortBy, user.ID)
	if err != nil {
		errs.HTTPErrorResponse(ctx, err)
		return
	}

	// Send list of comments together with total count
	ctx.IndentedJSON(http.StatusOK, gin.H{"data": comments, "total_count": totalCount})
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
	user, err := middleware.GetUser(ctx)
	if err != nil {
		errs.HTTPErrorResponse(ctx, err)
		return
	}

	// Map fields from request body to Comment model
	comment := models.Comment{
		Content:  requestBody.Content,
		PostID:   requestBody.PostID,
		AuthorID: user.ID,
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
	user, err := middleware.GetUser(ctx)
	if err != nil {
		errs.HTTPErrorResponse(ctx, err)
		return
	}

	// Validate request body
	var requestBody models.CommentUpdate
	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the comment
	updatedComment, err := controller.service.Update(uint(id), requestBody.Content, user.ID)
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
	user, err := middleware.GetUser(ctx)
	if err != nil {
		errs.HTTPErrorResponse(ctx, err)
		return
	}

	// Delete the comment
	if err := controller.service.Delete(uint(id), user.ID); err != nil {
		errs.HTTPErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
