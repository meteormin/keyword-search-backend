package auth

import (
	"github.com/miniyus/go-fiber/internal/users"
	"github.com/miniyus/go-fiber/pkg/jwt"
	"gorm.io/gorm"
)

func Factory(db *gorm.DB, generator jwt.Generator) *HandlerStruct {
	repo := NewRepository(db)
	service := NewService(repo, users.NewRepository(db), generator)
	handler := NewHandler(service)

	return handler
}
