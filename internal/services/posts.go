package services

import (
	"cvwo-backend/internal/models"
	"cvwo-backend/internal/repos"
	"errors"

	"gorm.io/gorm"
)

type PostService struct {
	postRepo repos.PostRepo
	userRepo repos.UserRepo
}

func NewPostService(postRepo repos.PostRepo, userRepo repos.UserRepo) *PostService {
	return &PostService{postRepo, userRepo}
}

func (service *PostService) GetAll() ([]models.Post, error) {
	return service.postRepo.GetAll()
}

func (service *PostService) GetByID(id uint) (*models.Post, error) {
	return service.postRepo.GetByID(id)
}

func (service *PostService) GetByTopic(topicID uint) ([]models.Post, error) {
	return service.postRepo.GetByTopic(topicID)
}

func (service *PostService) Create(postData *models.Post) (*models.Post, error) {
	post, err := service.postRepo.Create(postData)
	if err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return nil, NewValidationError("author_id not found")
		}
		return nil, err
	}
	return post, nil
}

// Update the title and content of the given post
func (service *PostService) Update(id uint, title string, content string) (*models.Post, error) {
	return service.postRepo.Update(id, title, content)
}

func (service *PostService) Delete(id uint) error {
	return service.postRepo.Delete(id)
}
