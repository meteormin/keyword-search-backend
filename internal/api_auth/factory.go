package api_auth

import (
	"github.com/miniyus/keyword-search-backend/auth"
	"github.com/miniyus/keyword-search-backend/internal/users"
	"github.com/miniyus/keyword-search-backend/pkg/jwt"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func New(db *gorm.DB, generator jwt.Generator, logger *zap.SugaredLogger) Handler {
	repo := auth.NewRepository(db)
	service := NewService(repo, users.NewRepository(db, logger), generator)
	handler := NewHandler(service)

	return handler
}
