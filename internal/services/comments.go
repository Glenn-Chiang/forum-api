package services

import (
	errs "cvwo-backend/internal/errors"
	"cvwo-backend/internal/models"
	"cvwo-backend/internal/repos"
	"errors"

	"gorm.io/gorm"
)

type CommentService struct {
	commentRepo repos.CommentRepo
	postRepo repos.PostRepo
	userRepo repos.UserRepo
}

func NewCommentService(commentRepo repos.CommentRepo, postRepo repos.PostRepo, userRepo repos.UserRepo) *CommentService {
	return &CommentService{commentRepo, postRepo, userRepo}
}

func (service *CommentService) GetAll() ([]models.Comment, error) {
	return service.commentRepo.GetAll()
}

// Get an individual comment by ID
func (service *CommentService) GetByID(id uint) (*models.Comment, error) {
	comment, err := service.commentRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.New(errs.ErrNotFound, "Comment not found")
		}
		return nil, err
	}
	return comment, nil
}

// Get all comments associated with a specified post
func (service *CommentService) GetByPostID(id uint) ([]models.Comment, error) {
	comments, err := service.commentRepo.GetByPostID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.New(errs.ErrNotFound, "Post not found")
		}
		return nil, err
	}
	return comments, nil
}

// Create a new comment associated with a specific post and user
func (service *CommentService) Create(commentData *models.Comment) (*models.Comment, error) {
	comment, err := service.commentRepo.Create(commentData)
	if err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return nil, errs.New(errs.ErrNotFound, "Post or author not found")
		}
		return nil, err
	}
	return comment, nil
}

// Update the content of the given comment
func (service *CommentService) Update(id uint, content string) (*models.Comment, error) {
	comment, err := service.commentRepo.Update(id, content)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.New(errs.ErrNotFound, "Comment not found")
		}
		return nil, err
	}
	return comment, nil
}

// Delete an individual comment
func (service *CommentService) Delete(id uint) error {
	if err := service.commentRepo.Delete(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.New(errs.ErrNotFound, "Comment not found")
		}
		return err
	}
	return nil
}

