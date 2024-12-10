package services

import (
	"cvwo-backend/models"
	"cvwo-backend/repos"
)

type CommentService struct {
	repo repos.CommentRepo
}

func NewCommentService(repo repos.CommentRepo) *CommentService {
	return &CommentService{repo}
}

func (service *CommentService) GetAll() ([]models.Comment, error) {
	return service.repo.GetAll()
}

func (service *CommentService) GetByID(id uint) (*models.Comment, error) {
	return service.repo.GetByID(id)
}

func (service *CommentService) GetByPostID(id uint) ([]models.Comment, error) {
	return service.repo.GetByPostID(id)
}

func (service *CommentService) Create(postData *models.Comment) (*models.Comment, error) {
	// TODO: Parse and validate the new post data
	return service.repo.Create(postData)
}

func (service *CommentService) Delete(id uint) error {
	return service.repo.Delete(id)
}

