package users

import (
	"github.com/miniyus/keyword-search-backend/repo"
	"gorm.io/gorm"
)

func New(db *gorm.DB) Handler {
	repository := repo.NewUserRepository(db)
	service := NewService(repository)
	return NewHandler(service)
}
