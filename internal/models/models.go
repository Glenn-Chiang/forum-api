package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex" json:"username"`
}

type Post struct {
	gorm.Model
	Title string `json:"title"`
	Content string `json:"content"`
	AuthorID uint `json:"authorId"`
	Author User `gorm:"constraint:OnDelete:SET NULL;"`
	Topics []*Topic `gorm:"many2many:post_topics"`
}

type Comment struct {
	gorm.Model
	Content string `json:"content"`
	PostID uint `json:"postId"`
	Post Post `gorm:"constraint:OnDelete:SET NULL;"`
	AuthorID uint `json:"authorId"`
	Author User `gorm:"constraint:OnDelete:SET NULL;"`
}

type Topic struct {
	ID uint `json:"id"`
	Name string `gorm:"uniqueIndex" json:"name"`
	Posts []*Post `gorm:"many2many:post_topics;"`
}
