package routes

import (
	"cvwo-backend/internal/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(router *gin.Engine, controller *controllers.UserController) {
	router.GET("/posts", controller.GetAll)
	router.GET("/posts/:id", controller.GetByID)
	router.POST("/posts", controller.Create)
}

func RegisterPostRoutes(router *gin.Engine, controller *controllers.PostController) {
	router.GET("/posts", controller.GetAll)
	router.GET("/posts/:id", controller.GetByID)
	router.POST("/posts", controller.Create)
	router.DELETE("/posts/:id", controller.Delete)
}

func RegisterCommentRoutes(router *gin.Engine, controller *controllers.CommentController) {
	router.POST("/posts", controller.Create)
	router.DELETE("/posts/:id", controller.Delete)
}

func RegisterTopicRoutes(router *gin.Engine, controller *controllers.TopicController) {
	router.GET("/posts", controller.GetAll)
	router.POST("/posts", controller.Create)
	router.DELETE("/posts/:id", controller.Delete)
}
