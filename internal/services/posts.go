package services

import (
	errs "cvwo-backend/internal/errors"
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

func (service *PostService) GetList(limit, offset int, sortBy string) ([]models.Post, error) {
	return service.postRepo.GetList(limit, offset, sortBy)
}

// Get all posts tagged with the specified topic
func (service *PostService) GetByTopic(topicID uint, limit, offset int, sortBy string) ([]models.Post, error) {
	return service.postRepo.GetByTopic(topicID, limit, offset, sortBy)
}

// Get an individual post by ID
func (service *PostService) GetByID(id uint) (*models.Post, error) {
	post, err := service.postRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.New(errs.ErrNotFound, "Post not found")
		}
	}
	return post, nil
}

// Create a new post
func (service *PostService) Create(postData *models.Post) (*models.Post, error) {
	post, err := service.postRepo.Create(postData)
	if err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return nil, errs.New(errs.ErrNotFound, "Author not found")
		}
		return nil, err
	}
	return post, nil
}

// Update the title and content of the given post
func (service *PostService) Update(id uint, title string, content string) (*models.Post, error) {
	post, err := service.postRepo.Update(id, title, content)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.New(errs.ErrNotFound, "Post not found")
		}
		return nil, err
	}
	return post, nil
}

// Delete an individual post
func (service *PostService) Delete(id uint) error {
	if err := service.postRepo.Delete(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.New(errs.ErrNotFound, "Post not found")
		}
		return err
	}
	return nil
}
