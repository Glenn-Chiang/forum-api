package services

import (
	"cvwo-backend/internal/models"
	"cvwo-backend/internal/repos"
	"fmt"
)

type CommentService struct {
	commentRepo repos.CommentRepo
	userRepo repos.UserRepo
}

func NewCommentService(commentRepo repos.CommentRepo, userRepo repos.UserRepo) *CommentService {
	return &CommentService{commentRepo, userRepo}
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

func (service *CommentService) Create(commentData *models.Comment) (*models.Comment, error) {
	// Check if authorID corresponds to an existing user
	if _, err := service.userRepo.GetByID(commentData.AuthorID); err != nil {
		return nil, fmt.Errorf("no author with ID %d", commentData.AuthorID)
	}
	return service.commentRepo.Create(commentData)
}

func (service *CommentService) Delete(id uint) error {
	return service.commentRepo.Delete(id)
}

