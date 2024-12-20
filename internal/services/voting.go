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

// Create a new vote associated to one user and one post, with a value of 1 indicating an upvote
func (service *VotingService) Vote(postID, userID uint, value int, currentUserID uint) error {
	// Check authorization
	if currentUserID != userID {
		return errs.New(errs.ErrUnauthorized, "Unauthorized")
	}

	return service.repo.Upsert(&models.Vote{PostID: postID, UserID: userID, Value: value})
}

// Remove the user's vote for a post
func (service *VotingService) RemoveVote(postID, userID uint, currentUserID uint) error {
	// Check authorization
	if currentUserID != userID {
		return errs.New(errs.ErrUnauthorized, "Unauthorized")
	}

	return service.repo.Delete(postID, userID)
}
