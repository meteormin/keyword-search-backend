package short_url

import (
	"github.com/miniyus/go-fiber/internal/app/search"
	"gorm.io/gorm"
)

func New(db *gorm.DB) Handler {
	repo := search.NewRepository(db)
	service := NewService(repo)
	return NewHandler(service)
}
