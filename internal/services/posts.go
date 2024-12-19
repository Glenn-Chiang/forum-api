package services

import (
	errs "cvwo-backend/internal/errors"
	"cvwo-backend/internal/models"
	"cvwo-backend/internal/repos"
	"errors"

	"gorm.io/gorm"
)

type PostService struct {
	postRepo  repos.PostRepo
	userRepo  repos.UserRepo
	topicRepo repos.TopicRepo
}

func NewPostService(postRepo repos.PostRepo, userRepo repos.UserRepo, topicRepo repos.TopicRepo) *PostService {
	return &PostService{postRepo, userRepo, topicRepo}
}

// Maps valid sort params to the corresponding SQL orderBy clause
var postSortFields = map[string]string{
	"new": "created_at DESC",
	"old": "created_at ASC",
}

// Get the SQL orderBy clause corresponding to the given sort param, if valid
func validPostSortField(sortBy string) (string, error) {
	sortField, exists := postSortFields[sortBy]
	if !exists {
		return "", errs.New(errs.ErrInvalid, "Invalid sort field")
	}
	return sortField, nil
}

// Get a list of posts
func (service *PostService) GetList(limit, offset int, sortBy string) ([]models.Post, int64, error) {
	// Validate sortBy param
	sortField, err := validPostSortField(sortBy)
	if err != nil {
		return nil, 0, err
	}
	return service.postRepo.GetList(limit, offset, sortField)
}

// Get all posts tagged with at least 1 of the given topics
func (service *PostService) GetByTags(topicIDs []uint, limit, offset int, sortBy string) ([]models.Post, int64, error) {
	// Validate sortBy param
	sortField, err := validPostSortField(sortBy)
	if err != nil {
		return nil, 0, err
	}

	return service.postRepo.GetByTopics(topicIDs, limit, offset, sortField)
}

// Get an individual post by ID
func (service *PostService) GetByID(id uint) (*models.Post, error) {
	post, err := service.postRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.New(errs.ErrNotFound, "Post not found")
		}
		return nil, err
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
