package search

import (
	"github.com/miniyus/go-fiber/internal/app/hosts"
	"github.com/miniyus/go-fiber/internal/entity"
)

type CreateSearch struct {
	UserId      uint   `json:"user_id"`
	HostId      uint   `json:"host_id"`
	Path        string `json:"path"`
	Query       string `json:"query"`
	Description string `json:"description"`
	Publish     bool   `json:"publish"`
}

type UpdateSearch struct {
	UserId      uint   `json:"user_id"`
	HostId      uint   `json:"host_id"`
	Path        string `json:"path"`
	Query       string `json:"query"`
	Description string `json:"description"`
	Publish     bool   `json:"publish"`
}

type SearchResponse struct {
	Host        *hosts.HostResponse `json:"host"`
	Path        string              `json:"path"`
	Query       string              `json:"query"`
	Description string              `json:"description"`
	Publish     bool                `json:"publish"`
}

func ToSearchResponse(search *entity.Search) *SearchResponse {
	host := hosts.ToHostResponse(search.Host)

	dto := &SearchResponse{
		Host:        host,
		Publish:     search.Publish,
		Path:        search.Path,
		Query:       search.Query,
		Description: search.Description,
	}

	return dto
}
