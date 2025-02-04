package post

import (
    "time"
    "gorm.io/gorm"
    "apicpt/internal/repository/user"
)

type Post struct {
    gorm.Model
    ID        string `gorm:"primaryKey"`
    Title     string
    Content   string
    Status    string
    AuthorID  uint
    Author    user.User `gorm:"foreignKey:AuthorID"`
    CreatedAt time.Time
    UpdatedAt time.Time
}

type Input struct {
    Title   string `json:"title"`
    Content string `json:"content"`
}
