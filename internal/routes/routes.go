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
	router.GET("/posts/:id", controller.GetByID)
	router.PUT("/posts/:id/topics", auth.CheckAuth, controller.UpdateTags)
	router.POST("/posts", auth.CheckAuth, controller.Create)
	router.PATCH("/posts/:id", auth.CheckAuth, controller.Update)
	router.DELETE("/posts/:id", auth.CheckAuth, controller.Delete)
}

func RegisterCommentRoutes(router *gin.Engine, auth *middleware.AuthMiddleware, controller *controllers.CommentController) {
	router.GET("/posts/:id/comments", controller.GetByPostID)
	router.POST("/comments", auth.CheckAuth, controller.Create)
	router.PATCH("/comments/:id", auth.CheckAuth, controller.Update)
	router.DELETE("/comments/:id", auth.CheckAuth, controller.Delete)
}

func RegisterTopicRoutes(router *gin.Engine, auth *middleware.AuthMiddleware, controller *controllers.TopicController) {
	router.GET("/topics", controller.GetAll)
}
