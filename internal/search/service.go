package search

import (
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/gofiber/utils"
	"github.com/miniyus/keyword-search-backend/entity"
	"strconv"
)

type Service interface {
	All(page utils.Page) (utils.Paginator[Response], error)
	GetByHostId(hostId uint, page utils.Page) (utils.Paginator[Response], error)
	GetDescriptionsByHostId(hostId uint, page utils.Page) (utils.Paginator[Description], error)
	Find(pk uint, userId uint) (*Response, error)
	Create(search *CreateSearch) (*Response, error)
	BatchCreate(hostId uint, search []*CreateSearch) ([]Response, error)
	Update(pk uint, userId uint, search *UpdateSearch) (*Response, error)
	Patch(pk uint, userId uint, search *PatchSearch) (*Response, error)
	Delete(pk uint, userId uint) (bool, error)
}

type ServiceStruct struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &ServiceStruct{
		repo: repo,
	}
}

func (s *ServiceStruct) All(page utils.Page) (utils.Paginator[Response], error) {
	data, count, err := s.repo.AllWithPage(page)
	res := make([]Response, 0)
	if err != nil {
		return utils.Paginator[Response]{
			Page:       page,
			TotalCount: 0,
			Data:       res,
		}, err
	}

	var sr Response
	for _, ent := range data {
		res = append(res, sr.FromEntity(ent))
	}

	return utils.Paginator[Response]{
		Page:       page,
		TotalCount: count,
		Data:       res,
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
	var sr Response
	for _, search := range data {
		searchRes = append(searchRes, sr.FromEntity(search))
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
	for _, search := range data {
		response := Description{
			Id:          search.ID,
			Description: search.Description,
			ShortUrl:    *search.ShortUrl,
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

	var sr Response
	searchRes := sr.FromEntity(*search)

	return &searchRes, err
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
	code := utils.Base64UrlEncode(idString)
	rs.ShortUrl = &code

	rs, err = s.repo.Update(rs.ID, *rs)
	if err != nil {
		return nil, err
	}

	var sr Response
	searchRes := sr.FromEntity(*rs)

	return &searchRes, err
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
	var sr Response
	for _, r := range rs {
		idString := strconv.Itoa(int(r.ID))
		code := utils.Base64UrlEncode(idString)
		r.ShortUrl = &code
		res := sr.FromEntity(r)
		resSlice = append(resSlice, res)
		updateSlice = append(updateSlice, r)
	}

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

	ent := search.ToEntity()

	updated, err := s.repo.Update(pk, ent)
	if err != nil {
		return nil, err
	}

	var sr Response
	searchRes := sr.FromEntity(*updated)

	return &searchRes, err
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

	ent := search.ToEntity()

	updated, err := s.repo.Update(pk, ent)
	if err != nil {
		return nil, err
	}

	var sr Response
	searchRes := sr.FromEntity(*updated)

	return &searchRes, err
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
