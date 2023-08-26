package files_test

import (
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/gofiber"
	"github.com/miniyus/gofiber/app"
	"github.com/miniyus/gofiber/database"
	"github.com/miniyus/keyword-search-backend/config"
	"github.com/miniyus/keyword-search-backend/entity"
	"github.com/miniyus/keyword-search-backend/internal/files"
	gormLogger "gorm.io/gorm/logger"
	"io"
	"net/http/httptest"
	"testing"
	"time"
)

func init() {
	cfg := config.GetConfigs()
	cfg.Database = map[string]database.Config{
		"default": {
			Name:        "default",
			Driver:      "postgres",
			Host:        "localhost",
			Dbname:      "go_fiber",
			Username:    "smyoo",
			Password:    "smyoo",
			Port:        "5432",
			TimeZone:    "Asia/Seoul",
			SSLMode:     false,
			AutoMigrate: []interface{}{&entity.File{}},
			Logger: gormLogger.Config{
				SlowThreshold:             time.Second,
				LogLevel:                  gormLogger.Silent,
				IgnoreRecordNotFoundError: true,
				Colorful:                  true,
			},
			MaxIdleConn: 10,
			MaxOpenConn: 100,
			MaxLifeTime: time.Hour,
		},
	}
	gofiber.New(*cfg.Configs)
}

func TestNew(t *testing.T) {
	a := gofiber.App()
	a.Route("/api", func(router app.Router, app app.Application) {
		db := gofiber.DB("default")
		f := files.New(db)
		router.Route(files.Prefix, f)
	})

	a.Bootstrap()
}

func TestHandler_GetFiles(t *testing.T) {
	req := httptest.NewRequest(fiber.MethodGet, "/api/files", nil)
	test, err := app.App().Test(req)
	if err != nil {
		t.Error(err)
	}
	all, err := io.ReadAll(test.Body)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(all))
}
