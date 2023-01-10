package worker_test

import (
	"github.com/go-redis/redis/v9"
	"github.com/miniyus/keyword-search-backend/pkg/worker"
	"log"
	"testing"
	"time"
)

func TestJobDispatcher_Dispatch(t *testing.T) {
	opt := worker.Option{
		Name: "test_workers",
	}
	dispatcherOpt := worker.DispatcherOption{
		WorkerOptions: []worker.Option{
			opt,
		},
		Redis: redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		}),
	}

	dispatcher := worker.NewDispatcher(dispatcherOpt)
	dispatcher.Dispatch("test_workers", worker.NewJob("t1", func(job worker.Job) error {
		log.Printf("id %s status %s", job.JobId, job.Status)
		job.Status = worker.SUCCESS
		return nil
	}))

	dispatcher.Dispatch("test_workers", worker.NewJob("t2", func(job worker.Job) error {
		log.Printf("id %s status %s", job.JobId, job.Status)
		return nil
	}))

	dispatcher.Run()

	time.Sleep(time.Second * 7)
}
