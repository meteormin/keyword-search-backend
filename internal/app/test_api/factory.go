package test_api

import (
	"gorm.io/gorm"
)

func Factory(db *gorm.DB) Handler {
	repo := NewRepository(db)
	service := NewService(repo)
	return NewHandler(service)
}
