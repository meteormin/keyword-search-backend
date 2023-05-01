package tasks

import (
	"errors"
	"github.com/miniyus/gofiber/app"
	"github.com/miniyus/gofiber/database"
	"github.com/miniyus/gofiber/entity"
	"github.com/miniyus/gofiber/jobqueue"
	"github.com/miniyus/gofiber/log"
	"github.com/miniyus/gorm-extension/gormrepo"
	worker "github.com/miniyus/goworker"
	"github.com/miniyus/keyword-search-backend/config"
	"github.com/miniyus/keyword-search-backend/internal/cache"
	"github.com/redis/go-redis/v9"
	"time"
)

var jobs *Jobs

var (
	redisCacheAll   = "cache.redis.all"
	pruneJobHistory = "job_history.prune"
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
	var redisClient *redis.Client
	app.Resolve(&redisClient)
	rc := cache.NewRedisCache(redisClient)

	container.AddJob(redisCacheAll, func(job *worker.Job) error {
		all, err := rc.All()
		if err != nil {
			return err
		}
		job.Meta[job.JobId] = map[string]interface{}{
			"string": all.String,
			"list":   all.List,
			"hash":   all.Hash,
			"sets":   all.Sets,
		}

		return nil
	})

	repo := gormrepo.NewGenericRepository(database.GetDB(), entity.JobHistory{})
	container.AddJob(pruneJobHistory, func(job *worker.Job) error {
		repo.DB().Where("created_at BETWEEN ? AND ?")
		return nil
	})

	jobs = &Jobs{container: container}
}

type Jobs struct {
	container jobqueue.Container
}

func GetJobs() *Jobs {
	return jobs
}

func (j *Jobs) Dispatch(jobId string) error {
	return j.container.Dispatch(jobId)
}

func (j *Jobs) SyncDispatch(jobId string) (interface{}, error) {
	dispatch, err := j.container.SyncDispatch(jobId)
	if err != nil {
		return nil, err
	}

	return dispatch.Meta[dispatch.JobId], nil
}

func (j *Jobs) GetRedisCacheAll() (*cache.RedisStruct, error) {
	data, err := j.SyncDispatch(redisCacheAll)
	if err != nil {
		return nil, err
	}

	if redisStruct, ok := data.(*cache.RedisStruct); ok {
		return redisStruct, nil
	} else {
		return nil, errors.New("invalid Type in GetRedisCacheAll()")
	}
}
