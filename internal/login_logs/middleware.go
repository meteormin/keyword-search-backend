package login_logs

import (
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/gofiber/auth"
	"github.com/miniyus/keyword-search-backend/entity"
	"gorm.io/gorm"
	"time"
)

func Middleware(db *gorm.DB) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		repo := NewRepository(db)
		user, err := auth.GetAuthUser(ctx)
		if err != nil {
			return err
		}

		log := entity.LoginLog{
			UserId:  user.Id,
			LoginAt: time.Now(),
			Ip:      ctx.IP(),
		}

		_, err = repo.Create(log)
		if err != nil {
			return err
		}

		return ctx.Next()
	}
}
