package services

import (
	"cvwo-backend/internal/repos"
)

type TaggingService struct {
	postRepo  repos.PostRepo
	topicRepo repos.TopicRepo
}

func NewTaggingService(postRepo repos.PostRepo, topicRepo repos.TopicRepo) *TaggingService {
	return &TaggingService{postRepo, topicRepo}
}

func (service *TaggingService) TagPostWithTopics(postId uint, topicIDs []uint) error {
	post, err := service.postRepo.GetByID(postId)
	if err != nil {
		return err
	}

	topics, err := service.topicRepo.GetByIDs(topicIDs)
	if err != nil {
		return err
	}

	return service.postRepo.AssociatePostWithTopics(post, topics)
}
