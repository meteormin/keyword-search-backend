package search

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/gofiber/database"
	"github.com/miniyus/gofiber/pagination"
	"github.com/miniyus/gofiber/utils"
	"github.com/miniyus/gorm-extension/gormrepo"
	"github.com/miniyus/keyword-search-backend/entity"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
)

type Service interface {
	All(page pagination.Page) (pagination.Paginator[Response], error)
	GetByHostId(hostId uint, userId uint, query Query) (pagination.Paginator[Response], error)
	GetDescriptionsByHostId(hostId uint, userId uint, query Query) (pagination.Paginator[Description], error)
	Find(pk uint, userId uint) (*Response, error)
	Create(userId uint, search *CreateSearch) (*Response, error)
	BatchCreate(hostId uint, userId uint, search []*CreateSearch) ([]Response, error)
	Update(pk uint, userId uint, search *UpdateSearch) (*Response, error)
	Patch(pk uint, userId uint, search *PatchSearch) (*Response, error)
	Delete(pk uint, userId uint) (bool, error)
	UploadImage(pk uint, userId uint, file *multipart.FileHeader) (*Response, error)
	FindImagePath(pk uint, userId uint) (string, error)
}

type ServiceStruct struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &ServiceStruct{
		repo: repo,
	}
}

func (s *ServiceStruct) All(page pagination.Page) (pagination.Paginator[Response], error) {
	data, count, err := s.repo.AllWithPage(page)
	res := make([]Response, 0)
	if err != nil {
		return pagination.Paginator[Response]{
			Page:       page,
			TotalCount: 0,
			Data:       res,
		}, err
	}

	var sr Response
	for _, ent := range data {
		res = append(res, sr.FromEntity(ent))
	}

	return pagination.Paginator[Response]{
		Page:       page,
		TotalCount: count,
		Data:       res,
	}, err
}

func (s *ServiceStruct) GetByHostId(hostId uint, userId uint, query Query) (pagination.Paginator[Response], error) {
	if !s.repo.HasHost(hostId, userId) {
		return pagination.Paginator[Response]{
			Page:       query.Page,
			TotalCount: 0,
			Data:       make([]Response, 0),
		}, fiber.ErrForbidden
	}

	data, count, err := s.repo.GetByHostId(hostId, Filter{
		Page:     query.Page,
		Query:    query.Query,
		QueryKey: query.QueryKey,
		SortBy:   query.SortBy,
		OrderBy:  query.OrderBy,
	})

	if err != nil {
		return pagination.Paginator[Response]{
			Page:       query.Page,
			TotalCount: 0,
			Data:       make([]Response, 0),
		}, err
	}

	searchRes := make([]Response, 0)
	var sr Response
	for _, search := range data {
		searchRes = append(searchRes, sr.FromEntity(search))
	}

	return pagination.Paginator[Response]{
		Page:       query.Page,
		TotalCount: count,
		Data:       searchRes,
	}, err
}

func (s *ServiceStruct) GetDescriptionsByHostId(hostId uint, userId uint, query Query) (pagination.Paginator[Description], error) {
	if !s.repo.HasHost(hostId, userId) {
		return pagination.Paginator[Description]{
			Page:       query.Page,
			TotalCount: 0,
			Data:       make([]Description, 0),
		}, fiber.ErrForbidden
	}

	data, count, err := s.repo.GetByHostId(hostId, Filter{
		Page:     query.Page,
		Query:    query.Query,
		QueryKey: query.QueryKey,
	})

	if err != nil {
		return pagination.Paginator[Description]{
			Page:       query.Page,
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

	return pagination.Paginator[Description]{
		Page:       query.Page,
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

func (s *ServiceStruct) Create(userId uint, search *CreateSearch) (*Response, error) {
	if !s.repo.HasHost(search.HostId, userId) {
		return nil, fiber.ErrForbidden
	}

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

func (s *ServiceStruct) BatchCreate(hostId uint, userId uint, search []*CreateSearch) ([]Response, error) {
	if !s.repo.HasHost(hostId, userId) {
		return make([]Response, 0), fiber.ErrForbidden
	}

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

func (s *ServiceStruct) UploadImage(pk uint, userId uint, file *multipart.FileHeader) (*Response, error) {
	fileRepo := gormrepo.NewGenericRepository(database.GetDB(), entity.File{})
	savePath := fmt.Sprintf("images/%s", file.Filename)

	ext := ""
	split := strings.Split(file.Filename, ".")
	if len(split) > 1 {
		ext = split[1]
	}

	search, err := s.repo.Find(pk)
	if err != nil {
		return nil, err
	}

	if search.Host.UserId != userId {
		return nil, fiber.ErrForbidden
	}

	f, err := file.Open()
	if err != nil {
		return nil, err
	}

	fh := make([]byte, 512)
	_, err = f.Read(fh)
	if err != nil {
		return nil, err
	}

	mimType := http.DetectContentType(fh)
	//if !strings.Contains(mimType, "image") {
	//	log.Print(mimType)
	//	return nil, fiber.NewError(400, "must upload only image file")
	//}

	create, err := fileRepo.Create(entity.File{
		Path:      savePath,
		MimeType:  mimType,
		Extension: ext,
	})

	if err != nil {
		return nil, err
	}

	search.FileId = &create.ID

	save, err := s.repo.Save(*search)
	if err != nil {
		return nil, err
	}

	var sr Response
	res := sr.FromEntity(*save)

	return &res, nil
}

func (s *ServiceStruct) FindImagePath(pk uint, userId uint) (string, error) {
	find, err := s.repo.Preload("File").Preload("Host").Find(pk)
	if err != nil {
		return "", err
	}

	if find.Host.UserId != userId {
		return "", fiber.ErrForbidden
	}

	if find.File == nil {
		return "", fiber.ErrNotFound
	}

	return find.File.Path, nil
}
