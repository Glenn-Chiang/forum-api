package controllers

import (
	"cvwo-backend/internal/models"
	"cvwo-backend/internal/services"
	"errors"
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

// GET /posts or /posts?topic_id=1
func (controller *PostController) GetAll(ctx *gin.Context) {
	topicIdParam := ctx.Query("topic_id")

	var posts []models.Post
	var err error

	if topicIdParam != "" { // If topicId is specified
		// Parse topic_id from request query param
		topicID, convErr := strconv.Atoi(topicIdParam)
		if convErr != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid topic_id"})
			return
		}
		posts, err = controller.postService.GetByTopic(uint(topicID))
	} else { // If no topicId is specified
		posts, err = controller.postService.GetAll()
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.IndentedJSON(http.StatusOK, posts)
}

// GET /posts/:id
func (controller *PostController) GetByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid post ID"})
		return
	}

	post, err := controller.postService.GetByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
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

	// // Retrieve the authenticated user from context
	// user, exists := ctx.Get("user")
	// // This should not happen as middleware already checks for valid user
	// if !exists {
	// 	ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	// 	return
	// }

	// // Check that the authorID of the post corresponds to the currently authenticated user's ID
	// userID := user.(*models.User).ID
	// if userID != requestBody.AuthorID {
	// 	ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	// 	return
	// }

	// Get the topics associated with the list of topic IDs
	topics, err := controller.topicService.GetByIDs(requestBody.TopicIDs)
	// Handle errors with fetching topics
	if err != nil {
		var notFoundErr *services.NotFoundError
		if errors.As(err, &notFoundErr) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": notFoundErr.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Map fields from request body to Post model
	post := models.Post{
		Title:    requestBody.Title,
		Content:  requestBody.Content,
		AuthorID: requestBody.AuthorID,
		Topics:   topics,
	}

	newPost, err := controller.postService.Create(&post)

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

	ctx.IndentedJSON(http.StatusCreated, newPost)
}

// PATCH /posts/:id
func (controller *PostController) Update(ctx *gin.Context) {
	// Validate postID
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid post ID"})
		return
	}

	// Validate request body
	var requestBody models.PostUpdate
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

	// TODO: Check that the post's authorID corresponds to the currently authenticated user's ID
	userID := user.(*models.User).ID
	if userID != uint(id) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	updatedPost, err := controller.postService.Update(uint(id), requestBody.Title, requestBody.Content)

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

	ctx.IndentedJSON(http.StatusOK, updatedPost)
}

// PUT /posts/:id/topics
// Replace the list of topics associated with a post with a new list of topics
func (controller *PostController) UpdateTags(ctx *gin.Context) {
	// Validate postID
	postId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid post ID"})
		return
	}

	// Validate request body
	var requestBody models.TagsUpdate
	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := controller.taggingService.TagPostWithTopics(uint(postId), requestBody.TopicIDs); err != nil {
		var notFoundErr *services.NotFoundError
		if errors.As(err, &notFoundErr) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": notFoundErr.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

// DELETE /posts/:id
func (controller *PostController) Delete(ctx *gin.Context) {
	// Validate postID
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid post ID"})
		return
	}

	// Retrieve the authenticated user from context
	// user, exists := ctx.Get("user")
	// // This should not happen as middleware already checks for valid user
	// if !exists {
	// 	ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	// 	return
	// }

	// // Check that the authorID of the post corresponds to the currently authenticated user's ID
	// userID := user.(*models.User).ID
	// if userID != uint(id) { // TODO: Compare with authorID of post, not with post ID!
	// 	ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	// 	return
	// }

	if err := controller.postService.Delete(uint(id)); err != nil {
		var notFoundErr *services.NotFoundError
		if errors.As(err, &notFoundErr) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": notFoundErr.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
