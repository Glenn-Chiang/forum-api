package services

import (
	errs "cvwo-backend/internal/errors"
	"cvwo-backend/internal/models"
	"cvwo-backend/internal/repos"
)

type VotingService struct {
	repo repos.VoteRepo
}

func NewVotingService(repo repos.VoteRepo) *VotingService {
	return &VotingService{repo}
}

// Create a new vote associated to one user and one post
func (service *VotingService) Vote(postID, userID uint, value int, currentUserID uint) error {
	// Check authorization
	if currentUserID != userID {
		return errs.New(errs.ErrUnauthorized, "Unauthorized")
	}

	// If vote value is 0, delete the vote record
	if value == 0 {
		return service.repo.Delete(postID, userID)
	}

	// Only allow value of 1 (upvote) or -1 (downvote)
	if value != 1 && value != -1 {
		return errs.New(errs.ErrInvalid, "Invalid vote value")
	}

	return service.repo.Upsert(&models.Vote{PostID: postID, UserID: userID, Value: value})
}
