package short_url

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/go-fiber/internal/app/search"
	"path"
	"strings"
)

type Service interface {
	FindRealUrl(code string, userId uint) (string, error)
}

type ServiceStruct struct {
	searchRepo search.Repository
}

func NewService(repository search.Repository) Service {
	return &ServiceStruct{searchRepo: repository}
}

func (s *ServiceStruct) FindRealUrl(code string, userId uint) (string, error) {
	searchEnt, err := s.searchRepo.FindByShortUrl(code, userId)
	if err != nil {
		return "", err
	}

	if searchEnt == nil {
		return "", fiber.ErrNotFound
	}

	host := searchEnt.Host.Host
	hostPath := searchEnt.Host.Path

	sep := "://"
	splitString := strings.Split(host, sep)
	hostPath = path.Join(splitString[1], hostPath)

	queryKey := searchEnt.QueryKey
	queryString := searchEnt.Query

	realUrl := fmt.Sprintf("%s://%s?%s=%s", splitString[0], hostPath, queryKey, queryString)

	return realUrl, err
}
