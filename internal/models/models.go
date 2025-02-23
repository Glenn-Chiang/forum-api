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
	// One post has one author (user).
	// If the associated user is deleted, set the author field to null
	Author *User `json:"author" gorm:"constraint:OnDelete:SET NULL;"`
	// Implicitly create a many2many join table between posts and topics. When a post is deleted, the post_topic record in the join table is deleted. The associated topics themselves are not deleted.
	Topics []Topic `json:"topics" gorm:"many2many:post_topics;constraint:OnDelete:CASCADE;"`
	// Array of votes associated with this post. Not included in json.
	Votes []PostVote `json:"-" gorm:"constraint:OnDelete:CASCADE;"` // When this post is deleted, the associated comments are deleted.

	// Upvotes minus downvotes
	// Computed field, not included in database
	NetVotes int `json:"votes" gorm:"->;-:migration"`

	// Indicates whether the current user has upvoted (1), downvoted (-1) or not voted (0) the post
	// Computed field, not included in database
	UserVote int `json:"user_vote" gorm:"->;-:migration"`
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
	Content   string    `json:"content" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	PostID    uint      `json:"post_id" gorm:"constraint:OnDelete:SET NULL;" ` // When the associated post is deleted, the comment remains but the post_id is set to null
	AuthorID  uint      `json:"author_id"`
	Author    User      `json:"author" gorm:"constraint:OnDelete:SET NULL;"` // When the associated user is deleted, set the author field to null

	// Array of votes associated with this comment. Not included in json.
	Votes []CommentVote `json:"-" gorm:"constraint:OnDelete:CASCADE;"` // When this comment is deleted, the associated votes are deleted.

	// Upvotes minus downvotes
	// Computed field, not included in database
	NetVotes int `json:"votes" gorm:"->;-:migration"`

	// Indicates whether the current user has upvoted (1), downvoted (-1) or not voted (0) the comment
	// Computed field, not included in database
	UserVote int `json:"user_vote" gorm:"->;-:migration"`
}

// Request body for creating a new comment
type NewComment struct {
	Content string `json:"content" binding:"required,max=1000"`
	PostID  uint   `json:"post_id" binding:"required"`
}

// Request body for updating a comment
type CommentUpdate struct {
	Content string `json:"content" binding:"required,max=1000"`
}

// Record for a user's vote for a post
type PostVote struct {
	// Composite primary key using post_id and user_id
	PostID uint `json:"post_id" gorm:"primaryKey;autoIncrement:false;"`
	UserID uint `json:"user_id" gorm:"primaryKey;autoIncrement:false;"`
	Value  int  `json:"value" gorm:"not null"` //upvote: 1, downvote: -1
}

// Record for a user's vote for a comment
type CommentVote struct {
	// Composite primary key using comment_id and user_id
	CommentID uint `json:"comment_id" gorm:"primaryKey;autoIncrement:false;" `
	UserID    uint `json:"user_id" gorm:"primaryKey;autoIncrement:false;"`
	Value     int  `json:"value" gorm:"not null"` //upvote: 1, downvote: -1
}

// Request body for voting on a post or comment
type VoteUpdate struct {
	Value int `json:"value"`
}

type Topic struct {
	ID   uint   `json:"id"`
	Name string `gorm:"uniqueIndex;not null" json:"name"`
}
