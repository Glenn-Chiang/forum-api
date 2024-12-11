package main

import (
	"github.com/gin-gonic/gin"

	"cvwo-backend/internal/controllers"
	"cvwo-backend/internal/data"
	"cvwo-backend/internal/repos"
	"cvwo-backend/internal/routes"
	"cvwo-backend/internal/services"
)

const databaseURI = "index.db"
const serverUrl = "localhost:8080"

func main() {
	// Initialize database
	db := data.InitDB(databaseURI)

	// Initialize feature layers
	// Users
	userRepo := repos.NewUserRepo(db)
	userService := services.NewUserService(*userRepo)
	userController := controllers.NewUserController(*userService)

	// Posts
	postRepo := repos.NewPostRepo(db)
	postService := services.NewPostService(*postRepo)
	postController := controllers.NewPostController(*postService)

	// Comments
	commentRepo := repos.NewCommentRepo(db)
	commentService := services.NewCommentService(*commentRepo)
	commentController := controllers.NewCommentController(*commentService)

	// Topics
	topicRepo := repos.NewTopicRepo(db)
	topicService := services.NewTopicService(*topicRepo)
	topicController := controllers.NewTopicController(*topicService)

	// Initialize router
	router := gin.Default()

	routes.RegisterUserRoutes(router, userController)
	routes.RegisterPostRoutes(router, postController)
	routes.RegisterCommentRoutes(router, commentController)
	routes.RegisterTopicRoutes(router, topicController)

	router.Run(serverUrl)
}
