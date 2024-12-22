package routes

import (
	"cvwo-backend/internal/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(router *gin.Engine, controller *controllers.UserController) {
	router.GET("/users", controller.GetAll)
	router.GET("/users/:id", controller.GetByID)
	router.POST("/users", controller.Create)
}

func RegisterAuthRoutes(router *gin.Engine, controller *controllers.AuthController) {
	router.POST("/login", controller.Login)
}

func RegisterPostRoutes(router *gin.Engine, controller *controllers.PostController) {
	router.GET("/posts", controller.GetList)
	// Get individual post
	router.GET("/posts/:post_id", controller.GetByID)
	// Create new post
	router.POST("/posts", controller.Create)
	// Update post title and content
	router.PATCH("/posts/:post_id", controller.Update)
	// Update post tags
	router.PUT("/posts/:post_id/topics", controller.UpdateTags)
	// Delete post
	router.DELETE("/posts/:post_id", controller.Delete)
	// Upvote/downvote post
	router.PUT("/posts/:post_id/votes/:user_id", controller.Vote)
}

func RegisterCommentRoutes(router *gin.Engine, controller *controllers.CommentController) {
	router.GET("/posts/:post_id/comments", controller.GetByPostID)
	router.POST("/comments", controller.Create)
	router.PATCH("/comments/:id", controller.Update)
	router.DELETE("/comments/:id", controller.Delete)
}

func RegisterTopicRoutes(router *gin.Engine, controller *controllers.TopicController) {
	router.GET("/topics", controller.GetAll)
}
