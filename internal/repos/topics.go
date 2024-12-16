package repos

import (
	"cvwo-backend/internal/models"

	"gorm.io/gorm"
)

type TopicRepo struct {
	DB *gorm.DB
}

func NewTopicRepo(db *gorm.DB) *TopicRepo {
	return &TopicRepo{DB: db}
}

func (repo *TopicRepo) GetAll() ([]models.Topic, error) {
	var topics []models.Topic
	if err := repo.DB.Find(&topics).Error; err != nil {
		return nil, err
	}
	return topics, nil
}

// Get the list of topics with the given IDs
func (repo *TopicRepo) GetByIDs(ids []uint) ([]models.Topic, error){
	var topics []models.Topic
	if err := repo.DB.Where("id IN ?", ids).Find(&topics).Error; err != nil {
		return nil, err
	}
	return topics, nil
}
