package entity

import (
	"context"
	"fmt"
	"github.com/miniyus/gofiber/app"
	"github.com/miniyus/gofiber/pkg/worker"
	"github.com/miniyus/gofiber/utils"
	"github.com/miniyus/keyword-search-backend/config"
	"gorm.io/gorm"
	"path"
	"strconv"
	"strings"
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

func (h *Host) AfterSave(tx *gorm.DB) (err error) {
	a := app.App()

	var dispatcher worker.Dispatcher
	a.Resolve(&dispatcher)

	rClientFn := utils.RedisClientMaker(config.GetConfigs().RedisConfig)
	jobId := fmt.Sprintf("%s.%d.%d", "hosts", h.ID, h.UserId)

	err = dispatcher.Dispatch(jobId, func(j *worker.Job) error {
		rClient := rClientFn()

		for _, s := range h.Search {
			if s.ShortUrl != nil {
				rKey := "short_url." + strconv.Itoa(int(h.UserId))
				cached, err := rClient.HGet(
					context.Background(),
					rKey,
					*s.ShortUrl,
				).Result()

				if cached != "" && err == nil {
					sep := ":/"
					splitString := strings.Split(h.Host, sep)
					hostPath := path.Join(splitString[1], h.Path)
					queryKey := s.QueryKey
					queryString := s.Query

					realUrl := fmt.Sprintf(
						"%s:/%s?%s=%s",
						splitString[0], hostPath, queryKey, queryString,
					)

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
