package auth

import (
	"github.com/miniyus/gofiber/pkg/jwt"
	"github.com/miniyus/keyword-search-backend/entity"
	"gorm.io/gorm"
)

type UserRepository interface {
	Find(pk uint) (*entity.User, error)
	FindByUsername(username string) (*entity.User, error)
	FindByEntity(user entity.User) (*entity.User, error)
	Create(user entity.User) (*entity.User, error)
	Update(pk uint, user entity.User) (*entity.User, error)
}

func New(db *gorm.DB, userRepository UserRepository, generator jwt.Generator) Handler {
	repo := NewRepository(db)
	service := NewService(repo, userRepository, generator)
	handler := NewHandler(service)

	return handler
}
