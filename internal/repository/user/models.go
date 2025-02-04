package user

import (
    "time"
    "gorm.io/gorm"
)

type User struct {
    gorm.Model
    Email     string    `gorm:"uniqueIndex"`
    Password  string
    Role      string
    CreatedAt time.Time
    UpdatedAt time.Time
}

type Token struct {
    gorm.Model
    UserID       uint
    RefreshToken string
    ExpiresAt    time.Time
    User         User `gorm:"foreignKey:UserID"`
}

type InputAuth struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=8"`
    Role     string `json:"role" binding:"required,oneof=Author Reader"`
}

type InputLogin struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}
