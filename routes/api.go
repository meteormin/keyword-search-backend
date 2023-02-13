package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/gofiber/app"
	"github.com/miniyus/gofiber/auth"
	configure "github.com/miniyus/gofiber/config"
	"github.com/miniyus/gofiber/database"
	"github.com/miniyus/gofiber/jobs"
	"github.com/miniyus/gofiber/log"
	"github.com/miniyus/gofiber/permission"
	"github.com/miniyus/gofiber/pkg/jwt"
	rsGen "github.com/miniyus/gofiber/pkg/rs256"
	"github.com/miniyus/gofiber/pkg/worker"
	"github.com/miniyus/gofiber/utils"
	"github.com/miniyus/keyword-search-backend/internal/host_search"
	"github.com/miniyus/keyword-search-backend/internal/hosts"
	"github.com/miniyus/keyword-search-backend/internal/login_logs"
	"github.com/miniyus/keyword-search-backend/internal/search"
	"github.com/miniyus/keyword-search-backend/internal/short_url"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"path"
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

	var zLogger *zap.SugaredLogger
	a.Resolve(&zLogger)

	if zLogger == nil {
		zLogger = log.GetLogger()
	}

	authMiddlewaresParameter := auth.MiddlewaresParameter{
		Cfg:    cfg.Auth.Jwt,
		Logger: zLogger,
	}

	apiRouter.Route("/auth", func(router fiber.Router) {
		privateKey := rsGen.PrivatePemDecode(path.Join(cfg.Path.DataPath, "secret/private.pem"))
		tokenGenerator := jwt.NewGenerator(privateKey, privateKey.Public(), cfg.Auth.Exp)
		authHandler := auth.New(db, tokenGenerator)

		router.Post("/token", login_logs.Middleware(db), authHandler.SignIn).Name("auth.token")
	}).Name("api.auth")

	hasPermission := permission.HasPermission(permission.HasPermissionParameter{
		DB:           db,
		DefaultPerms: cfg.Permission,
		FilterFunc:   nil,
	})

	apiRouter.Route(
		hosts.Prefix,
		hosts.Register(hosts.New(db)),
		auth.Middlewares(authMiddlewaresParameter, hasPermission())...,
	).Name("api.hosts")

	apiRouter.Route(
		search.Prefix,
		search.Register(search.New(db)),
		auth.Middlewares(authMiddlewaresParameter, hasPermission())...,
	).Name("api.search")

	hostSearchHandler := host_search.New(db, jDispatcher)

	apiRouter.Route(
		host_search.Prefix,
		host_search.Register(hostSearchHandler),
		auth.Middlewares(
			authMiddlewaresParameter,
			jobs.AddJobMeta(jDispatcher, db),
			hasPermission(),
		)...,
	).Name("api.hosts.search")

	apiRouter.Route(
		short_url.Prefix,
		short_url.Register(short_url.New(
			db,
			utils.RedisClientMaker(cfg.RedisConfig),
		)),
		auth.Middlewares(authMiddlewaresParameter, hasPermission())...,
	).Name("api.short_url")

}
