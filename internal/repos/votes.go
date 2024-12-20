package repos

import (
	"cvwo-backend/internal/models"
	"errors"

	"gorm.io/gorm"
)

type VoteRepo struct {
	DB *gorm.DB
}

func NewVoteRepo(db *gorm.DB) *VoteRepo {
	return &VoteRepo{DB: db}
}

// Update existing vote or create new vote if the user has not voted for the post
func (repo *VoteRepo) Upsert(vote *models.Vote) error {
	var existingVote models.Vote
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
func (repo *VoteRepo) Delete(postID, userID uint) error {
	return repo.DB.Delete(&models.Vote{}, "post_id = ? AND user_id = ?", postID, userID).Error
}
