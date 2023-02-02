package host_search

import (
	"github.com/miniyus/keyword-search-backend/internal/search"
	"github.com/miniyus/keyword-search-backend/pkg/worker"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func New(db *gorm.DB, logger *zap.SugaredLogger, dispatcher worker.Dispatcher) Handler {
	repo := search.NewRepository(db, logger)
	service := search.NewService(repo)
	return NewHandler(service, dispatcher)
}
