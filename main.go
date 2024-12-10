package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"cvwo-backend/controllers"
	"cvwo-backend/data"
	"cvwo-backend/models"
	"cvwo-backend/repos"
	"cvwo-backend/services"
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

	// Initialize layers
	postRepo := repos.NewPostRepo(db)
	postService := services.NewPostService(*postRepo)
	postController := controllers.NewPostController(*postService)

	// Configure routes
	router := gin.Default()
	router.GET("/posts", postController.GetAll)
	router.GET("/posts/:id", postController.GetByID)
	router.POST("/posts", postController.Create)

	router.Run(serverUrl)
}
