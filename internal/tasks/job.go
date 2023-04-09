package tasks

import (
	"github.com/miniyus/gofiber/app"
	"github.com/miniyus/gofiber/jobqueue"
	worker "github.com/miniyus/goworker"
	"github.com/miniyus/keyword-search-backend/internal/cache"
	"github.com/redis/go-redis/v9"
)

var (
	RedisCacheAll = "cache.redis.all"
)

func RegisterJob(workerOption jobqueue.WorkerOption) app.Register {
	jobqueue.NewContainer(workerOption)

	return func(app app.Application) {
		container := jobqueue.GetContainer()
		container.AddJob(RedisCacheAll, func(job *worker.Job) error {
			var redisClient *redis.Client
			app.Resolve(&redisClient)
			rc := cache.NewRedisCache(redisClient)
			all, err := rc.All()
			if err != nil {
				return err
			}
			job.Meta[job.JobId] = all

			return nil
		})
	}
}
