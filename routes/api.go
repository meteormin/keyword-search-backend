package routes

import (
	"github.com/miniyus/gofiber/app"
	configure "github.com/miniyus/gofiber/config"
	"github.com/miniyus/gofiber/database"
	"github.com/miniyus/gofiber/permission"
	"github.com/miniyus/gofiber/pkg/worker"
	"github.com/miniyus/gofiber/utils"
	"github.com/miniyus/keyword-search-backend/internal/group_detail"
	"github.com/miniyus/keyword-search-backend/internal/host_search"
	"github.com/miniyus/keyword-search-backend/internal/hosts"
	"github.com/miniyus/keyword-search-backend/internal/search"
	"github.com/miniyus/keyword-search-backend/internal/short_url"
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
	a.Resolve(&jDispatcher)

	permissions := permission.NewPermissionsFromConfig(cfg.Permission)
	permissionCollection := permission.NewPermissionCollection(permissions...)

	hasPermParam := permission.HasPermissionParameter{
		DB:           db,
		DefaultPerms: permissionCollection,
		FilterFunc:   group_detail.FilterFunc,
	}

	apiRouter.Route(
		hosts.Prefix,
		hosts.Register(hosts.New(db)),
	).Name("api.hosts")

	apiRouter.Route(
		search.Prefix,
		search.Register(search.New(db)),
	).Name("api.search")

	hostSearchHandler := host_search.New(db, jDispatcher)

	apiRouter.Route(
		host_search.Prefix,
		host_search.Register(hostSearchHandler, hasPermParam),
	).Name("api.hosts.search")

	apiRouter.Route(
		short_url.Prefix,
		short_url.Register(short_url.New(
			db,
			utils.RedisClientMaker(cfg.RedisConfig),
		)),
	).Name("api.short_url")

}
