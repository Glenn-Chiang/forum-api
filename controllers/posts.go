package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"cvwo-backend/models"
	"cvwo-backend/services"
)

func GetPosts(ctx *gin.Context) {
	posts := services.GetPosts()
	ctx.IndentedJSON(http.StatusOK, posts)
}

func GetPostByID(ctx *gin.Context) {
	id := ctx.Param("id")

	posts := services.GetPosts()
	for _, post := range posts {
		if post.ID == id {
			ctx.IndentedJSON(http.StatusOK, post)
			return
		}
	}
	ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": "post not found"})
}

func CreatePost(ctx *gin.Context) {
	var postData models.Post

	// TODO: Parse and validate post data
	if err := ctx.BindJSON(&postData); err != nil {
		return
	}
	
	newPost := services.CreatePost(postData)
		
	ctx.IndentedJSON(http.StatusCreated, newPost)
}
