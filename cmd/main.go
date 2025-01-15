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

func main() {
	// Load .env file in development but not in production
	if os.Getenv("ENV") == "development" {
		if err := godotenv.Load(".env.development"); err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	// Initialize database
	db := data.InitDB(os.Getenv("DB_URL"))

	// Initialize application layers
	// Repositories (data access)
	userRepo := repos.NewUserRepo(db)
	postRepo := repos.NewPostRepo(db)
	commentRepo := repos.NewCommentRepo(db)
	topicRepo := repos.NewTopicRepo(db)
	postVoteRepo := repos.NewPostVoteRepo(db)
	commentVoteRepo := repos.NewCommentVoteRepo(db)

	// Services (business logic)
	userService := services.NewUserService(*userRepo)
	postService := services.NewPostService(*postRepo, *userRepo, *topicRepo)
	commentService := services.NewCommentService(*commentRepo, *postRepo, *userRepo)
	topicService := services.NewTopicService(*topicRepo)
	taggingService := services.NewTaggingService(*postRepo, *topicRepo)
	votingService := services.NewVotingService(*postVoteRepo, *commentVoteRepo)
	authService := services.NewAuthService(*userRepo)

	// Controllers (route handlers)
	userController := controllers.NewUserController(*userService)
	postController := controllers.NewPostController(*postService, *taggingService, *votingService)
	commentController := controllers.NewCommentController(*commentService, *votingService)
	topicController := controllers.NewTopicController(*topicService)
	authController := controllers.NewAuthController(authService)

	// Initialize router
	router := gin.Default()

	// Logger middleware to log request body of all requests
	router.Use(middleware.ResponseLogger)

	// CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{os.Getenv("FRONTEND_URL")},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Authentication middleware to validate jwt from requests
	authMiddleware := middleware.NewAuthMiddleware(authService)
	router.Use(authMiddleware.Authenticate())

	// Register route handlers
	routes.RegisterUserRoutes(router, userController)
	routes.RegisterPostRoutes(router, postController)
	routes.RegisterCommentRoutes(router, commentController)
	routes.RegisterTopicRoutes(router, topicController)
	routes.RegisterAuthRoutes(router, authController)

	router.Run()
}
