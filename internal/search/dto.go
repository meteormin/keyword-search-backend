package search

import (
	"github.com/miniyus/gofiber/pagination"
	"github.com/miniyus/gofiber/utils"
	"github.com/miniyus/keyword-search-backend/entity"
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

func (us UpdateSearch) ToEntity() entity.Search {
	return entity.Search{
		HostId:      us.HostId,
		QueryKey:    us.QueryKey,
		Query:       us.Query,
		Description: us.Description,
		Publish:     us.Publish,
	}
}

type PatchSearch struct {
	QueryKey    *string `json:"query_key,omitempty"`
	Query       *string `json:"query,omitempty"`
	Description *string `json:"description,omitempty"`
	Publish     *bool   `json:"publish,omitempty" validate:"omitempty,required,boolean"`
}

func (ps PatchSearch) ToEntity() entity.Search {
	var ent entity.Search

	if ps.Query != nil {
		ent.Query = *ps.Query
	}

	if ps.QueryKey != nil {
		ent.QueryKey = *ps.QueryKey
	}

	if ps.Publish != nil {
		ent.Publish = *ps.Publish
	}

	if ps.Description != nil {
		ent.Description = *ps.Description
	}

	return ent
}

type Response struct {
	Id          uint    `json:"id"`
	HostId      uint    `json:"host_id"`
	ShortUrl    *string `json:"short_url"`
	QueryKey    string  `json:"query_key"`
	Query       string  `json:"query"`
	Description string  `json:"description"`
	Publish     bool    `json:"publish"`
	Views       uint    `json:"views"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

type Description struct {
	Id          uint   `json:"id"`
	Description string `json:"description"`
	ShortUrl    string `json:"short_url"`
}

type ResponseByHost struct {
	pagination.Paginator[Response]
	Data []Response `json:"data"`
}

type DescriptionWithPage struct {
	pagination.Paginator[Description]
	Data []Description `json:"data"`
}

type ResponseAll struct {
	pagination.Paginator[Response]
	Data []Response `json:"data"`
}

func (r Response) FromEntity(search entity.Search) Response {
	createdAt := utils.TimeIn(search.CreatedAt, "Asia/Seoul")
	updatedAt := utils.TimeIn(search.UpdatedAt, "Asia/Seoul")

	return Response{
		Id:          search.ID,
		HostId:      search.HostId,
		ShortUrl:    search.ShortUrl,
		Publish:     search.Publish,
		QueryKey:    search.QueryKey,
		Query:       search.Query,
		Description: search.Description,
		CreatedAt:   createdAt.Format(utils.DefaultDateLayout),
		UpdatedAt:   updatedAt.Format(utils.DefaultDateLayout),
	}
}

type Query struct {
	Page     pagination.Page `query:"-"`
	QueryKey *string         `json:"query_key" query:"query_key"`
	Query    *string         `json:"query" query:"query"`
	SortBy   *string         `json:"sort_by" query:"sort_by"`
	OrderBy  *string         `json:"order_by" query:"order_by"`
}
