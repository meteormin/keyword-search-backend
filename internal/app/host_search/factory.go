package host_search

import (
	"github.com/miniyus/go-fiber/internal/app/search"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func New(db *gorm.DB, logger *zap.SugaredLogger) Handler {
	repo := search.NewRepository(db, logger)
	service := search.NewService(repo)
	return NewHandler(service)
}
