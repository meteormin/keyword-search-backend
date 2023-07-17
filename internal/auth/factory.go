package auth

import (
	"github.com/miniyus/gofiber/pkg/jwt"
	"github.com/miniyus/keyword-search-backend/entity"
	"github.com/miniyus/keyword-search-backend/repo"
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
	repository := repo.NewAuthRepository(db)
	service := NewService(repository, userRepository, generator)
	handler := NewHandler(service)

	return handler
}
