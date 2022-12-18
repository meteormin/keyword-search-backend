package search

import (
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/go-fiber/internal/entity"
	"github.com/miniyus/go-fiber/internal/utils"
	"log"
	"strconv"
)

type Service interface {
	All(page utils.Page) (utils.Paginator, error)
	GetByHostId(hostId uint, page utils.Page) (utils.Paginator, error)
	Find(pk uint, userId uint) (*Response, error)
	Create(search *CreateSearch) (*Response, error)
	BatchCreate(hostId uint, search []*CreateSearch) ([]Response, error)
}

type ServiceStruct struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &ServiceStruct{
		repo: repo,
	}
}

func (s *ServiceStruct) All(page utils.Page) (utils.Paginator, error) {
	data, count, err := s.repo.All(page)
	if err != nil {
		data = make([]entity.Search, 0)
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
		data = make([]entity.Search, 0)
		return utils.Paginator{
			Page:       page,
			TotalCount: 0,
			Data:       data,
		}, err
	}

	var searchRes []Response
	for _, s := range data {
		response := ToSearchResponse(&s)
		searchRes = append(searchRes, *response)
	}

	return utils.Paginator{
		Page:       page,
		TotalCount: count,
		Data:       searchRes,
	}, err
}

func (s *ServiceStruct) Find(pk uint, userId uint) (*Response, error) {
	search, err := s.repo.Find(pk)
	if err != nil {
		return nil, err
	}

	if search.Host.UserId != userId {
		return nil, fiber.ErrForbidden
	}

	searchRes := ToSearchResponse(search)

	return searchRes, err
}

func (s *ServiceStruct) Create(search *CreateSearch) (*Response, error) {
	ent := entity.Search{
		HostId:      search.HostId,
		QueryKey:    search.QueryKey,
		Query:       search.Query,
		Description: search.Description,
		Publish:     search.Publish,
	}

	rs, err := s.repo.Create(ent)

	if err != nil {
		return nil, err
	}
	idString := strconv.Itoa(int(rs.ID))
	code := utils.B64UrlEncode(idString)
	rs.ShortUrl = &code

	rs, err = s.repo.Update(rs.ID, *rs)
	if err != nil {
		return nil, err
	}

	searchRes := ToSearchResponse(rs)
	return searchRes, err
}

func (s *ServiceStruct) BatchCreate(hostId uint, search []*CreateSearch) ([]Response, error) {
	entities := make([]entity.Search, 0)
	for _, s := range search {
		ent := entity.Search{
			HostId:      hostId,
			QueryKey:    s.QueryKey,
			Query:       s.Query,
			Description: s.Description,
			Publish:     s.Publish,
		}
		entities = append(entities, ent)
	}

	rs, err := s.repo.BatchCreate(entities)

	resSlice := make([]Response, 0)
	if err != nil {
		return resSlice, err
	}

	updateSlice := make([]entity.Search, 0)
	for _, r := range rs {
		idString := strconv.Itoa(int(r.ID))
		code := utils.B64UrlEncode(idString)
		r.ShortUrl = &code
		res := ToSearchResponse(&r)
		resSlice = append(resSlice, *res)
		updateSlice = append(updateSlice, r)
	}
	log.Print(updateSlice)
	rs, err = s.repo.BatchCreate(updateSlice)

	if err != nil {
		return make([]Response, 0), err
	}

	return resSlice, nil
}
