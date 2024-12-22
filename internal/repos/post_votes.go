package repos

import (
	"cvwo-backend/internal/models"
	"errors"

	"gorm.io/gorm"
)

type PostVoteRepo struct {
	DB *gorm.DB
}

func NewPostVoteRepo(db *gorm.DB) *PostVoteRepo {
	return &PostVoteRepo{DB: db}
}

// Update existing vote or create new vote if the user has not voted for the post
func (repo *PostVoteRepo) Upsert(vote *models.PostVote) error {
	var existingVote models.PostVote
	if err := repo.DB.First(&existingVote, "post_id = ? AND user_id = ?", vote.PostID, vote.UserID).Error; err != nil {
		// If this user has not voted for this post, create new vote
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return repo.DB.Create(&vote).Error
		}
		return err
	}

	// Update existing vote
	existingVote.Value = vote.Value
	return repo.DB.Save(&existingVote).Error
}

// Delete a vote, i.e. user removes their vote for a post
func (repo *PostVoteRepo) Delete(postID, userID uint) error {
	return repo.DB.Delete(&models.PostVote{}, "post_id = ? AND user_id = ?", postID, userID).Error
}
