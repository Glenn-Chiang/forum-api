package services

import (
	"cvwo-backend/internal/models"
	"cvwo-backend/internal/repos"
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
	// Check if authorID corresponds to an existing user
	if _, err := service.userRepo.GetByID(postData.AuthorID); err != nil {
		return nil, NewValidationError("author_id", "not found")
	}
	return service.postRepo.Create(postData)
}

// Update the title and content of the given post
func (service *PostService) Update(id uint, title string, content string) (*models.Post, error) {
	return service.postRepo.Update(id, title, content)
}

func (service *PostService) Delete(id uint) error {
	return service.postRepo.Delete(id)
}
