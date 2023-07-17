package repo

import (
	"github.com/miniyus/gorm-extension/gormrepo"
	"github.com/miniyus/keyword-search-backend/entity"
	"gorm.io/gorm"
)

// AuthRepository
// JWT 관련 토큰 Repository
type AuthRepository interface {
	gormrepo.GenericRepository[entity.AccessToken]
	FindByToken(token string) (*entity.AccessToken, error)
	FindByUserId(userId uint) (*entity.AccessToken, error)
	Delete(pk uint) (bool, error)
}

// AuthRepositoryStruct
// Repository 인터페이스 구현 구조체
type AuthRepositoryStruct struct {
	gormrepo.GenericRepository[entity.AccessToken]
}

// NewAuthRepository
// Repository 생성 함수
func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &AuthRepositoryStruct{gormrepo.NewGenericRepository(db, entity.AccessToken{})}
}

// FindByToken
// find token by token(string)
func (repo *AuthRepositoryStruct) FindByToken(token string) (*entity.AccessToken, error) {
	return repo.FindByEntity(entity.AccessToken{Token: token})
}

// FindByUserId
// find token by user id
func (repo *AuthRepositoryStruct) FindByUserId(userId uint) (*entity.AccessToken, error) {
	return repo.FindByEntity(entity.AccessToken{UserId: userId})
}
