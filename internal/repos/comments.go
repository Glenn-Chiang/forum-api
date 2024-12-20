package repos

import (
	"cvwo-backend/internal/models"

	"gorm.io/gorm"
)

type CommentRepo struct {
	DB *gorm.DB
}

func NewCommentRepo(db *gorm.DB) *CommentRepo {
	return &CommentRepo{DB: db}
}

// Get all comments associated with the given post. Each comment includes the associated author.
func (repo *CommentRepo) GetByPostID(postId uint, limit int, offset int, sortBy string, currentUserID uint) ([]models.Comment, int64, error) {
	var comments []models.Comment

	// Apply filter
	query := repo.DB
	query = query.Where("post_id = ?", postId)

	// Get the filtered, sorted and paginated comments
	if err := query.Preload("Author").Limit(limit).Offset(offset).Order(sortBy).Find(&comments).Error; err != nil {
		return nil, 0, err
	}

	// Get the total number of comments associated with the given post
	var count int64
	if err := query.Model(&models.Comment{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}
	
	return comments, count, nil
}

// Get a particular comment including the associated author
func (repo *CommentRepo) GetByID(id uint) (*models.Comment, error) {
	var comment models.Comment
	if err := repo.DB.Preload("Author").First(&comment, id).Error; err != nil {
		return nil, err
	}
	return &comment, nil
}

// Create a new comment
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

// Delete an individual comment
func (repo *CommentRepo) Delete(id uint) error {
	result := repo.DB.Delete(&models.Comment{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
