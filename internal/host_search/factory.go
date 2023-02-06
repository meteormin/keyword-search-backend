package host_search

import (
	"github.com/miniyus/keyword-search-backend/internal/search"
	"github.com/miniyus/keyword-search-backend/pkg/worker"
	"gorm.io/gorm"
)

func New(db *gorm.DB, dispatcher worker.Dispatcher) Handler {
	repo := search.NewRepository(db)
	service := search.NewService(repo)
	return NewHandler(service, dispatcher)
}
