package short_url

import (
	"github.com/miniyus/keyword-search-backend/repo"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func New(db *gorm.DB, redisClient func() *redis.Client) Handler {
	repository := repo.NewSearchRepository(db)
	service := NewService(repository, redisClient)
	return NewHandler(service)
}
