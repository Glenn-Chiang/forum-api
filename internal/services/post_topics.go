package services

import (
"cvwo-backend/internal/repos"
)

type PostTopicService struct {
	postRepo  repos.PostRepo
	topicRepo repos.TopicRepo
}

func (service *PostTopicService) TagPostWithTopics(postId uint, topicIDs []uint) error {
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
