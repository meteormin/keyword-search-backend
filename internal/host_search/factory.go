package host_search

import (
	worker "github.com/miniyus/goworker"
	"github.com/miniyus/keyword-search-backend/internal/search"
	"gorm.io/gorm"
)

func New(db *gorm.DB, dispatcher worker.Dispatcher) Handler {
	repo := search.NewRepository(db)
	service := search.NewService(repo)
	return NewHandler(service, dispatcher)
}
