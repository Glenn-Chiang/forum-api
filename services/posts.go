package services

import (
	"cvwo-backend/models"
	"cvwo-backend/repos"
)

type PostService struct {
	repo repos.PostRepo
}

func NewPostService(repo repos.PostRepo) *PostService {
	return &PostService{repo}
}

func (service *PostService) GetAll() ([]models.Post, error) {
	return service.repo.GetAll()
}

func (service *PostService) GetByID(id string) (*models.Post, error) {
	return service.repo.GetByID(id)
}

func (service *PostService) Create(postData *models.Post) (*models.Post, error) {
	// TODO: Parse and validate the new post data
	return service.repo.Create(postData)
}
