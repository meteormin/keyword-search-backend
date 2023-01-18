package register_test

import (
	"github.com/go-redis/redis/v9"
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/keyword-search-backend/config"
	"github.com/miniyus/keyword-search-backend/internal/core"
	"github.com/miniyus/keyword-search-backend/internal/core/permission"
	"github.com/miniyus/keyword-search-backend/pkg/jwt"
	"github.com/miniyus/keyword-search-backend/pkg/worker"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"log"
	"testing"
)

func TestRegister(t *testing.T) {
	c := core.New()

	log.Print(c.Instances())

	var app *fiber.App
	c.Resolve(&app)
	log.Print(app)

	var cfg *config.Configs
	c.Resolve(&cfg)
	log.Print(cfg)

	var db *gorm.DB
	c.Resolve(&db)
	log.Print(db)

	var rc *redis.Client
	c.Resolve(&rc)
	log.Print(rc)

	var tg jwt.Generator
	c.Resolve(&tg)
	log.Print(tg)

	var logs *zap.SugaredLogger
	c.Resolve(&logs)
	logs.Info("logs")

	var permCollect permission.Collection
	c.Resolve(&permCollect)
	log.Print(permCollect)

	var jobDispatcher worker.Dispatcher
	c.Resolve(&jobDispatcher)
	log.Print(jobDispatcher)

}
