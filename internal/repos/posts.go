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

// Helper function that can be used by all repository functions that involve getting a list of posts
// Helps to calculate computed fields and preload associations
func buildPostsQuery(db *gorm.DB, limit, offset int, sortBy string, currentUserID uint) *gorm.DB {
	return db.Model(&models.Post{}).
		Preload("Author").Preload("Topics"). // Include these fields in returned posts
		Select("posts.*, "+
			"SUM(votes.value) AS net_votes, "+ // Calculate net votes of the post
			"user_votes.value AS user_vote").  // Get the current user's vote for the post
		Joins("LEFT JOIN votes ON posts.id = votes.post_id").                                                              // Get all vote records associated to the post
		Joins("LEFT JOIN votes AS user_votes ON posts.id = user_votes.post_id AND user_votes.user_id = ?", currentUserID). // Get the single vote record made by the current user, that is associated to the post
		Group("posts.id").
		Limit(limit).            // Apply pagination
		Offset(offset).          // Apply pagination
		Order(sortBy).           // Apply sorting
		Session(&gorm.Session{}) // Prevent query contamination
}

// Get a list of all posts including their associated topics
// Also returns the total number of posts
func (repo *PostRepo) GetList(limit, offset int, sortBy string, currentUserID uint) ([]models.Post, int64, error) {
	var posts []models.Post
	if err := buildPostsQuery(repo.DB, limit, offset, sortBy, currentUserID).Find(&posts).Error; err != nil {
		return nil, 0, err
	}

	// Get the total number of posts
	var count int64
	if err := repo.DB.Model(&models.Post{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	return posts, count, nil
}

// Get all posts associated with at least 1 of the topics in the given list of topics
// Also returns the total number of posts filtered
func (repo *PostRepo) GetByTopics(topicIDs []uint, limit, offset int, sortBy string, currentUserID uint) ([]models.Post, int64, error) {
	var posts []models.Post

	// Filter out the posts associated with the given topics
	filteredDB := repo.DB.Joins("JOIN post_topics ON posts.id = post_topics.post_id").
		Where("post_topics.topic_id IN ?", topicIDs)

	if err := buildPostsQuery(filteredDB, limit, offset, sortBy, currentUserID).Find(&posts).Error; err != nil {
		return nil, 0, err
	}

	// Get the total count of filtered posts
	var count int64
	if err := filteredDB.Model(&models.Post{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	return posts, count, nil
}

// Get an individual post including the associated author and topics
func (repo *PostRepo) GetByID(id uint) (*models.Post, error) {
	var post models.Post

	err := repo.DB.Model(&models.Post{}).
		Preload("Topics").Preload("Author"). // Include these fields in the returned post
		Select("posts.*, SUM(votes.value) AS net_votes"). // Compute net votes
		Joins("LEFT JOIN votes ON votes.post_id = posts.id").
		Where("id = ?", id).
		Group("posts.id").
		Find(&post).Error

	if err != nil {
		return nil, err
	}

	return &post, nil
}

// Authenticated version of GetByID function, which also takes in currentUserID in order to compute user_votes field
func (repo *PostRepo) GetByIDWithAuth(postID uint, currentUserID uint) (*models.Post, error) {
	var post models.Post

	err := repo.DB.Model(&models.Post{}).
		Preload("Topics").Preload("Author"). // Include these fields in the returned post
		Select("posts.*, "+
		"COALESCE(SUM(votes.value),0) AS net_votes, "+ // Compute net votes of the post
		"COALESCE(user_votes.value,0) AS user_vote"). // Get the current user's vote for the post
		Joins("LEFT JOIN votes ON posts.id = votes.post_id").                                                              // Get all vote records associated to the post
		Joins("LEFT JOIN votes AS user_votes ON posts.id = user_votes.post_id AND user_votes.user_id = ?", currentUserID). // Get the single vote record made by the current user, that is associated to the post
		Where("id = ?", postID).
		Group("posts.id").
		Find(&post).Error

	if err != nil {
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
