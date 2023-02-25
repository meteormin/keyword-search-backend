package entity

import (
	"context"
	"fmt"
	"github.com/miniyus/gofiber/app"
	"github.com/miniyus/gorm-extension/gormhooks"
	worker "github.com/miniyus/goworker"
	"github.com/miniyus/keyword-search-backend/utils"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"strconv"
)

type Host struct {
	gorm.Model
	UserId      uint      `json:"user_id"`
	User        *User     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"user"`
	Host        string    `gorm:"column:host;type:varchar(100);uniqueIndex" json:"host"`
	Subject     string    `gorm:"column:subject;type:varchar(100)" json:"subject"`
	Description string    `gorm:"column:description;type:varchar(255)" json:"description"`
	Path        string    `gorm:"column:path;type:varchar(255)" json:"path"`
	Publish     bool      `gorm:"column:publish;type:bool" json:"publish"`
	Search      []*Search `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"search"`
}

func (h *Host) Hooks() *gormhooks.Hooks[*Host] {
	return gormhooks.GetHooks(h)
}

func (h *Host) AfterSave(tx *gorm.DB) (err error) {
	return h.Hooks().AfterSave(tx)
}

type hostHookHandler struct {
	app app.Application
}

func newHostHookHandler(a app.Application) *hostHookHandler {
	return &hostHookHandler{app: a}
}

func (handler hostHookHandler) AfterSave(h *Host, tx *gorm.DB) (err error) {
	a := handler.app
	var rClient *redis.Client
	a.Resolve(&rClient)

	var dispatcher worker.Dispatcher
	a.Resolve(&dispatcher)

	jobId := fmt.Sprintf("%s.%d.%d", "hosts", h.ID, h.UserId)

	var search []Search
	err = tx.Where(&Search{HostId: h.ID}).Find(&search).Error
	if err != nil {
		return err
	}

	err = dispatcher.Dispatch(jobId, func(j *worker.Job) error {
		j.Meta["user_id"] = h.UserId
		for _, s := range search {
			if s.ShortUrl != nil {
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
		}

		return nil
	})

	return err
}
