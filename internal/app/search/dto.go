package search

import (
	"github.com/miniyus/keyword-search-backend/internal/entity"
	"github.com/miniyus/keyword-search-backend/internal/utils"
)

type MultiCreateSearch struct {
	Search []*CreateSearch `json:"search" validate:"required"`
}

type CreateSearch struct {
	HostId      uint   `json:"host_id" validate:"required"`
	QueryKey    string `json:"query_key" validate:"required"`
	Query       string `json:"query" validate:"required"`
	Description string `json:"description" validate:"required"`
	Publish     bool   `json:"publish" validate:"required,boolean"`
}

type UpdateSearch struct {
	HostId      uint   `json:"host_id" validate:"required"`
	QueryKey    string `json:"query_key" validate:"required"`
	Query       string `json:"query" validate:"required"`
	Description string `json:"description" validate:"required"`
	Publish     bool   `json:"publish" validate:"required,boolean"`
}

type PatchSearch struct {
	HostId      uint    `json:"host_id" validate:"required"`
	QueryKey    *string `json:"query_key,omitempty"`
	Query       *string `json:"query,omitempty"`
	Description *string `json:"description,omitempty"`
	Publish     *bool   `json:"publish,omitempty" validate:"omitempty,required,boolean"`
}

type Response struct {
	Id          uint    `json:"id"`
	ShortUrl    *string `json:"short_url"`
	QueryKey    string  `json:"query_key"`
	Query       string  `json:"query"`
	Description string  `json:"description"`
	Publish     bool    `json:"publish"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

type ResponseAll struct {
	utils.Paginator
	Data []entity.Search `json:"data"`
}

func ToSearchResponse(search *entity.Search) *Response {
	createdAt := utils.TimeIn(search.CreatedAt, "Asia/Seoul")
	updatedAt := utils.TimeIn(search.UpdatedAt, "Asia/Seoul")

	dto := &Response{
		Id:          search.ID,
		ShortUrl:    search.ShortUrl,
		Publish:     search.Publish,
		QueryKey:    search.QueryKey,
		Query:       search.Query,
		Description: search.Description,
		CreatedAt:   createdAt.Format(utils.DefaultDateLayout),
		UpdatedAt:   updatedAt.Format(utils.DefaultDateLayout),
	}

	return dto
}
