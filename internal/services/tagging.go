package services

import (
	errs "cvwo-backend/internal/errors"
	"cvwo-backend/internal/repos"
)

type TaggingService struct {
	postRepo  repos.PostRepo
	topicRepo repos.TopicRepo
}

func NewTaggingService(postRepo repos.PostRepo, topicRepo repos.TopicRepo) *TaggingService {
	return &TaggingService{postRepo, topicRepo}
}

func (service *TaggingService) TagPostWithTopics(postId uint, topicIDs []uint, currentUserID uint) error {
	post, err := service.postRepo.GetByID(postId)
	if err != nil {
		return err
	}

	topics, err := service.topicRepo.GetByIDs(topicIDs)
	if err != nil {
		return err
	}

	// Check authorization
	if currentUserID != post.AuthorID {
		return errs.New(errs.ErrUnauthorized, "Unauthorized")
	}

	return service.postRepo.AssociatePostWithTopics(post, topics)
}
