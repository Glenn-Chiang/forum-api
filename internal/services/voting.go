package services

import (
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
func (service *VotingService) Upvote(postID, userID uint) error {
	return service.repo.Upsert(&models.Vote{PostID: postID, UserID: userID, Value: 1})
}

// Create a new vote associated to one user and one post, with a value of -1 indicating a downvote
func (service *VotingService) Downvote(postID, userID uint) error {
	return service.repo.Upsert(&models.Vote{PostID: postID, UserID: userID, Value: -1})
}

// Remove the user's vote for a post
func (service *VotingService) RemoveVote(postID, userID uint) error {
	return service.repo.Delete(postID, userID)
}
