package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type post struct {
	ID string `json:"id"`
	Title string `json:"title"`
	Content string `json:"content"`
	Author string `json:"author"`
}

type user struct {
	ID string `json:"id"`
	Username string `json:"username"`
}

var posts = []post{
	{ID: "1", Title: "What Could Have Been"},
	{ID: "2", Title: "My Dearest Enemy"},
	{ID: "3", Title: "The Glorious Evolution"},
}

func main() {
	router := gin.Default()
	router.GET("/posts", getPosts)
	router.GET("/posts/:id", getPostByID)
	router.POST("/posts", createPost)

	router.Run("localhost:8080")
}

func getPosts(ctx *gin.Context) {
	ctx.IndentedJSON(http.StatusOK, posts)
}

func createPost(ctx *gin.Context) {
	var newPost post

	if err := ctx.BindJSON(&newPost); err != nil {
		return
	}

	posts = append(posts, newPost)
	ctx.IndentedJSON(http.StatusCreated, newPost)
}

func getPostByID(ctx *gin.Context) {
	id := ctx.Param("id")

	for _, post := range posts {
		if post.ID == id {
			ctx.IndentedJSON(http.StatusOK, post)
			return
		}
	}
	ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": "post not found"})
}

