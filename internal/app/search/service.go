package search

import (
	"github.com/miniyus/go-fiber/internal/entity"
	"github.com/miniyus/go-fiber/internal/utils"
)

type Service interface {
	All(page utils.Page) (utils.Paginator, error)
	GetByHostId(hostId uint, page utils.Page) (utils.Paginator, error)
	Find(pk uint) (*SearchResponse, error)
	Create(search *CreateSearch) (*SearchResponse, error)
}

type ServiceStruct struct {
	repo Repository
}

func NewService(repo Repository) *ServiceStruct {
	return &ServiceStruct{
		repo: repo,
	}
}

func (s *ServiceStruct) All(page utils.Page) (utils.Paginator, error) {
	data, count, err := s.repo.All(page)
	if err != nil {
		data = make([]*entity.Search, 0)
		return utils.Paginator{
			Page:       page,
			TotalCount: 0,
			Data:       data,
		}, err
	}

	return utils.Paginator{
		Page:       page,
		TotalCount: count,
		Data:       data,
	}, err
}

func (s *ServiceStruct) GetByHostId(hostId uint, page utils.Page) (utils.Paginator, error) {
	data, count, err := s.repo.GetByHostId(hostId, page)

	if err != nil {
		data = make([]*entity.Search, 0)
		return utils.Paginator{
			Page:       page,
			TotalCount: 0,
			Data:       data,
		}, err
	}

	var searchRes []*SearchResponse
	for _, s := range data {
		response := ToSearchResponse(s)
		searchRes = append(searchRes, response)
	}

	return utils.Paginator{
		Page:       page,
		TotalCount: count,
		Data:       data,
	}, err
}

func (s *ServiceStruct) Find(pk uint) (*SearchResponse, error) {
	search, err := s.repo.Find(pk)
	if err != nil {
		return nil, err
	}

	searchRes := ToSearchResponse(search)

	return searchRes, err
}

func (s *ServiceStruct) Create(search *CreateSearch) (*SearchResponse, error) {
	ent := entity.Search{
		HostId:      search.HostId,
		Path:        search.Path,
		Query:       search.Query,
		Description: search.Description,
		Publish:     search.Publish,
	}

	rs, err := s.repo.Create(&ent)
	if err != nil {
		return nil, err
	}

	searchRes := ToSearchResponse(rs)
	return searchRes, err
}
