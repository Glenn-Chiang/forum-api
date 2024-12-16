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

func (service *TopicService) GetByIDs(ids []uint) ([]models.Topic, error) {
	return service.repo.GetByIDs(ids)
}
