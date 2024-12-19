package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"cvwo-backend/internal/controllers"
	"cvwo-backend/internal/data"
	"cvwo-backend/internal/middleware"
	"cvwo-backend/internal/repos"
	"cvwo-backend/internal/routes"
	"cvwo-backend/internal/services"
)

const serverUrl = "localhost:8080"
const clientUrl = "*" // TODO: Set to frontend domain

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize database
	db := data.InitDB(os.Getenv("DB_URL"))

	// Initialize application layers
	// Repositories (data access)
	userRepo := repos.NewUserRepo(db)
	postRepo := repos.NewPostRepo(db)
	commentRepo := repos.NewCommentRepo(db)
	topicRepo := repos.NewTopicRepo(db)
	
	// Services (business logic)
	userService := services.NewUserService(*userRepo)
	postService := services.NewPostService(*postRepo, *userRepo, *topicRepo)
	commentService := services.NewCommentService(*commentRepo, *postRepo, *userRepo)
	topicService := services.NewTopicService(*topicRepo)
	taggingService := services.NewTaggingService(*postRepo, *topicRepo)
	authService := services.NewAuthService(userRepo)
	
	// Controllers (route handlers)
	userController := controllers.NewUserController(*userService)
	postController := controllers.NewPostController(*postService, *topicService, *taggingService)
	commentController := controllers.NewCommentController(*commentService)
	topicController := controllers.NewTopicController(*topicService)
	authController := controllers.NewAuthController(authService)

	// Authentication middleware to validate jwt from requests
	authMiddleware := middleware.NewAuthMiddleware(authService)

	// Initialize router
	router := gin.Default()

	// Logger middleware to log request body of all requests
	router.Use(middleware.ResponseLogger)

	// CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{clientUrl},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders: []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge: 12 * time.Hour,
	}))

	// Register route handlers
	routes.RegisterUserRoutes(router, authMiddleware, userController)
	routes.RegisterPostRoutes(router, authMiddleware, postController)
	routes.RegisterCommentRoutes(router, authMiddleware, commentController)
	routes.RegisterTopicRoutes(router, authMiddleware, topicController)
	routes.RegisterAuthRoutes(router, authController)

	router.Run(serverUrl)
}
