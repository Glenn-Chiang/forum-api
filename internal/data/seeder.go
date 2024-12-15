package data

import (
	"log"

	"gorm.io/gorm"

	"cvwo-backend/internal/models"
)

var users = []models.User{
	{Username: "Viktor"},
	{Username: "Friedrich Nietzsche"},
	{Username: "Hamlet"},
	{Username: "Macbeth"},
}

var posts = []models.Post{
	{Title: "The Glorious Evolution", Content: "I thought I could put an end to the world's suffering. But when every equation was solved, all that remained were fields of dreamless solitude. There is no prize to perfection, only an end to pursuit.", AuthorID: 1},
	{Title: "God is Dead", Content: "God is dead. God remains dead. And we have killed him. How shall we ever comfort ourselves, the murderers of all murderers?", AuthorID: 2},
	{Title: "The Abyss", Content: "Beware that when fighting monsters, you do not yourself become a monster. For when you gaze long enough into the abyss, the abyss also gazes into you", AuthorID: 2},
	{Title: "To be or not to be", Content: "To be, or not to be - that is the question. Whether tis nobler in the mind to suffer, the slings and arrows of outrageous fortune, or to take arms against a sea of troubles, and by opposing - end them.", AuthorID: 3},
	{Title: "Out, brief candle", Content: "Out, out brief candle. Life's but a walking shadow. A poor player that struts and frets his hour upon the stage, and then is heard no more.", AuthorID: 4},
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
