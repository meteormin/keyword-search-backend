package entity

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v9"
	"github.com/miniyus/gofiber/app"
	"gorm.io/gorm"
	"path"
	"strconv"
	"strings"
)

type Search struct {
	gorm.Model
	Host        *Host   `json:"host" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	HostId      uint    `json:"host_id" gorm:"index:query_unique,unique"`
	QueryKey    string  `gorm:"column:query_key;type:varchar(50);index:query_unique,unique" json:"query_key"`
	Query       string  `gorm:"column:query;type:varchar(50);index:query_unique,unique" json:"query"`
	Description string  `gorm:"column:description;type:varchar(50)" json:"description"`
	Publish     bool    `gorm:"column:publish;type:bool" json:"publish"`
	ShortUrl    *string `gorm:"column:short_url;type:varchar(255);uniqueIndex" json:"short_url"`
}

func (s *Search) AfterSave(tx *gorm.DB) (err error) {
	a := app.App()

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

	return err
}
