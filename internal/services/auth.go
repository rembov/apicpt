package services

import (
	"gorm.io/gorm"
	"time"
)

type AuthService struct {
	db         *gorm.DB
	tokenTTL   time.Duration
	refreshTTL time.Duration
}

func NewAuthService(db *gorm.DB, tokenTTL, refreshTTL time.Duration) *AuthService {
	return &AuthService{
		db:         db,
		tokenTTL:   tokenTTL,
		refreshTTL: refreshTTL,
	}
}
