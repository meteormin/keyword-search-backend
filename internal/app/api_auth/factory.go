package api_auth

import (
	"github.com/miniyus/go-fiber/internal/app/users"
	"github.com/miniyus/go-fiber/internal/core/auth"
	"github.com/miniyus/go-fiber/pkg/jwt"
	"gorm.io/gorm"
)

func New(db *gorm.DB, generator jwt.Generator) *HandlerStruct {
	repo := auth.NewRepository(db)
	service := NewService(repo, users.NewRepository(db), generator)
	handler := NewHandler(service)

	return handler
}
