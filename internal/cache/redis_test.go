package cache_test

import (
	"github.com/miniyus/keyword-search-backend/internal/cache"
	"github.com/redis/go-redis/v9"
	"log"
	"testing"
)

var redisClient *redis.Client

func init() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})
}

func TestRedisCache_All(t *testing.T) {
	rc := cache.NewRedisCache(redisClient)

	all, err := rc.All()
	if err != nil {
		t.Error(err)
	}
	marshal, err := all.Marshal()
	if err != nil {
		t.Error(err)
	}

	log.Print(marshal)
}
