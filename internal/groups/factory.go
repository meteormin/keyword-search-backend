package groups

import (
	"github.com/miniyus/keyword-search-backend/repo"
	"gorm.io/gorm"
)

func New(db *gorm.DB) Handler {
	repository := repo.NewGroupRepository(db)
	service := NewService(repository)
	return NewHandler(service)
}
