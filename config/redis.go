package config

import (
	"github.com/go-redis/redis/v9"
	"os"
	"strconv"
)

func redisConfig() *redis.Options {
	redisDb, err := strconv.Atoi(os.Getenv("REDIS_DATABASE"))
	if err != nil {
		redisDb = 0
	}

	return &redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       redisDb,
	}
}
