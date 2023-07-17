package host_search

import (
	worker "github.com/miniyus/goworker"
	"github.com/miniyus/keyword-search-backend/internal/search"
	"github.com/miniyus/keyword-search-backend/repo"
	"gorm.io/gorm"
)

func New(db *gorm.DB, dispatcher worker.Dispatcher) Handler {
	repository := repo.NewSearchRepository(db)
	service := search.NewService(repository)
	return NewHandler(service, dispatcher)
}
