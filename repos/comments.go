package repos

import (
	"cvwo-backend/models"
	"fmt"

	"gorm.io/gorm"
)

type CommentRepo struct {
	DB *gorm.DB
}

func NewCommentRepo(db *gorm.DB) *CommentRepo {
	return &CommentRepo{DB: db}
}

func (repo *CommentRepo) GetAll() ([]models.Comment, error) {
	var comments []models.Comment
	if err := repo.DB.Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}

func (repo *CommentRepo) GetByID(id uint) (*models.Comment, error) {
	var comment models.Comment
	if err := repo.DB.First(&comment, id).Error; err != nil {
		return nil, err
	}
	return &comment, nil
}

func (repo *CommentRepo) Create(comment *models.Comment) (*models.Comment, error) {
	if err := repo.DB.Create(comment).Error; err != nil {
		return nil, err
	}
	return comment, nil
}

func (repo *CommentRepo) Delete(id uint) error {
	result := repo.DB.Delete(&models.Comment{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("comment not found")
	}
	return nil
}
