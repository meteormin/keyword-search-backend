package entity

import (
	"context"
	"github.com/miniyus/gofiber/app"
	"github.com/miniyus/gorm-extension/gormhooks"
	"github.com/miniyus/keyword-search-backend/utils"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"strconv"
)

type Search struct {
	gorm.Model
	Host        *Host   `json:"host" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	HostId      uint    `json:"host_id" gorm:"index:query_unique,unique"`
	QueryKey    string  `gorm:"column:query_key;type:varchar(50);index:query_unique,unique" json:"query_key"`
	Query       string  `gorm:"column:query;type:varchar(50);index:query_unique,unique" json:"query"`
	Description string  `gorm:"column:description;type:varchar(50)" json:"description"`
	Publish     bool    `gorm:"column:publish;type:bool" json:"publish"`
	Views       uint    `json:"views" gorm:"column:views;default:0"`
	ShortUrl    *string `gorm:"column:short_url;type:varchar(255);uniqueIndex" json:"short_url"`
	File        *File   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	FileId      uint
}

func (s *Search) Hooks() *gormhooks.Hooks[*Search] {
	return gormhooks.GetHooks(s)
}

func (s *Search) AfterSave(tx *gorm.DB) (err error) {
	return s.Hooks().AfterSave(tx)
}

type SearchHookHandler struct {
	app app.Application
}

func newSearchHookHandler(app app.Application) *SearchHookHandler {
	return &SearchHookHandler{app: app}
}

func (sh *SearchHookHandler) AfterSave(s *Search, tx *gorm.DB) (err error) {
	a := sh.app

	var rClient *redis.Client
	a.Resolve(&rClient)

	if s.ShortUrl != nil {
		var h Host

		err = tx.First(&h, s.HostId).Error
		if err != nil {
			return err
		}

		rKey := "short_url." + strconv.Itoa(int(h.UserId))
		cached, err := rClient.HGet(
			context.Background(),
			rKey,
			*s.ShortUrl,
		).Result()

		if cached != "" && err == nil {
			realUrl := utils.MakeRealUrl(h.Host, h.Path, s.QueryKey, s.Query)
			rClient.HSet(
				context.Background(),
				rKey,
				*s.ShortUrl,
				realUrl,
			)
		}
	}

	return err
}
