package register_test

import (
	"github.com/gofiber/fiber/v2"
	configure "github.com/miniyus/keyword-search-backend/config"
	"github.com/miniyus/keyword-search-backend/internal/database"
	"github.com/miniyus/keyword-search-backend/pkg/IOContainer"
	"github.com/miniyus/keyword-search-backend/pkg/IOContainer/register"
	"github.com/miniyus/keyword-search-backend/pkg/worker"
	"log"
	"testing"
)

func TestRegister(t *testing.T) {
	cfg := configure.GetConfigs()
	c := IOContainer.NewContainer(fiber.New(), database.DB(cfg.Database), cfg)

	register.Resister(c)

	var jobDispatcher worker.Dispatcher

	c.Resolve(&jobDispatcher)

	log.Print(jobDispatcher)
}
