package repos

import (
	"cvwo-backend/models"
	"fmt"

	"gorm.io/gorm"
)

type PostRepo struct {
	DB *gorm.DB
}

func NewPostRepo(db *gorm.DB) *PostRepo {
	return &PostRepo{DB: db}
}

func (repo *PostRepo) GetAll() ([]models.Post, error) {
	var posts []models.Post
	if err := repo.DB.Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

func (repo *PostRepo) GetByID(id uint) (*models.Post, error) {
	var post models.Post
	if err := repo.DB.First(&post, id).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

// Get all posts made by the given user
func (repo *PostRepo) GetByUserID(userId uint) ([]models.Post, error) {
	var posts []models.Post
	if err := repo.DB.Find(&posts, models.Post{AuthorID: userId}).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

func (repo *PostRepo) Create(post *models.Post) (*models.Post, error) {
	if err := repo.DB.Create(post).Error; err != nil {
		return nil, err
	}
	return post, nil
}

// Update only the title of the given post
func (repo *PostRepo) UpdateTitle(id uint, title string) (*models.Post, error) {
	var post models.Post
	if err := repo.DB.First(&post, id).Error; err != nil {
		return nil, err
	}

	if err := repo.DB.Model(&post).Updates(models.Post{Title: title}).Error; err != nil {
		return nil, err
	}

	return &post, nil
}

// Update only the content of the given post
func (repo *PostRepo) UpdateContent(id uint, content string) (*models.Post, error) {
	var post models.Post
	if err := repo.DB.First(&post, id).Error; err != nil {
		return nil, err
	}

	if err := repo.DB.Model(&post).Updates(models.Post{Content: content}).Error; err != nil {
		return nil, err
	}

	return &post, nil
}

func (repo *PostRepo) AssociatePostWithTopics(post *models.Post, topics []models.Topic) error {
	return repo.DB.Model(post).Association("Topics").Append(topics)
}

func (repo *PostRepo) Delete(id uint) error {
	result := repo.DB.Delete(&models.User{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("post not found")
	}
	return nil
}
