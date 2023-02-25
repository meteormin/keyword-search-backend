package config

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
	"strconv"
)

func redisConfig() *redis.Options {
	addr := os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT")

	redisDb, err := strconv.Atoi(os.Getenv("REDIS_DATABASE"))
	if err != nil {
		redisDb = 0
	}

	return &redis.Options{
		Addr:     addr,
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       redisDb,
		OnConnect: func(ctx context.Context, cn *redis.Conn) error {
			log.Print("connect redis")
			err := ctx.Err()
			if err != nil {
				log.Print(err)
			}

			return err
		},
	}
}
