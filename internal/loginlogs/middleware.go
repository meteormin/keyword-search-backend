package loginlogs

import (
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/gorm-extension/gormhooks"
	"github.com/miniyus/keyword-search-backend/entity"
	"github.com/miniyus/keyword-search-backend/repo"
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

		repository := repo.NewLoginLogRepository(db)

		at := entity.AccessToken{}
		atHooks := gormhooks.New(&at)
		atHooks.HandleAfterCreate(func(c *entity.AccessToken, tx *gorm.DB) (err error) {
			log := entity.LoginLog{
				UserId:  c.UserId,
				LoginAt: time.Now(),
				Ip:      ctx.IP(),
			}

			_, err = repository.Create(log)
			return err
		})

		return ctx.Next()
	}
}
