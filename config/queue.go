package config

import "github.com/miniyus/keyword-search-backend/pkg/worker"

func QueueConfig() worker.DispatcherOption {
	return worker.DispatcherOption{
		WorkerOptions: []worker.Option{
			{
				Name:        "default",
				MaxJobCount: 12,
			},
		},
	}
}
