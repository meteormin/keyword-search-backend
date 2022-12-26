package api_auth

import (
	"github.com/miniyus/go-fiber/internal/app/users"
	"github.com/miniyus/go-fiber/internal/core/auth"
	"github.com/miniyus/go-fiber/pkg/jwt"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func New(db *gorm.DB, generator jwt.Generator, logger *zap.SugaredLogger) Handler {
	repo := auth.NewRepository(db)
	service := NewService(repo, users.NewRepository(db, logger), generator)
	handler := NewHandler(service)

	return handler
}
