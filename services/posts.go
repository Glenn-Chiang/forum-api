package services

import (
	"cvwo-backend/models"
	"cvwo-backend/repos"
)

func GetPosts() []models.Post {
	return repos.GetPosts()
}

func GetPostByID(id string) models.Post {
	return repos.GetPostById(id)
}

func CreatePost(postData models.Post) models.Post {
	return repos.CreatePost(postData)
}
