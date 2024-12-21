package routes

import (
	"cvwo-backend/internal/controllers"
	"cvwo-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(router *gin.Engine, auth *middleware.AuthMiddleware, controller *controllers.UserController) {
	router.GET("/users", controller.GetAll)
	router.GET("/users/:id", controller.GetByID)
	router.POST("/users", controller.Create)
}

func RegisterAuthRoutes(router *gin.Engine, controller *controllers.AuthController) {
	router.POST("/login", controller.Login)
}

func RegisterPostRoutes(router *gin.Engine, auth *middleware.AuthMiddleware, controller *controllers.PostController) {
	router.GET("/posts", controller.GetList)
	// Get individual post
	router.GET("/posts/:post_id", controller.GetByID)
	// Create new post
	router.POST("/posts", auth.CheckAuth, controller.Create)
	// Update post title and content
	router.PATCH("/posts/:post_id", auth.CheckAuth, controller.Update)
	// Update post tags
	router.PUT("/posts/:post_id/topics", auth.CheckAuth, controller.UpdateTags)
	// Delete post
	router.DELETE("/posts/:post_id", auth.CheckAuth, controller.Delete)
	// Upvote/downvote post
	router.PUT("/posts/:post_id/votes/:user_id", auth.CheckAuth, controller.Vote)
	// Remove user's vote for post
	router.DELETE("/posts/:post_id/votes/:user_id", auth.CheckAuth, controller.DeleteVote)
}

func RegisterCommentRoutes(router *gin.Engine, auth *middleware.AuthMiddleware, controller *controllers.CommentController) {
	router.GET("/posts/:post_id/comments", controller.GetByPostID)
	router.POST("/comments", auth.CheckAuth, controller.Create)
	router.PATCH("/comments/:id", auth.CheckAuth, controller.Update)
	router.DELETE("/comments/:id", auth.CheckAuth, controller.Delete)
}

func RegisterTopicRoutes(router *gin.Engine, auth *middleware.AuthMiddleware, controller *controllers.TopicController) {
	router.GET("/topics", controller.GetAll)
}
