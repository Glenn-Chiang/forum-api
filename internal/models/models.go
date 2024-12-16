package models

import "time"

type User struct {
	ID       uint   `json:"id"`
	Username string `gorm:"uniqueIndex" json:"username"`
	Password string `json:"-"` // Hashed password, excluded from JSON
}

type AuthInput struct {
	Username string `gorm:"uniqueIndex" json:"username" binding:"required,min=5,max=20"`
	Password string `binding:"required"`
}

type Post struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	AuthorID  uint      `json:"author_id"`
	Author    *User     `gorm:"constraint:OnDelete:SET NULL;" json:"author,omitempty"`
	Topics    []Topic  `gorm:"many2many:post_topics;constraint:OnDelete:CASCADE;" json:"topics"`
}

// Structure of request body for creating a new post
type NewPost struct {
	Title    string `json:"title" binding:"required,max=200"`
	Content  string `json:"content" binding:"required,min=10,max=1000"`
	AuthorID uint   `json:"author_id" binding:"required"`
	TopicIDs []uint `json:"topic_ids"`
}

// Structure of request body for updating a post
type PostUpdate struct {
	Title   string `json:"title" binding:"required,max=200"`
	Content string `json:"content" binding:"required,min=10,max=1000"`
}

// Structure of request body for updating the topics associated with a post
type TagsUpdate struct {
	TopicIDs []uint `json:"topic_ids"`
}

type Comment struct {
	ID        uint      `json:"id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	PostID    uint      `json:"post_id"`
	AuthorID  uint      `json:"author_id"`
	Author    User      `gorm:"constraint:OnDelete:SET NULL;" json:"author,omitempty"`
}

// Structure of request body for creating a new comment
type NewComment struct {
	Content  string `json:"content" binding:"required,max=1000"`
	PostID   uint   `json:"post_id" binding:"required"`
	AuthorID uint   `json:"author_id" binding:"required"`
}

// Structure of request body for updating a comment
type CommentUpdate struct {
	Content string `json:"content" binding:"required,max=1000"`
}

type Topic struct {
	ID    uint    `json:"id"`
	Name  string  `gorm:"uniqueIndex" json:"name" binding:"required"`
	Posts []Post `gorm:"many2many:post_topics;constraint:OnDelete:CASCADE" json:"-"`
}
