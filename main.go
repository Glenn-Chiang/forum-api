package main

import (
	"github.com/gin-gonic/gin"

	"cvwo-backend/controllers"
)


func main() {
	router := gin.Default()
	router.GET("/posts", controllers.GetPosts)
	router.GET("/posts/:id", controllers.GetPostByID)
	router.POST("/posts", controllers.CreatePost)

	router.Run("localhost:8080")
}

