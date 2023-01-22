package routes

import (
	"github.com/miniyus/keyword-search-backend/internal/api/api_auth"
	"github.com/miniyus/keyword-search-backend/internal/api/groups"
	"github.com/miniyus/keyword-search-backend/internal/api/host_search"
	"github.com/miniyus/keyword-search-backend/internal/api/hosts"
	"github.com/miniyus/keyword-search-backend/internal/api/search"
	"github.com/miniyus/keyword-search-backend/internal/api/short_url"
	"github.com/miniyus/keyword-search-backend/internal/api/users"
	"github.com/miniyus/keyword-search-backend/internal/app"
	"github.com/miniyus/keyword-search-backend/internal/auth"
	"github.com/miniyus/keyword-search-backend/internal/permission"
	"github.com/miniyus/keyword-search-backend/internal/resolver"
)

const ApiPrefix = "/api"

func Api(apiRouter app.Router, a app.Application) {
	zapLogger := resolver.MakeLogger(a.Config().CustomLogger)
	tokenGenerator := resolver.MakeJwtGenerator(resolver.JwtGeneratorConfig{
		DataPath: a.Config().Path.DataPath,
		Exp:      a.Config().Auth.Exp,
	})
	permissionCollection := resolver.MakePermissionCollection(a.Config().Permission)
	jDispatcher := resolver.MakeJobDispatcher(resolver.JobDispatcherConfig{
		RedisCfg:  a.Config().RedisConfig,
		WorkerCfg: a.Config().QueueConfig,
	})

	authMiddlewareParam := auth.MiddlewaresParameter{
		Cfg: a.Config().Auth.Jwt,
		DB:  a.DB(),
	}

	hasPermParam := permission.HasPermissionParameter{
		DB:           a.DB(),
		DefaultPerms: permissionCollection(),
	}

	apiRouter.Route(
		api_auth.Prefix,
		api_auth.Register(
			api_auth.New(
				a.DB(),
				tokenGenerator(),
				zapLogger(),
			),
			authMiddlewareParam,
		),
	).Name("api.auth")

	apiRouter.Route(
		groups.Prefix,
		groups.Register(groups.New(a.DB(), zapLogger())),
		auth.Middlewares(authMiddlewareParam, permission.HasPermission(hasPermParam))...,
	).Name("api.groups")

	apiRouter.Route(
		users.Prefix,
		users.Register(users.New(a.DB(), zapLogger())),
		auth.Middlewares(authMiddlewareParam)...,
	).Name("api.users")

	apiRouter.Route(
		hosts.Prefix,
		hosts.Register(hosts.New(a.DB(), zapLogger())),
		auth.Middlewares(authMiddlewareParam)...,
	).Name("api.hosts")

	apiRouter.Route(
		search.Prefix,
		search.Register(search.New(a.DB(), zapLogger())),
		auth.Middlewares(authMiddlewareParam)...,
	).Name("api.search")

	hostSearch := host_search.New(a.DB(), zapLogger())
	hostSearchParameter := host_search.RegisterParameter{
		HasPerm:       hasPermParam,
		JobDispatcher: jDispatcher(),
	}
	apiRouter.Route(
		host_search.Prefix,
		host_search.Register(hostSearch, hostSearchParameter),
		auth.Middlewares(authMiddlewareParam)...,
	).Name("api.hosts.search")

	apiRouter.Route(
		short_url.Prefix,
		short_url.Register(short_url.New(
			a.DB(),
			resolver.MakeRedisClient(a.Config().RedisConfig),
			zapLogger(),
		)),
		auth.Middlewares(authMiddlewareParam)...,
	).Name("api.short_url")

}
