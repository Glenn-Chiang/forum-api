package repos

import (
	"cvwo-backend/internal/models"
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

// Get all comments associated with given post
func (repo *CommentRepo) GetByPostID(postId uint) ([]models.Comment, error) {
	var comments []models.Comment
	if err := repo.DB.Find(&comments, models.Comment{PostID: postId}).Error; err != nil {
		return nil, err
	}
	return comments, nil
}

// Get all comments made by given user
func (repo *CommentRepo) GetByUserID(userId uint) ([]models.Comment, error) {
	var comments []models.Comment
	if err := repo.DB.Find(&comments, models.Comment{AuthorID: userId}).Error; err != nil {
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

// Update the content of the given comment
func (repo *CommentRepo) Update(id uint, content string) (*models.Comment, error) {
	var comment models.Comment
	if err := repo.DB.First(&comment, id).Error; err != nil {
		return nil, err
	}

	if err := repo.DB.Model(&comment).Updates(models.Comment{Content: content}).Error; err != nil {
		return nil, err
	}

	return &comment, nil
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
