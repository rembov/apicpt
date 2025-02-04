package user

import (
    "gorm.io/gorm"
)

type Repository struct {
    db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
    return &Repository{db: db}
}

func (r *Repository) CreateUser(user *User) error {
    return r.db.Create(user).Error
}

func (r *Repository) GetUser(email string) (*User, error) {
    var user User
    err := r.db.First(&user, "email = ?", email).Error
    return &user, err
}

func (r *Repository) UpdateUser(user *User) error {
    return r.db.Save(user).Error
}

func (r *Repository) EmailExists(email string) bool {
    var count int64
    r.db.Model(&User{}).Where("email = ?", email).Count(&count)
    return count > 0
}

// Методы для работы с токенами
func (r *Repository) SaveToken(token *Token) error {
    return r.db.Create(token).Error
}

func (r *Repository) GetToken(refreshToken string) (*Token, error) {
    var token Token
    err := r.db.First(&token, "refresh_token = ?", refreshToken).Error
    return &token, err
}

func (r *Repository) DeleteToken(refreshToken string) error {
    return r.db.Delete(&Token{}, "refresh_token = ?", refreshToken).Error
}
