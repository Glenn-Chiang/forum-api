package models

import "time"

type User struct {
	ID       uint   `json:"id"`
	Username string `gorm:"uniqueIndex;not null" json:"username"`
	Password string `gorm:"not null" json:"-"` // Hashed password, excluded from JSON
}

// Structure of request body for login/register
type AuthInput struct {
	Username string `binding:"required,max=20"`
	Password string `binding:"required,min=5,max=20"`
}

type Post struct {
	ID        uint      `json:"id"`
	Title     string    `gorm:"not null" json:"title"`
	Content   string    `gorm:"not null" json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	AuthorID  uint      `json:"author_id"`
	// One post has one author (user). When the associated user is deleted, set the author field to null
	Author    *User     `gorm:"constraint:OnDelete:SET NULL;" json:"author,omitempty"` 
	// Implicitly create a many2many join table between posts and topics. When a post is deleted, the post_topic record in the join table is deleted. The associated topics themselves are not deleted.
	Topics    []Topic   `gorm:"many2many:post_topics;constraint:OnDelete:CASCADE;" json:"topics"` 
	// One post has many votes
	Votes []Vote `json:"votes"`
}

// Record for one user's vote on one post
type Vote struct {
	// Composite primary key using post_id and user_id
	PostID uint `gorm:"primaryKey;autoIncrement:false" json:"post_id"`
	UserID uint `gorm:"primaryKey;autoIncrement:false" json:"user_id"` 
	Value int `gorm:"not null" json:"vote"` //upvote: 1, downvote: -1
}

// Structure of request body for creating a new post
type NewPost struct {
	Title    string `binding:"required,max=200"`
	Content  string `binding:"required,min=10,max=1000"`
	TopicIDs []uint 
}

// Structure of request body for updating a post
type PostUpdate struct {
	Title   string `binding:"required,max=200"`
	Content string `binding:"required,min=10,max=1000"`
}

// Structure of request body for updating the topics associated with a post
type PostTagsUpdate struct {
	TopicIDs []uint 
}

type Comment struct {
	ID        uint      `json:"id"`
	Content   string    `gorm:"not null" json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	PostID    uint      `json:"post_id"` // By default, when the associated post is deleted, the comment remains but the post_id is set to null
	AuthorID  uint      `json:"author_id"`
	Author    User      `gorm:"constraint:OnDelete:SET NULL;" json:"author,omitempty"` // When the associated user is deleted, set the author field to null
}

// Structure of request body for creating a new comment
type NewComment struct {
	Content  string `binding:"required,max=1000"`
	PostID   uint   `binding:"required"`
}

// Structure of request body for updating a comment
type CommentUpdate struct {
	Content string `binding:"required,max=1000"`
}

type Topic struct {
	ID    uint   `json:"id"`
	Name  string `gorm:"uniqueIndex;not null" json:"name"`
}
