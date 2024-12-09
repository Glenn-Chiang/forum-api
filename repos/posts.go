package repos

import (
	"cvwo-backend/models"
)

var posts = []models.Post{
	{ID: "1", Title: "What Could Have Been"},
	{ID: "2", Title: "Goodbye"},
	{ID: "3", Title: "The Glorious Evolution"},
}

func GetPosts() []models.Post {
	return posts
}

func GetPostById(id string) models.Post {
	// TODO: Find post by id
	return posts[0]
}

func CreatePost(postData models.Post) models.Post {
	var newPost = models.Post{ID: "4", Title: "Oh the misery"}
	posts = append(posts, newPost)
	return newPost
}
