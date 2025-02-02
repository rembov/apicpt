package models

import (
	"gorm.io/gorm"
	"time"
)

// Модель поста
type Post struct {
	gorm.Model
	ID        string `gorm:"primaryKey"`
	Title     string
	Content   string
	Status    string
	AuthorID  uint
	Author    User `gorm:"foreignKey:AuthorID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
