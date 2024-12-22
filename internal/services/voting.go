package services

import (
	errs "cvwo-backend/internal/errors"
	"cvwo-backend/internal/models"
	"cvwo-backend/internal/repos"
)

type VotingService struct {
	postVoteRepo    repos.PostVoteRepo
	commentVoteRepo repos.CommentVoteRepo
}

func NewVotingService(postRepo repos.PostVoteRepo, commentRepo repos.CommentVoteRepo) *VotingService {
	return &VotingService{postRepo, commentRepo}
}

// Update a user's vote for a post
func (service *VotingService) VotePost(postID, userID uint, value int, currentUserID uint) error {
	// Check authorization
	if currentUserID != userID {
		return errs.New(errs.ErrUnauthorized, "Unauthorized")
	}

	// If vote value is 0, delete the vote record
	if value == 0 {
		return service.postVoteRepo.Delete(postID, userID)
	}

	// Only allow value of 1 (upvote) or -1 (downvote)
	if value != 1 && value != -1 {
		return errs.New(errs.ErrInvalid, "Invalid vote value")
	}

	return service.postVoteRepo.Upsert(&models.PostVote{PostID: postID, UserID: userID, Value: value})
}

// Update a user's vote for a comment
func (service *VotingService) VoteComment(commentID, userID uint, value int, currentUserID uint) error {
	// Check authorization
	if currentUserID != userID {
		return errs.New(errs.ErrUnauthorized, "Unauthorized")
	}

	// If vote value is 0, delete the vote record
	if value == 0 {
		return service.commentVoteRepo.Delete(commentID, userID)
	}

	// Only allow value of 1 (upvote) or -1 (downvote)
	if value != 1 && value != -1 {
		return errs.New(errs.ErrInvalid, "Invalid vote value")
	}

	return service.commentVoteRepo.Upsert(&models.CommentVote{CommentID: commentID, UserID: userID, Value: value})
}
