package api_jobs

import (
	"github.com/go-redis/redis/v9"
	"github.com/miniyus/keyword-search-backend/pkg/worker"
	"go.uber.org/zap"
)

func New(redis func() *redis.Client, dispatcher worker.Dispatcher, zapLogger *zap.SugaredLogger) Handler {
	s := NewService(redis, dispatcher, zapLogger)

	return NewHandler(s)
}
