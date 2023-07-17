package hosts

import (
	"github.com/miniyus/keyword-search-backend/repo"
	"gorm.io/gorm"
)

func New(db *gorm.DB) Handler {
	repository := repo.NewHostRepository(db)
	service := NewService(repository)
	return NewHandler(service)
}
