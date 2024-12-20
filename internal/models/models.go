package models

import "time"

type User struct {
	ID       uint   `json:"id"`
	Username string `gorm:"uniqueIndex;not null" json:"username"`
	Password string `gorm:"not null" json:"-"` // Hashed password, excluded from JSON
}

// Request body for login/register
type AuthInput struct {
	Username string `binding:"required,max=20"`
	Password string `binding:"required,min=5,max=20"`
}

type Post struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title" gorm:"not null" `
	Content   string    `json:"content" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	AuthorID  uint      `json:"author_id"`
	// One post has one author (user). When the associated user is deleted, set the author field to null
	Author    *User     `json:"author" gorm:"constraint:OnDelete:SET NULL;"` 
	// Implicitly create a many2many join table between posts and topics. When a post is deleted, the post_topic record in the join table is deleted. The associated topics themselves are not deleted.
	Topics    []Topic   `json:"topics" gorm:"many2many:post_topics;constraint:OnDelete:CASCADE;"` 
	// Array of votes associated with this post. Not included in json.
	Votes []Vote `json:"-"`
	// Upvotes - downvotes. net_votes is a computed field that is not included in the database schema
	NetVotes int64 `json:"votes" gorm:"-"`
}

// Record for one user's vote on one post
type Vote struct {
	// Composite primary key using post_id and user_id
	PostID uint `json:"post_id" gorm:"primaryKey;autoIncrement:false" `
	UserID uint `json:"user_id" gorm:"primaryKey;autoIncrement:false"` 
	Value int `json:"value" gorm:"not null"` //upvote: 1, downvote: -1
}

// Request body for voting on a post
type PostVote struct {
	Value int `json:"value"`
}

// Request body for creating a new post
type NewPost struct {
	Title    string `json:"title" binding:"required,max=200"`
	Content  string `json:"content" binding:"required,min=10,max=1000"`
	TopicIDs []uint `json:"topic_ids"`
}

// Request body for updating a post
type PostUpdate struct {
	Title   string `json:"title" binding:"required,max=200"`
	Content string `json:"content" binding:"required,min=10,max=1000"`
}

// Request body for updating the topics associated with a post
type PostTagsUpdate struct {
	TopicIDs []uint `json:"topic_ids"`
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

// Request body for creating a new comment
type NewComment struct {
	Content  string `json:"content" binding:"required,max=1000"`
	PostID   uint   `json:"post_id" binding:"required"`
}

// Request body for updating a comment
type CommentUpdate struct {
	Content string `json:"content" binding:"required,max=1000"`
}

type Topic struct {
	ID    uint   `json:"id"`
	Name  string `gorm:"uniqueIndex;not null" json:"name"`
}
