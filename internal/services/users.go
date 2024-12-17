package services

import (
	errs "cvwo-backend/internal/errors"
	"cvwo-backend/internal/models"
	"cvwo-backend/internal/repos"
	"errors"

	"gorm.io/gorm"
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
	user, err := service.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.New(errs.ErrNotFound, "User not found")
		}
	}
	return user, nil
}

func (service *UserService) Create(userData *models.AuthInput) (*models.User, error) {
	// Hash password
	passwordHash, err := HashPassword(userData.Password)
	if err != nil {
		return nil, err
	}

	user, err := service.repo.Create(&models.User{
		Username: userData.Username,
		Password: passwordHash, // Create user using hashed password
	})

	// Check if username is already in use
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, errs.New(errs.ErrConflict, "Username already in use")
		}
		return nil, err
	}
	return user, nil
}

func (service *UserService) Delete(id uint) error {
	return service.repo.Delete(id)
}
