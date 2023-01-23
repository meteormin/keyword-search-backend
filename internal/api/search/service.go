package search

import (
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/keyword-search-backend/entity"
	"github.com/miniyus/keyword-search-backend/logger"
	"github.com/miniyus/keyword-search-backend/utils"
	"log"
	"strconv"
)

type Service interface {
	All(page utils.Page) (utils.Paginator[entity.Search], error)
	GetByHostId(hostId uint, page utils.Page) (utils.Paginator[Response], error)
	GetDescriptionsByHostId(hostId uint, page utils.Page) (utils.Paginator[Description], error)
	Find(pk uint, userId uint) (*Response, error)
	Create(search *CreateSearch) (*Response, error)
	BatchCreate(hostId uint, search []*CreateSearch) ([]Response, error)
	Update(pk uint, userId uint, search *UpdateSearch) (*Response, error)
	Patch(pk uint, userId uint, search *PatchSearch) (*Response, error)
	Delete(pk uint, userId uint) (bool, error)
	logger.HasLogger
}

type ServiceStruct struct {
	repo Repository
	logger.HasLoggerStruct
}

func NewService(repo Repository) Service {
	return &ServiceStruct{
		repo: repo,
		HasLoggerStruct: logger.HasLoggerStruct{
			Logger: repo.GetLogger(),
		},
	}
}

func (s *ServiceStruct) All(page utils.Page) (utils.Paginator[entity.Search], error) {
	data, count, err := s.repo.All(page)
	if err != nil {
		data = make([]entity.Search, 0)
		return utils.Paginator[entity.Search]{
			Page:       page,
			TotalCount: 0,
			Data:       data,
		}, err
	}

	return utils.Paginator[entity.Search]{
		Page:       page,
		TotalCount: count,
		Data:       data,
	}, err
}

func (s *ServiceStruct) GetByHostId(hostId uint, page utils.Page) (utils.Paginator[Response], error) {
	data, count, err := s.repo.GetByHostId(hostId, page)

	if err != nil {
		return utils.Paginator[Response]{
			Page:       page,
			TotalCount: 0,
			Data:       make([]Response, 0),
		}, err
	}

	searchRes := make([]Response, 0)
	for _, s := range data {
		response := ToSearchResponse(&s)
		searchRes = append(searchRes, *response)
	}

	return utils.Paginator[Response]{
		Page:       page,
		TotalCount: count,
		Data:       searchRes,
	}, err
}

func (s *ServiceStruct) GetDescriptionsByHostId(hostId uint, page utils.Page) (utils.Paginator[Description], error) {
	data, count, err := s.repo.GetByHostId(hostId, page)

	if err != nil {
		return utils.Paginator[Description]{
			Page:       page,
			TotalCount: 0,
			Data:       make([]Description, 0),
		}, err
	}

	searchRes := make([]Description, 0)
	for _, s := range data {
		response := Description{
			Id:          s.ID,
			Description: s.Description,
			ShortUrl:    *s.ShortUrl,
		}

		searchRes = append(searchRes, response)
	}

	return utils.Paginator[Description]{
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
	for _, crateSearch := range search {
		ent := entity.Search{
			HostId:      hostId,
			QueryKey:    crateSearch.QueryKey,
			Query:       crateSearch.Query,
			Description: crateSearch.Description,
			Publish:     crateSearch.Publish,
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

func (s *ServiceStruct) Update(pk uint, userId uint, search *UpdateSearch) (*Response, error) {
	exists, err := s.repo.Find(pk)
	if err != nil {
		return nil, err
	}

	if exists == nil {
		return nil, fiber.ErrNotFound
	}

	if exists.Host.UserId != userId {
		return nil, fiber.ErrForbidden
	}

	ent := entity.Search{
		HostId:      search.HostId,
		QueryKey:    search.QueryKey,
		Query:       search.Query,
		Description: search.Description,
		Publish:     search.Publish,
	}

	updated, err := s.repo.Update(pk, ent)
	if err != nil {
		return nil, err
	}

	return ToSearchResponse(updated), err
}

func (s *ServiceStruct) Patch(pk uint, userId uint, search *PatchSearch) (*Response, error) {
	exists, err := s.repo.Find(pk)
	if err != nil {
		return nil, err
	}

	if exists == nil {
		return nil, fiber.ErrNotFound
	}

	if exists.Host.UserId != userId {
		return nil, fiber.ErrForbidden
	}

	ent := entity.Search{}
	if search.HostId != 0 {
		ent.HostId = search.HostId
	}

	if search.Query != nil {
		ent.Query = *search.Query
	}

	if search.QueryKey != nil {
		ent.QueryKey = *search.QueryKey
	}

	if search.Publish != nil {
		ent.Publish = *search.Publish
	}

	if search.Description != nil {
		ent.Description = *search.Description
	}

	updated, err := s.repo.Update(pk, ent)
	if err != nil {
		return nil, err
	}

	return ToSearchResponse(updated), err
}

func (s *ServiceStruct) Delete(pk uint, userId uint) (bool, error) {
	exists, err := s.repo.Find(pk)
	if err != nil {
		return false, err
	}

	if exists == nil {
		return false, fiber.ErrNotFound
	}

	if exists.Host.UserId != userId {
		return false, fiber.ErrForbidden
	}

	return s.repo.Delete(pk)
}
