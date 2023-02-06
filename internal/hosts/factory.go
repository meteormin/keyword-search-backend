package hosts

import (
	"gorm.io/gorm"
)

func New(db *gorm.DB) Handler {
	repo := NewRepository(db)
	service := NewService(repo)
	return NewHandler(service)
}
