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

func RegisterPostRoutes(router *gin.Engine, controller *controllers.PostController) {
	router.GET("/posts", controller.GetAll)
	router.GET("/posts/:id", controller.GetByID)
	router.POST("/posts", controller.Create)
	router.PATCH("/posts/:id", controller.Update)
	router.DELETE("/posts/:id", controller.Delete)
}

func RegisterCommentRoutes(router *gin.Engine, controller *controllers.CommentController) {
	router.GET("/posts/:id/comments", controller.GetByPostID)
	router.POST("/comments", controller.Create)
	router.PATCH("/comments/:id", controller.Update)
	router.DELETE("/comments/:id", controller.Delete)
}

func RegisterTopicRoutes(router *gin.Engine, controller *controllers.TopicController) {
	router.GET("/topics", controller.GetAll)
	router.POST("/topics", controller.Create)
	router.DELETE("/topics/:id", controller.Delete)
}
