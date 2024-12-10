package services

import (
	"cvwo-backend/models"
	"cvwo-backend/repos"
)

type UserService struct {
	repo repos.UserRepo
}

func NewUserService(repo repos.UserRepo) *UserService {
	return &UserService{repo}
}

func (service *UserService) GetAll() ([]models.User, error) {
	return service.repo.GetAll()
}

func (service *UserService) GetByID(id uint) (*models.User, error) {
	return service.repo.GetByID(id)
}

func (service *UserService) Create(postData *models.User) (*models.User, error) {
	// TODO: Parse and validate the new post data
	return service.repo.Create(postData)
}

func (service *UserService) Delete(id uint) error {
	return service.repo.Delete(id)
}

