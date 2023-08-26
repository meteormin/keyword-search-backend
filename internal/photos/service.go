package photos

import (
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/go-restapi"
	"github.com/miniyus/keyword-search-backend/entity"
)

type Service struct {
	restapi.Service[entity.Photo, *Request, *Response]
	parentRepository restapi.Repository[entity.Gallery]
}

func NewService(repository restapi.Repository[entity.Photo], parent restapi.Repository[entity.Gallery]) *Service {
	return &Service{
		restapi.NewService[entity.Photo, *Request, *Response](repository, &Response{}),
		parent,
	}
}

func (s *Service) BatchCreate(dtos []Request, userId uint) ([]Response, error) {
	res := make([]Response, 0)
	entities := make([]entity.Photo, 0)
	find, err := s.parentRepository.Find(dtos[0].GalleryId)
	if err != nil {
		return nil, err
	}
	if find.UserId != userId {
		return nil, fiber.ErrForbidden
	}

	for _, dto := range dtos {

		ent := &entity.Photo{}
		err = dto.ToEntity(ent)
		if err != nil {
			return nil, err
		}
		entities = append(entities, *ent)
		r := s.Response()
		restapi.DeepCopy(r)
		err = r.FromEntity(*ent)
		if err != nil {
			return nil, err
		}
		res = append(res, *r)
	}

	db := s.Repo().DB().CreateInBatches(entities, 100)
	if db.Error != nil {
		return nil, db.Error
	}
	return res, nil
}

func (s *Service) Create(dto *Request) (*Response, error) {
	_, err := s.parentRepository.Find(dto.GalleryId)
	if err != nil {
		return nil, err
	}

	ent := &entity.Photo{}
	err = dto.ToEntity(ent)
	if err != nil {
		return nil, err
	}

	create, err := s.Repo().Create(*ent)
	if err != nil {
		return nil, err
	}
	res := s.Response()
	err = res.FromEntity(*create)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *Service) FindImage(pk uint, userId uint) (*Response, error) {
	find, err := s.Repo().Preload("Gallery").Find(pk)
	if err != nil {
		return nil, err
	}
	if find.Gallery.UserId != userId {
		return nil, fiber.ErrForbidden
	}

	res := s.Response()
	err = res.FromEntity(*find)
	if err != nil {
		return nil, err
	}

	return res, nil
}
