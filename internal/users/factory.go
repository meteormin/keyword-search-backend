package users

import (
	"gorm.io/gorm"
)

func Factory(db *gorm.DB) Handler {
	repo := NewRepository(db)
	service := NewService(repo)
	return NewHandler(service)
}
