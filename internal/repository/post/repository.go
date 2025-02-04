package post

import (
    "gorm.io/gorm"
)

type Repository struct {
    db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
    return &Repository{db: db}
}

func (r *Repository) CreatePost(post *Post) error {
    return r.db.Create(post).Error
}

func (r *Repository) GetPost(id string) (*Post, error) {
    var post Post
    err := r.db.First(&post, "id = ?", id).Error
    return &post, err
}

func (r *Repository) UpdatePost(post *Post) error {
    return r.db.Save(post).Error
}

func (r *Repository) DeletePost(id string) error {
    return r.db.Delete(&Post{}, "id = ?", id).Error
}

func (r *Repository) PostExists(id string) bool {
    var count int64
    r.db.Model(&Post{}).Where("id = ?", id).Count(&count)
    return count > 0
}
