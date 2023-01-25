package routes

import (
	"github.com/miniyus/keyword-search-backend/app"
	"github.com/miniyus/keyword-search-backend/auth"
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
	"github.com/miniyus/keyword-search-backend/resolver"
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
		FilterFunc: group_detail.FilterFunc(group_detail.FilterParameter{
			DB: a.DB(),
		}),
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

	hostSearchHandler := host_search.New(a.DB(), zapLogger())
	hostSearchParameter := host_search.RegisterParameter{
		HasPerm:       hasPermParam,
		JobDispatcher: jDispatcher(),
	}

	apiRouter.Route(
		host_search.Prefix,
		host_search.Register(hostSearchHandler, hostSearchParameter),
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

	apiRouter.Route(
		test_api.Prefix,
		test_api.Register(jDispatcher(), zapLogger()),
	).Name("api.test_api")
}
