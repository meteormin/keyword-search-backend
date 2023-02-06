package short_url

import (
	"github.com/go-redis/redis/v9"
	"github.com/miniyus/keyword-search-backend/internal/search"
	"gorm.io/gorm"
)

func New(db *gorm.DB, redisClient func() *redis.Client) Handler {
	repo := search.NewRepository(db)
	service := NewService(repo, redisClient)
	return NewHandler(service)
}
