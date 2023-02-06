package jobs

import (
	"github.com/go-redis/redis/v9"
	"github.com/miniyus/keyword-search-backend/pkg/worker"
)

func New(redis func() *redis.Client, dispatcher worker.Dispatcher) Handler {
	s := NewService(redis, dispatcher)

	return NewHandler(s)
}
