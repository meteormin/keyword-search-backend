package utils

import "github.com/go-redis/redis/v9"

type StatusResponse struct {
	Status bool `json:"status"`
}

type DataResponse[T interface{}] struct {
	Data T `json:"data"`
}

func RedisClientMaker(options *redis.Options) func() *redis.Client {
	return func() *redis.Client {
		return redis.NewClient(options)
	}
}
