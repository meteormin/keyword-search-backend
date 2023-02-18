package routes

import (
	"github.com/miniyus/gofiber/app"
	"github.com/miniyus/gofiber/auth"
	configure "github.com/miniyus/gofiber/config"
	"github.com/miniyus/gofiber/database"
	"github.com/miniyus/gofiber/jobqueue"
	"github.com/miniyus/gofiber/jobs"
	"github.com/miniyus/gofiber/log"
	"github.com/miniyus/gofiber/permission"
	"github.com/miniyus/gofiber/pkg/worker"
	"github.com/miniyus/gofiber/utils"
	"github.com/miniyus/keyword-search-backend/internal/host_search"
	"github.com/miniyus/keyword-search-backend/internal/hosts"
	"github.com/miniyus/keyword-search-backend/internal/search"
	"github.com/miniyus/keyword-search-backend/internal/short_url"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const ApiPrefix = "/api"

func Api(apiRouter app.Router, a app.Application) {
	var cfg *configure.Configs
	a.Resolve(&cfg)

	if cfg == nil {
		configs := configure.GetConfigs()
		cfg = &configs
	}

	var db *gorm.DB
	a.Resolve(&db)

	if db == nil {
		db = database.GetDB()
	}

	var jDispatcher worker.Dispatcher
	jDispatcher = jobqueue.GetDispatcher()

	if jDispatcher == nil {
		a.Resolve(&jDispatcher)
	}

	var zLogger *zap.SugaredLogger
	a.Resolve(&zLogger)

	if zLogger == nil {
		zLogger = log.GetLogger()
	}

	hasPermission := permission.HasPermission(permission.HasPermissionParameter{
		DB:           db,
		DefaultPerms: cfg.Permission,
		FilterFunc:   nil,
	})

	apiRouter.Route(
		hosts.Prefix,
		hosts.Register(hosts.New(db)),
		auth.JwtMiddleware(cfg.Auth.Jwt), auth.Middlewares(), hasPermission(),
	).Name("api.hosts")

	apiRouter.Route(
		search.Prefix,
		search.Register(search.New(db)),
		auth.JwtMiddleware(cfg.Auth.Jwt), auth.Middlewares(), hasPermission(),
	).Name("api.search")

	hostSearchHandler := host_search.New(db, jDispatcher)

	apiRouter.Route(
		host_search.Prefix,
		host_search.Register(hostSearchHandler),
		auth.JwtMiddleware(cfg.Auth.Jwt), auth.Middlewares(),
		jobs.AddJobMeta(),
		hasPermission(),
	).Name("api.hosts.search")

	apiRouter.Route(
		short_url.Prefix,
		short_url.Register(short_url.New(
			db,
			utils.RedisClientMaker(cfg.RedisConfig),
		)),
		auth.JwtMiddleware(cfg.Auth.Jwt), auth.Middlewares(), hasPermission(),
	).Name("api.short_url")

}
