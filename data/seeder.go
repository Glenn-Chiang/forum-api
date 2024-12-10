package data

import (
	"log"

	"gorm.io/gorm"

	"cvwo-backend/models"
)

var posts = []models.Post{
	{Title: "What Could Have Been"},
	{Title: "Goodbye"},
	{Title: "The Glorious Evolution"},
}

func SeedData(db *gorm.DB) error {
	var count int64
	if err := db.Model(&models.Post{}).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		log.Println("Database already seeded")
		return nil
	}

	if err := db.Create(&posts).Error; err != nil {
		return err
	}

	log.Println("Database seeded successfully")
	return nil
}
