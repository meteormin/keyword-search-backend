package config

import worker "github.com/miniyus/goworker"

type WorkerName string

const (
	DefaultWorker WorkerName = "default"
)

func jobQueueConfig() worker.DispatcherOption {
	return worker.DispatcherOption{
		WorkerOptions: []worker.Option{
			{
				Name:        string(DefaultWorker),
				MaxJobCount: 12,
			},
		},
	}
}
