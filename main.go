package main

import (
	"github.com/gin-gonic/gin"

	"cvwo-backend/controllers"
	"cvwo-backend/data"
)

const databaseURI = "index.db"
const serverUrl = "localhost:8080"

func main() {
	data.MustOpenDB(databaseURI)

	router := gin.Default()
	router.GET("/posts", controllers.GetPosts)
	router.GET("/posts/:id", controllers.GetPostByID)
	router.POST("/posts", controllers.CreatePost)

	router.Run(serverUrl)
}

