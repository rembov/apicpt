package services

import "apicpt/internal/entites"

type postRepo interface {
	GetPost(id string) (*models.Post, error)
}
