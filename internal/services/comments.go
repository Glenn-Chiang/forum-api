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
	postRepo    repos.PostRepo
	userRepo    repos.UserRepo
}

func NewCommentService(commentRepo repos.CommentRepo, postRepo repos.PostRepo, userRepo repos.UserRepo) *CommentService {
	return &CommentService{commentRepo, postRepo, userRepo}
}

// Maps valid sort params to the corresponding SQL orderBy clause
var commentSortFields = map[string]string{
	"new": "created_at DESC",
	"old": "created_at ASC",
}

// Get the SQL orderBy clause corresponding to the given sort param, if valid
func validCommentSortField(sortBy string) (string, error) {
	sortField, exists := commentSortFields[sortBy]
	if !exists {
		return "", errs.New(errs.ErrInvalid, "Invalid sort field")
	}
	return sortField, nil
}

// Get all comments associated with a specified post
func (service *CommentService) GetByPostID(postId uint, limit int, offset int, sortBy string) ([]models.Comment, error) {
	// Validate sortBy param
	sortField, err := validCommentSortField(sortBy)
	if err != nil {
		return nil, err
	}

	comments, err := service.commentRepo.GetByPostID(postId, limit, offset, sortField)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.New(errs.ErrNotFound, "Post not found")
		}
		return nil, err
	}
	return comments, nil
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

// Get total number of comments
func (service *CommentService) GetTotalCount() (int64, error) {
	return service.commentRepo.GetTotalCount()
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
