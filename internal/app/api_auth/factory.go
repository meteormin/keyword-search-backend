package api_auth

import (
	"github.com/miniyus/go-fiber/internal/app/users"
	"github.com/miniyus/go-fiber/internal/core/auth"
	"github.com/miniyus/go-fiber/pkg/jwt"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func New(db *gorm.DB, generator jwt.Generator, logger *zap.SugaredLogger) Handler {
	repo := auth.NewRepository(db, logger)
	service := NewService(repo, users.NewRepository(db, repo.GetLogger()), generator)
	handler := NewHandler(service)

	return handler
}
