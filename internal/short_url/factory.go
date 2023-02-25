package short_url

import (
	"fmt"
	"github.com/miniyus/keyword-search-backend/internal/search"
	"github.com/miniyus/keyword-search-backend/utils"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func New(db *gorm.DB, redisClient func() *redis.Client) Handler {
	repo := search.NewRepository(db)
	service := NewService(repo, redisClient)
	return NewHandler(service)
}

type realUrlMaker struct {
	Host     string
	Path     string
	QueryKey string
	Query    string
}

func (rm *realUrlMaker) makeRealUrl() string {
	var realUrl string
	host := utils.JoinHostPath(rm.Host, rm.Path)

	if rm.QueryKey == "" {
		realUrl = fmt.Sprintf("%s/%s", host, rm.Query)
	} else {
		realUrl = utils.AddQueryString(host, map[string]interface{}{
			rm.QueryKey: rm.Query,
		})
	}

	return realUrl
}

func MakeRealUrl(host string, path string, queryKey string, query string) string {
	return (&realUrlMaker{Host: host, Path: path, QueryKey: queryKey, Query: query}).makeRealUrl()
}
