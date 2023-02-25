package short_url

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/gofiber/log"
	"github.com/miniyus/keyword-search-backend/internal/search"
	"github.com/redis/go-redis/v9"
	"path"
	"strconv"
	"strings"
	"time"
)

var redisContext = context.Background()

type Service interface {
	FindRealUrl(code string, userId uint) (string, error)
}

type ServiceStruct struct {
	searchRepo search.Repository
	redis      func() *redis.Client
}

func NewService(repository search.Repository, redisClient func() *redis.Client) Service {
	return &ServiceStruct{
		searchRepo: repository,
		redis:      redisClient,
	}
}

func (s *ServiceStruct) hGet(r *redis.Client, rKey string, rField string) string {
	result, err := r.HGet(redisContext, rKey, rField).Result()
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

	r.ExpireGT(redisContext, rKey, time.Hour)

	return result
}

func (s *ServiceStruct) hSet(r *redis.Client, rKey string, rField string, rValue string) error {

	err := r.HSet(redisContext, rKey, rField, rValue).Err()
	if err != nil {
		return err
	}

	r.ExpireNX(redisContext, rKey, time.Hour)

	return nil
}

func (s *ServiceStruct) FindRealUrl(code string, userId uint) (string, error) {
	r := s.redis()
	rKey := "short_url." + strconv.Itoa(int(userId))
	rField := code

	result := s.hGet(r, rKey, rField)
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

	host := searchEnt.Host.Host
	hostPath := searchEnt.Host.Path

	sep := ":/"
	splitString := strings.Split(host, sep)
	hostPath = path.Join(splitString[1], hostPath)

	queryKey := searchEnt.QueryKey
	queryString := searchEnt.Query

	realUrl := fmt.Sprintf("%s:/%s?%s=%s", splitString[0], hostPath, queryKey, queryString)

	err = s.hSet(r, rKey, rField, realUrl)
	if err != nil {
		log.GetLogger().Error(err)
		return realUrl, err
	}

	return realUrl, err
}
