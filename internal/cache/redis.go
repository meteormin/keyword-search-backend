package cache

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type RedisStruct struct {
	String map[string]string            `json:"string"`
	List   map[string][]string          `json:"list"`
	Hash   map[string]map[string]string `json:"hash"`
	Sets   map[string][]string          `json:"sets"`
}

func (rs *RedisStruct) Marshal() (string, error) {
	marshal, err := json.Marshal(rs)
	if err != nil {
		return "", err
	}

	return string(marshal), nil
}

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{
		client: client,
	}
}

func (rc *RedisCache) All() (*RedisStruct, error) {
	redisStruct := &RedisStruct{}

	allKeys, err := rc.client.Keys(ctx, "*").Result()
	if err != nil {
		return nil, err
	}

	for _, key := range allKeys {
		result, err := rc.client.Type(ctx, key).Result()
		if err != nil {
			return nil, err
		}
		switch result {
		case "string":
			redisStruct.String = map[string]string{}
			str, err := rc.client.Get(ctx, key).Result()
			if err != nil {
				return nil, err
			}
			redisStruct.String[key] = str
		case "list":
			redisStruct.List = map[string][]string{}
			lLen, err := rc.client.LLen(ctx, key).Result()
			if err != nil {
				return nil, err
			}
			list, err := rc.client.LRange(ctx, key, 0, lLen-1).Result()
			redisStruct.List[key] = list
			break
		case "hash":
			redisStruct.Hash = map[string]map[string]string{}
			hGetAll, err := rc.client.HGetAll(ctx, key).Result()
			if err != nil {
				return nil, err
			}
			redisStruct.Hash[key] = hGetAll
			break
		default:
			break
		}
	}

	return redisStruct, nil
}
