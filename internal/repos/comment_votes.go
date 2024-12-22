package repos

import (
	"cvwo-backend/internal/models"
	"errors"

	"gorm.io/gorm"
)

type CommentVoteRepo struct {
	DB *gorm.DB
}

func NewCommentVoteRepo(db *gorm.DB) *CommentVoteRepo {
	return &CommentVoteRepo{DB: db}
}

// Update existing vote or create new vote if the user has not voted for the comment
func (repo *CommentVoteRepo) Upsert(vote *models.CommentVote) error {
	var existingVote models.CommentVote
	if err := repo.DB.First(&existingVote, "comment_id = ? AND user_id = ?", vote.CommentID, vote.UserID).Error; err != nil {
		// If this user has not voted for this comment, create new vote
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return repo.DB.Create(&vote).Error
		}
		return err
	}

	// Update existing vote
	existingVote.Value = vote.Value
	return repo.DB.Save(&existingVote).Error
}

// Delete a vote, i.e. user removes their vote for a comment
func (repo *CommentVoteRepo) Delete(commentID, userID uint) error {
	return repo.DB.Delete(&models.CommentVote{}, "comment_id = ? AND user_id = ?", commentID, userID).Error
}
