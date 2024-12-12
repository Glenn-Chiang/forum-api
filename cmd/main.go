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
const clientUrl = "*"

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize database
	db := data.InitDB(os.Getenv("DB_URL"))

	// Initialize feature layers
	// Users
	userRepo := repos.NewUserRepo(db)
	userService := services.NewUserService(*userRepo)
	userController := controllers.NewUserController(*userService)
	
	// Posts
	postRepo := repos.NewPostRepo(db)
	postService := services.NewPostService(*postRepo, *userRepo)
	postController := controllers.NewPostController(*postService)
	
	// Comments
	commentRepo := repos.NewCommentRepo(db)
	commentService := services.NewCommentService(*commentRepo, *postRepo, *userRepo)
	commentController := controllers.NewCommentController(*commentService)
	
	// Topics
	topicRepo := repos.NewTopicRepo(db)
	topicService := services.NewTopicService(*topicRepo)
	topicController := controllers.NewTopicController(*topicService)

	// Authentication
	authService := services.NewAuthService(userService)
	authController := controllers.NewAuthController(authService)
	authMiddleware := middleware.NewAuthMiddleware(authService)

	// Initialize router
	router := gin.Default()

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
