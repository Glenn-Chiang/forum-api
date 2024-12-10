package services

import (
	"cvwo-backend/internal/models"
	"cvwo-backend/internal/repos"
)

type TopicService struct {
	repo repos.TopicRepo
}

func NewTopicService(repo repos.TopicRepo) *TopicService {
	return &TopicService{repo}
}

func (service *TopicService) GetAll() ([]models.Topic, error) {
	return service.repo.GetAll()
}

func (service *TopicService) GetByID(id uint) (*models.Topic, error) {
	return service.repo.GetByID(id)
}

func (service *TopicService) Create(postData *models.Topic) (*models.Topic, error) {
	// TODO: Parse and validate the new post data
	return service.repo.Create(postData)
}

func (service *TopicService) Delete(id uint) error {
	return service.repo.Delete(id)
}

