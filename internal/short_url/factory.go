package short_url

import (
	"github.com/miniyus/keyword-search-backend/internal/search"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func New(db *gorm.DB, redisClient func() *redis.Client) Handler {
	repo := search.NewRepository(db)
	service := NewService(repo, redisClient)
	return NewHandler(service)
}
