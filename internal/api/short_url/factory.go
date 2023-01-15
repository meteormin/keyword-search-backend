package short_url

import (
	"github.com/go-redis/redis/v9"
	"github.com/miniyus/keyword-search-backend/internal/api/search"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func New(db *gorm.DB, redisClient func() *redis.Client, logger *zap.SugaredLogger) Handler {
	repo := search.NewRepository(db, logger)
	service := NewService(repo, redisClient, logger)
	return NewHandler(service, logger)
}
