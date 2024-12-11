package data

import (
	"cvwo-backend/internal/models"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDB(dsn string) *gorm.DB {
	// Open database
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Migrate tables based on models
	if err := db.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{}, &models.Topic{}); err != nil {
		log.Fatalf("Failed to migrate tables: %v", err)
	}

	// Seed database with initial data
	if err := SeedData(db); err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}

	return db
}
