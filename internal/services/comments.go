package services

import (
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

func (service *CommentService) GetByID(id uint) (*models.Comment, error) {
	return service.commentRepo.GetByID(id)
}

func (service *CommentService) GetByPostID(id uint) ([]models.Comment, error) {
	return service.commentRepo.GetByPostID(id)
}

// Create a new comment associated with a specific post and user
func (service *CommentService) Create(commentData *models.Comment) (*models.Comment, error) {
	comment, err := service.commentRepo.Create(commentData)
	if err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return nil, NewValidationError("post_id or author_id not found")
		}
	}
	return service.commentRepo.Create(comment)
}

// Update the content of the given comment
func (service *CommentService) Update(id uint, content string) (*models.Comment, error) {
	return service.commentRepo.Update(id, content)
}

func (service *CommentService) Delete(id uint) error {
	return service.commentRepo.Delete(id)
}

