package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `json:"username"`
}

type Post struct {
	gorm.Model
	Title string `json:"title"`
	Content string `json:"content"`
	AuthorID string `json:"authorId"`
}

type Comment struct {
	gorm.Model
	Content string `json:"content"`
	PostID string `json:"postId"`
}
