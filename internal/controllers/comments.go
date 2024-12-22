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
	commentService services.CommentService
	votingService  services.VotingService
}

func NewCommentController(commentService services.CommentService, votingService services.VotingService) *CommentController {
	return &CommentController{commentService, votingService}
}

// GET /posts/:id/comments
func (controller *CommentController) GetByPostID(ctx *gin.Context) {
	// Validate postId param
	postId, err := strconv.Atoi(ctx.Param("post_id"))
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

	// Retrieve the authenticated userID from context
	userID := middleware.GetUserIDOrZero(ctx)

	// Get the list of comments
	comments, totalCount, err := controller.commentService.GetByPostID(uint(postId), limit, offset, sortBy, userID)
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

	// Retrieve the authenticated userID from context
	userID, err := middleware.GetUserID(ctx)
	if err != nil {
		errs.HTTPErrorResponse(ctx, err)
		return
	}

	// Map fields from request body to Comment model
	comment := models.Comment{
		Content:  requestBody.Content,
		PostID:   requestBody.PostID,
		AuthorID: userID,
	}

	// Create new comment
	newComment, err := controller.commentService.Create(&comment)
	if err != nil {
		errs.HTTPErrorResponse(ctx, err)
		return
	}

	ctx.IndentedJSON(http.StatusCreated, newComment)
}

// PATCH /comments/:id
func (controller *CommentController) Update(ctx *gin.Context) {
	// Validate commentID
	id, err := strconv.Atoi(ctx.Param("comment_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	// Retrieve the authenticated userID from context
	userID, err := middleware.GetUserID(ctx)
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
	updatedComment, err := controller.commentService.Update(uint(id), requestBody.Content, userID)
	if err != nil {
		errs.HTTPErrorResponse(ctx, err)
		return
	}

	ctx.IndentedJSON(http.StatusOK, updatedComment)
}

// DELETE /comments/:id
func (controller *CommentController) Delete(ctx *gin.Context) {
	// Validate commentID
	id, err := strconv.Atoi(ctx.Param("comment_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	// Retrieve the authenticated userID from context
	userID, err := middleware.GetUserID(ctx)
	if err != nil {
		errs.HTTPErrorResponse(ctx, err)
		return
	}

	// Delete the comment
	if err := controller.commentService.Delete(uint(id), userID); err != nil {
		errs.HTTPErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

// PUT /comments/:comment_id/votes/:user_id
// Upvote or downvote a comment. The vote is associated with a particular user.
func (controller *CommentController) Vote(ctx *gin.Context) {
	// Validate comment_id param
	commentID, err := strconv.Atoi(ctx.Param("comment_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	// Validate user_id param
	userID, err := strconv.Atoi(ctx.Param("user_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
	}

	// Validate request body
	var requestBody models.VoteUpdate
	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Retrieve the authenticated userID from context
	authenticatedUserID, err := middleware.GetUserID(ctx)
	if err != nil {
		errs.HTTPErrorResponse(ctx, err)
		return
	}

	if err := controller.votingService.VoteComment(uint(commentID), uint(userID), requestBody.Value, authenticatedUserID); err != nil {
		errs.HTTPErrorResponse(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}
