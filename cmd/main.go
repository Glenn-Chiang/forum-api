package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"cvwo-backend/internal/controllers"
	"cvwo-backend/internal/data"
	"cvwo-backend/internal/models"
	"cvwo-backend/internal/repos"
	"cvwo-backend/internal/services"
)

const databaseURI = "index.db"
const serverUrl = "localhost:8080"

func main() {
	// Initialize database
	db := data.MustOpenDB(databaseURI)

	// Migrate schema
	if err := db.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{}); err != nil {
		log.Fatalf("Failed to migrate tables: %v", err)
	}

	// Seed database with initial data
	if err := data.SeedData(db); err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}

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

	// Configure routes

	// Posts
	router.GET("/posts", postController.GetAll)
	router.GET("/posts/:id", postController.GetByID)
	router.POST("/posts", postController.Create)
	router.DELETE("/posts/:id", postController.Delete)

	// Users
	router.GET("/users", userController.GetAll)
	router.GET("/users/:id", userController.GetByID)
	router.POST("/users", userController.Create)

	// Comments
	router.GET("/posts/:id/comments", commentController.GetByPostID)
	router.POST("/comments", commentController.Create)
	router.DELETE("/comments/:id", commentController.Delete)

	// Topics
	router.GET("/topics", topicController.GetAll)
	router.POST("/topics", topicController.Create)
	router.DELETE("/topics/:id", topicController.Delete)

	router.Run(serverUrl)
}
