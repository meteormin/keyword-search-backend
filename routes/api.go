package routes

import (
	"github.com/miniyus/keyword-search-backend/app"
	"github.com/miniyus/keyword-search-backend/auth"
	configure "github.com/miniyus/keyword-search-backend/config"
	"github.com/miniyus/keyword-search-backend/internal/api_auth"
	"github.com/miniyus/keyword-search-backend/internal/group_detail"
	"github.com/miniyus/keyword-search-backend/internal/groups"
	"github.com/miniyus/keyword-search-backend/internal/host_search"
	"github.com/miniyus/keyword-search-backend/internal/hosts"
	"github.com/miniyus/keyword-search-backend/internal/search"
	"github.com/miniyus/keyword-search-backend/internal/short_url"
	"github.com/miniyus/keyword-search-backend/internal/test_api"
	"github.com/miniyus/keyword-search-backend/internal/users"
	"github.com/miniyus/keyword-search-backend/permission"
	"github.com/miniyus/keyword-search-backend/pkg/jwt"
	rsGen "github.com/miniyus/keyword-search-backend/pkg/rs256"
	"github.com/miniyus/keyword-search-backend/pkg/worker"
	"github.com/miniyus/keyword-search-backend/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"path"
)

const ApiPrefix = "/api"

func Api(apiRouter app.Router, a app.Application) {
	var cfg *configure.Configs
	a.Resolve(&cfg)

	var zapLogger *zap.SugaredLogger
	a.Resolve(&zapLogger)

	var db *gorm.DB
	a.Resolve(&db)

	privateKey := rsGen.PrivatePemDecode(path.Join(cfg.Path.DataPath, "secret/private.pem"))
	tokenGenerator := jwt.NewGenerator(privateKey, privateKey.Public(), cfg.Auth.Exp)

	permissions := permission.NewPermissionsFromConfig(cfg.Permission)
	permissionCollection := permission.NewPermissionCollection(permissions...)

	opts := cfg.QueueConfig
	opts.Redis = utils.RedisClientMaker(cfg.RedisConfig)

	jDispatcher := worker.NewDispatcher(opts)
	jDispatcher.Run()

	authMiddlewareParam := auth.MiddlewaresParameter{
		Cfg: cfg.Auth.Jwt,
		DB:  db,
	}

	hasPermParam := permission.HasPermissionParameter{
		DB:           db,
		DefaultPerms: permissionCollection,
		FilterFunc: group_detail.FilterFunc(group_detail.FilterParameter{
			DB: db,
		}),
	}

	apiRouter.Route(
		api_auth.Prefix,
		api_auth.Register(
			api_auth.New(
				db,
				tokenGenerator,
				zapLogger,
			),
			authMiddlewareParam,
		),
	).Name("api.auth")

	apiRouter.Route(
		groups.Prefix,
		groups.Register(groups.New(db, zapLogger)),
		auth.Middlewares(authMiddlewareParam, permission.HasPermission(hasPermParam))...,
	).Name("api.groups")

	apiRouter.Route(
		users.Prefix,
		users.Register(users.New(db, zapLogger)),
		auth.Middlewares(authMiddlewareParam)...,
	).Name("api.users")

	apiRouter.Route(
		hosts.Prefix,
		hosts.Register(hosts.New(db, zapLogger)),
		auth.Middlewares(authMiddlewareParam)...,
	).Name("api.hosts")

	apiRouter.Route(
		search.Prefix,
		search.Register(search.New(db, zapLogger)),
		auth.Middlewares(authMiddlewareParam)...,
	).Name("api.search")

	hostSearchHandler := host_search.New(db, zapLogger, jDispatcher)

	apiRouter.Route(
		host_search.Prefix,
		host_search.Register(hostSearchHandler, hasPermParam),
		auth.Middlewares(authMiddlewareParam)...,
	).Name("api.hosts.search")

	apiRouter.Route(
		short_url.Prefix,
		short_url.Register(short_url.New(
			db,
			utils.RedisClientMaker(cfg.RedisConfig),
			zapLogger,
		)),
		auth.Middlewares(authMiddlewareParam)...,
	).Name("api.short_url")

	apiRouter.Route(
		test_api.Prefix,
		test_api.Register(jDispatcher, zapLogger),
	).Name("api.test_api")
}
