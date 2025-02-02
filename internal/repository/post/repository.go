package post

import (
	"apicpt/internal/entites"
	"gorm.io/gorm"
)

type PostRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (p PostRepository) GetPost(id string) (*models.Post, error) {
	var post modePost
	err := db.DB.First(&post, "id = ?", id).Error
	return &post, err
}
