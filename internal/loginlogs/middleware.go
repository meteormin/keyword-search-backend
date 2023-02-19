package loginlogs

import (
	"github.com/gofiber/fiber/v2"
	mEntity "github.com/miniyus/gofiber/entity"
	"github.com/miniyus/gofiber/pkg/gormhooks"
	"github.com/miniyus/keyword-search-backend/entity"
	"gorm.io/gorm"
	"strings"
	"time"
)

func Middleware(db *gorm.DB, method string, path string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		if strings.ToUpper(method) != strings.ToUpper(ctx.Method()) {
			return ctx.Next()
		}

		if strings.ToLower(path) != strings.ToLower(ctx.Path()) {
			return ctx.Next()
		}

		repo := NewRepository(db)

		at := mEntity.AccessToken{}
		atHooks := gormhooks.New(&at)
		atHooks.HandleAfterCreate(func(c *mEntity.AccessToken, tx *gorm.DB) (err error) {
			log := entity.LoginLog{
				UserId:  c.UserId,
				LoginAt: time.Now(),
				Ip:      ctx.IP(),
			}

			_, err = repo.Create(log)
			return err
		})

		return ctx.Next()
	}
}
