package groups

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func New(db *gorm.DB, log *zap.SugaredLogger) Handler {
	repo := NewRepository(db, log)
	service := NewService(repo)
	return NewHandler(service)
}
