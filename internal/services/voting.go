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
func (service *VotingService) Vote(postID, userID uint, value bool) error {
	return service.repo.Upsert(&models.Vote{PostID: postID, UserID: userID, Value: value})
}

// Remove the user's vote for a post
func (service *VotingService) RemoveVote(postID, userID uint) error {
	return service.repo.Delete(postID, userID)
}
