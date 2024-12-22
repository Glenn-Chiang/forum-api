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

// Get all comments associated with the given post
func (repo *CommentRepo) GetByPostID(postID uint, limit int, offset int, sortBy string, currentUserID uint) ([]models.Comment, int64, error) {
	var comments []models.Comment

	// Apply filter
	filteredDB := repo.DB.Where("post_id = ?", postID)

	err := repo.DB.Model(&models.Comment{}).
		Preload("Author"). // Include comment author
		Select("comments.*, "+
			// Compute net votes of the comment
			"COALESCE(SUM(votes.value),0) AS net_votes, "+
			// Get the current user's vote for the comment
			"COALESCE(user_votes.value,0) AS user_vote").
		// Filter comments corresponding to the post
		Where("comments.post_id = ?", postID). 
		// Get all vote records associated to the comment
		Joins("LEFT JOIN comment_votes AS votes ON comments.id = votes.comment_id"). 
		// Get the single vote record made by the current user, that is associated to the comment
		Joins("LEFT JOIN comment_votes AS user_votes ON comments.id = user_votes.comment_id AND user_votes.user_id = ?", currentUserID). 
		Group("comments.id").
		Limit(limit).Offset(offset).Order(sortBy).
		Find(&comments).Error

	if err != nil {
		return nil, 0, err
	}
		
	// Get the total number of comments associated with the given post
	var count int64
	if err := filteredDB.Model(&models.Comment{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}
	
	return comments, count, nil
}

// Get an individual comment
func (repo *CommentRepo) GetByID(id uint) (*models.Comment, error) {
	var comment models.Comment
	if err := repo.DB.First(&comment, id).Error; err != nil {
		return nil, err
	}
	return &comment, nil
}

// Similar to GetByID but includes additional computed fields and preloaded associations
// Takes in currentUserID in order to compute user_vote field
func (repo *CommentRepo) GetByIDWithAuth(commentID uint, currentUserID uint) (*models.Comment, error) {
	var comment models.Comment

	err := repo.DB.Model(&models.Comment{}).
		Preload("Author").
		Select("comments.*, " +
			// Compute net votes for the comment
			"COALESCE(SUM(votes.value),0) AS net_votes, " +
			// Get the current user's vote for the comment
			"COALESCE(user_votes.value, 0) AS user_vote").
		Joins("LEFT JOIN comment_votes AS votes ON comments.id = votes.comment_id").
		Joins("LEFT JOIN comment_votes AS user_votes ON comments.id = user_votes.comment_id AND user_votes.user_id = ?", currentUserID).
		Where("comments.id = ?", commentID).
		Group("comments.id").
		Find(&comment).Error
	
	if err != nil {
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
