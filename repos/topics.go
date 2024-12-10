package repos

import (
	"cvwo-backend/models"
	"fmt"

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

func (repo *TopicRepo) GetByID(id uint) (*models.Topic, error) {
	var topic models.Topic
	if err := repo.DB.First(&topic, id).Error; err != nil {
		return nil, err
	}
	return &topic, nil
}

func (repo *TopicRepo) Create(topic *models.Topic) (*models.Topic, error) {
	if err := repo.DB.Create(topic).Error; err != nil {
		return nil, err
	}
	return topic, nil
}

func (repo *TopicRepo) Delete(id uint) error {
	result := repo.DB.Delete(&models.Topic{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("topic not found")
	}
	return nil
}
