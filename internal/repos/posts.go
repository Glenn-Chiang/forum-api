package repos

import (
	"cvwo-backend/internal/models"

	"gorm.io/gorm"
)

type PostRepo struct {
	DB *gorm.DB
}

func NewPostRepo(db *gorm.DB) *PostRepo {
	return &PostRepo{DB: db}
}

// Get a list of all posts including their associated topics
func (repo *PostRepo) GetList(limit, offset int) ([]models.Post, error) {
	var posts []models.Post
	if err := repo.DB.Preload("Topics").Limit(limit).Offset(offset).Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

// Get all posts associated with a particular topic. Includes the associated topics of each post.
func (repo *PostRepo) GetByTopic(topicId uint, limit, offset int) ([]models.Post, error) {
	var posts []models.Post
	err := repo.DB.Preload("Topics").Joins("JOIN post_topics ON posts.id = post_topics.post_id").
		Where("post_topics.topic_id = ?", topicId).
		Find(&posts).Error
	if err != nil {
		return nil, err
	}
	return posts, nil
}

// Get all posts made by a particular user, including the associated topics of each post
func (repo *PostRepo) GetByUserID(userId uint) ([]models.Post, error) {
	var posts []models.Post
	if err := repo.DB.Preload("Topics").Find(&posts, models.Post{AuthorID: userId}).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

// Get an individual post including the associated author and topics
func (repo *PostRepo) GetByID(id uint) (*models.Post, error) {
	var post models.Post
	if err := repo.DB.Preload("Topics").Preload("Author").First(&post, id).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

// Create a new post
func (repo *PostRepo) Create(post *models.Post) (*models.Post, error) {
	if err := repo.DB.Create(post).Error; err != nil {
		return nil, err
	}
	return post, nil
}

// Update the title and content of the given post
func (repo *PostRepo) Update(id uint, title string, content string) (*models.Post, error) {
	var post models.Post
	if err := repo.DB.First(&post, id).Error; err != nil {
		return nil, err
	}

	if err := repo.DB.Model(&post).Updates(models.Post{Title: title, Content: content}).Error; err != nil {
		return nil, err
	}

	return &post, nil
}

// Replace the current list of topics associated with the given post with the given new list of topics
func (repo *PostRepo) AssociatePostWithTopics(post *models.Post, topics []models.Topic) error {
	return repo.DB.Model(post).Association("Topics").Replace(topics)
}

func (repo *PostRepo) Delete(id uint) error {
	result := repo.DB.Delete(&models.Post{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
