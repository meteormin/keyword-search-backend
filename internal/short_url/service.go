package short_url

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/gofiber/log"
	"github.com/miniyus/keyword-search-backend/repo"
	"github.com/miniyus/keyword-search-backend/utils"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

var redisContext = context.Background()

type Service interface {
	FindRealUrl(code string, userId uint) (string, error)
}

type ServiceStruct struct {
	searchRepo repo.SearchRepository
	makeRedis  func() *redis.Client
	redis      *redis.Client
}

func NewService(repository repo.SearchRepository, redisClient func() *redis.Client) Service {
	return &ServiceStruct{
		searchRepo: repository,
		makeRedis:  redisClient,
		redis:      redisClient(),
	}
}

func (s *ServiceStruct) hGet(rKey string, rField string) string {
	result, err := s.redis.HGet(redisContext, rKey, rField).Result()
	if err == redis.Nil {
		return ""
	}

	if err != nil {
		log.GetLogger().Error(err)
		return ""
	}

	if result == "" {
		return ""
	}

	s.redis.ExpireGT(redisContext, rKey, time.Hour)

	return result
}

func (s *ServiceStruct) hSet(rKey string, rField string, rValue string) error {

	err := s.redis.HSet(redisContext, rKey, rField, rValue).Err()
	if err != nil {
		return err
	}

	s.redis.ExpireNX(redisContext, rKey, time.Hour)

	return nil
}

func (s *ServiceStruct) FindRealUrl(code string, userId uint) (string, error) {
	rKey := "short_url." + strconv.Itoa(int(userId))
	rField := code

	result := s.hGet(rKey, rField)
	if result != "" {
		return result, nil
	}

	searchEnt, err := s.searchRepo.FindByShortUrl(code, userId)
	if err != nil {
		return "", err
	}

	if searchEnt == nil {
		return "", fiber.ErrNotFound
	}

	searchEnt.Views += 1
	save, err := s.searchRepo.Save(*searchEnt)
	if err != nil {
		return "", err
	}

	host := save.Host.Host
	hostPath := save.Host.Path
	queryKey := save.QueryKey
	queryString := save.Query

	realUrl := utils.MakeRealUrl(host, hostPath, queryKey, queryString)

	err = s.hSet(rKey, rField, realUrl)
	if err != nil {
		log.GetLogger().Error(err)
		return realUrl, err
	}

	return realUrl, err
}
