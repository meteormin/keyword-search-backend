package api_jobs

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v9"
	"github.com/miniyus/keyword-search-backend/logger"
	"github.com/miniyus/keyword-search-backend/pkg/worker"
	"go.uber.org/zap"
)

type Service interface {
	GetJobs(workerName string) ([]worker.Job, error)
	GetJob(workerName string, jobId string) (*worker.Job, error)
	Status() *worker.StatusInfo
	logger.HasLogger
}

type ServiceStruct struct {
	redis      func() *redis.Client
	dispatcher worker.Dispatcher
	logger.HasLoggerStruct
}

func NewService(redis func() *redis.Client, dispatcher worker.Dispatcher, zapLogger *zap.SugaredLogger) Service {
	return &ServiceStruct{
		redis:           redis,
		dispatcher:      dispatcher,
		HasLoggerStruct: logger.HasLoggerStruct{Logger: zapLogger},
	}
}

func (s *ServiceStruct) GetJobs(workerName string) ([]worker.Job, error) {
	keys := fmt.Sprintf("%s.*", workerName)
	redisClient := s.redis()
	result, err := redisClient.Keys(context.Background(), keys).Result()
	if err != nil {
		return nil, err
	}

	jobs := make([]worker.Job, 0)
	for _, r := range result {
		job := worker.Job{}

		val, err := redisClient.Get(context.Background(), r).Result()

		if err == redis.Nil {
			continue
		}

		if err != nil {
			return nil, err
		}

		err = job.UnMarshal(val)
		if err != nil {
			return nil, err
		}

		jobs = append(jobs, job)
	}

	return jobs, nil
}

func (s *ServiceStruct) GetJob(workerName string, jobId string) (*worker.Job, error) {
	key := fmt.Sprintf("%s.%s", workerName, jobId)
	redisClient := s.redis()

	val, err := redisClient.Get(context.Background(), key).Result()

	if err == redis.Nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	job := worker.Job{}

	err = job.UnMarshal(val)
	if err != nil {
		return nil, err
	}

	return &job, nil
}

func (s *ServiceStruct) Status() *worker.StatusInfo {
	return s.dispatcher.Status()
}
