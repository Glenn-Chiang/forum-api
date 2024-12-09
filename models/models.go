package models

type User struct {
	ID string `json:"id"`
	Username string `json:"username"`
}

type Post struct {
	ID string `json:"id"`
	Title string `json:"title"`
	Content string `json:"content"`
	AuthorID string `json:"authorId"`
}

type Comment struct {
	ID string `json:"id"`
	Content string `json:"content"`
	PostID string `json:"postId"`
}
