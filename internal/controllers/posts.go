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

type PostController struct {
	postService    services.PostService
	topicService   services.TopicService
	taggingService services.TaggingService
}

func NewPostController(postService services.PostService, topicService services.TopicService, taggingService services.TaggingService) *PostController {
	return &PostController{postService, topicService, taggingService}
}

// GET /posts or /posts?topic_id=1&page=1&limit=10
// Get a list of posts that is paginated, sorted, and filtered by topic
func (controller *PostController) GetList(ctx *gin.Context) {
	// Get the "page" query param and validate it
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

	var posts []models.Post

	topicIdParam := ctx.Query("topic_id")

	// If no topicId is specified
	if topicIdParam == "" {
		posts, err = controller.postService.GetList(limit, offset, sortBy)
		if err != nil {
			errs.HTTPErrorResponse(ctx, err)
			return
		}
		// If topicId is specified, check if it is valid
	} else {
		topicId, err := strconv.Atoi(topicIdParam)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid topic_id"})
			return
		}
		posts, err = controller.postService.GetByTopic(uint(topicId), limit, offset, sortBy)
		if err != nil {
			errs.HTTPErrorResponse(ctx, err)
			return
		}
	}

	ctx.IndentedJSON(http.StatusOK, posts)
}

// GET /posts/:id
func (controller *PostController) GetByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	post, err := controller.postService.GetByID(uint(id))
	if err != nil {
		errs.HTTPErrorResponse(ctx, err)
		return
	}
	ctx.IndentedJSON(http.StatusOK, post)
}

// POST /posts
func (controller *PostController) Create(ctx *gin.Context) {
	// Validate request body
	var requestBody models.NewPost
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

	// Check that the authorID of the post corresponds to the currently authenticated user's ID
	if user.ID != requestBody.AuthorID {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get the topics associated with the list of topic IDs
	topics, err := controller.topicService.GetByIDs(requestBody.TopicIDs)
	// Handle errors with fetching topics
	if err != nil {
		errs.HTTPErrorResponse(ctx, err)
		return
	}

	// Map fields from request body to Post model
	post := models.Post{
		Title:    requestBody.Title,
		Content:  requestBody.Content,
		AuthorID: requestBody.AuthorID,
		Topics:   topics,
	}

	// Create the post
	newPost, err := controller.postService.Create(&post)
	if err != nil {
		errs.HTTPErrorResponse(ctx, err)
		return
	}

	ctx.IndentedJSON(http.StatusCreated, newPost)
}

// PATCH /posts/:id
func (controller *PostController) Update(ctx *gin.Context) {
	// Validate postID
	postId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	// Validate request body
	var requestBody models.PostUpdate
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

	// Fetch the post to check its authorID
	post, err := controller.postService.GetByID(uint(postId))
	if err != nil {
		errs.HTTPErrorResponse(ctx, err)
		return
	}

	// Check that the post's authorID corresponds to the currently authenticated user's ID
	if user.ID != post.AuthorID {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Update the post
	updatedPost, err := controller.postService.Update(uint(postId), requestBody.Title, requestBody.Content)
	if err != nil {
		errs.HTTPErrorResponse(ctx, err)
		return
	}

	ctx.IndentedJSON(http.StatusOK, updatedPost)
}

// PUT /posts/:id/topics
// Replace the list of topics associated with a post with a new list of topics
func (controller *PostController) UpdateTags(ctx *gin.Context) {
	// Validate postID
	postId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	// Validate request body
	var requestBody models.PostTagsUpdate
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

	// Fetch the post to check its authorID
	post, err := controller.postService.GetByID(uint(postId))
	if err != nil {
		errs.HTTPErrorResponse(ctx, err)
		return
	}

	// Check that the post's authorID corresponds to the currently authenticated user's ID
	if user.ID != post.AuthorID {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Update the post tags
	if err := controller.taggingService.TagPostWithTopics(uint(postId), requestBody.TopicIDs); err != nil {
		errs.HTTPErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

// DELETE /posts/:id
func (controller *PostController) Delete(ctx *gin.Context) {
	// Validate postID
	postId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	// Retrieve the authenticated user from context
	user, err := middleware.GetUserFromContext(ctx)
	if err != nil {
		errs.HTTPErrorResponse(ctx, err)
		return
	}

	// Fetch the post to check its authorID
	post, err := controller.postService.GetByID(uint(postId))
	if err != nil {
		errs.HTTPErrorResponse(ctx, err)
		return
	}

	// Check that the post's authorID corresponds to the currently authenticated user's ID
	if user.ID != post.AuthorID {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Delete the post
	if err := controller.postService.Delete(uint(postId)); err != nil {
		errs.HTTPErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
