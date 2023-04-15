package tasks

import (
	"errors"
	"github.com/miniyus/gofiber/app"
	"github.com/miniyus/gofiber/jobqueue"
	"github.com/miniyus/gofiber/log"
	worker "github.com/miniyus/goworker"
	"github.com/miniyus/keyword-search-backend/config"
	"github.com/miniyus/keyword-search-backend/internal/cache"
	"github.com/redis/go-redis/v9"
	"time"
)

var jobs *Jobs

var (
	redisCacheAll = "cache.redis.all"
)

func RegisterJob(app app.Application) {
	var cfg *config.Configs
	app.Resolve(&cfg)

	log.New(log.Config{
		Name:     "tasks_job",
		FilePath: cfg.Path.LogPath,
		Filename: "tasks_job.log",
	})

	jobqueue.NewContainer(jobqueue.WorkerOption{
		Logger:      log.GetLogger("tasks_job"),
		MaxJobCount: 100,
		Delay:       time.Second,
	})

	container := jobqueue.GetContainer()

	container.AddJob(redisCacheAll, func(job *worker.Job) error {
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

	jobs = &Jobs{container: container}
}

type Jobs struct {
	container jobqueue.Container
}

func (j *Jobs) GetRedisCacheAll() (*cache.RedisStruct, error) {
	job, err := j.container.SyncDispatch(redisCacheAll)
	if err != nil {
		return nil, err
	}

	all := job.Meta[job.JobId]
	if redisStruct, ok := all.(*cache.RedisStruct); ok {
		return redisStruct, nil
	} else {
		return nil, errors.New("invalid Type in GetRedisCacheAll()")
	}
}
