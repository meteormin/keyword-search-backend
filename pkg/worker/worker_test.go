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
		Name:        "default",
		MaxJobCount: 10,
	}
	dispatcherOpt := worker.DispatcherOption{
		WorkerOptions: []worker.Option{
			opt,
		},
		Redis: func() *redis.Client {
			return redis.NewClient(&redis.Options{
				Addr:     "localhost:6379",
				Password: "",
				DB:       0,
			})
		},
	}

	dispatcher := worker.NewDispatcher(dispatcherOpt)
	err := dispatcher.Dispatch("t1", func(job worker.Job) error {
		log.Printf("id %s status %s", job.JobId, job.Status)
		job.Status = worker.SUCCESS
		return nil
	})
	if err != nil {
		log.Print(err)
		return
	}

	err = dispatcher.Dispatch("t2", func(job worker.Job) error {
		log.Printf("id %s status %s", job.JobId, job.Status)
		return nil
	})

	if err != nil {
		log.Print(err)
		return
	}

	dispatcher.Run()
	time.Sleep(time.Second * 3)
	err = dispatcher.Dispatch("t3", func(job worker.Job) error {
		log.Printf("id %s status %s", job.JobId, job.Status)
		return nil
	})

	if err != nil {
		log.Print(err)
		return
	}

	err = dispatcher.Dispatch("t3", func(job worker.Job) error {
		log.Printf("id %s status %s", job.JobId, job.Status)
		return nil
	})

	if err != nil {
		log.Print(err)
		return
	}

	time.Sleep(time.Second * 7)
}
