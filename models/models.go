package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID string `json:"id"`
	Username string `json:"username"`
}

type Post struct {
	gorm.Model
	ID string `json:"id"`
	Title string `json:"title"`
	Content string `json:"content"`
	AuthorID string `json:"authorId"`
}

type Comment struct {
	gorm.Model
	ID string `json:"id"`
	Content string `json:"content"`
	PostID string `json:"postId"`
}
