package hosts

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func New(db *gorm.DB, logger *zap.SugaredLogger) Handler {
	repo := NewRepository(db, logger)
	service := NewService(repo)
	return NewHandler(service)
}
