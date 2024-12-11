package data

import (
	"log"

	"gorm.io/gorm"

	"cvwo-backend/internal/models"
)

var users = []models.User{
	{Username: "Viktor"},
	{Username: "Vi"},
	{Username: "Jinx"},
}

var posts = []models.Post{
	{Title: "The Glorious Evolution", AuthorID: 1},
	{Title: "Goodbye", AuthorID: 2},
	{Title: "What Could Have Been", AuthorID: 3},
}

func SeedData(db *gorm.DB) error {
	// If there is at least 1 user, we assume the database is already populated and do not seed it
	var count int64
	if err := db.Model(&models.User{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		log.Println("Database already seeded")
		return nil
	}

	if err := db.Create(&users).Error; err != nil {
		return err
	}
	if err := db.Create(&posts).Error; err != nil {
		return err
	}

	log.Println("Database seeded successfully")
	return nil
}
