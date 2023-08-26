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
	"gorm.io/gorm"
	goLog "log"
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
	err := app.Resolve(&redisClient)
	if err != nil {
		goLog.Print(err)
	}
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
		get, err := repo.Get(func(tx *gorm.DB) (*gorm.DB, error) {
			now := time.Now()
			year, month, _ := now.Date()
			loc := now.Location()
			currentMonth := time.Date(year, month, 1, 0, 0, 0, 0, loc)
			tx.Where("created_at < ?", currentMonth)
			return tx, nil
		})

		if err != nil {
			return err
		}

		job.Meta[job.JobId] = get

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

func (j *Jobs) GetPruneJobHistory() ([]entity.JobHistory, error) {
	data, err := j.SyncDispatch(pruneJobHistory)
	if err != nil {
		return nil, err
	}

	if histories, ok := data.([]entity.JobHistory); ok {
		return histories, nil
	} else {
		return nil, errors.New("invalid Type in GetPruneJobHistory")
	}
}
