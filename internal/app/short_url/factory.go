package short_url

import (
	"github.com/miniyus/go-fiber/internal/app/search"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func New(db *gorm.DB, logger *zap.SugaredLogger) Handler {
	repo := search.NewRepository(db, logger)
	service := NewService(repo, logger)
	return NewHandler(service, logger)
}
