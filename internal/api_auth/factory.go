package api_auth

import (
	"github.com/miniyus/gofiber/auth"
	"github.com/miniyus/gofiber/pkg/jwt"
	"github.com/miniyus/keyword-search-backend/internal/users"
	"gorm.io/gorm"
)

func New(db *gorm.DB, generator jwt.Generator) Handler {
	repo := auth.NewRepository(db)
	service := NewService(repo, users.NewRepository(db), generator)
	handler := NewHandler(service)

	return handler
}
