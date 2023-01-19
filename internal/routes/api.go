package routes

import (
	"github.com/miniyus/keyword-search-backend/internal/api/api_auth"
	"github.com/miniyus/keyword-search-backend/internal/api/groups"
	"github.com/miniyus/keyword-search-backend/internal/api/host_search"
	"github.com/miniyus/keyword-search-backend/internal/api/hosts"
	"github.com/miniyus/keyword-search-backend/internal/api/search"
	"github.com/miniyus/keyword-search-backend/internal/api/short_url"
	"github.com/miniyus/keyword-search-backend/internal/api/users"
	"github.com/miniyus/keyword-search-backend/internal/core/auth"
	"github.com/miniyus/keyword-search-backend/internal/core/container"
	"github.com/miniyus/keyword-search-backend/internal/core/permission"
	"github.com/miniyus/keyword-search-backend/internal/core/register/resolver"
	"github.com/miniyus/keyword-search-backend/internal/core/router"
	"github.com/miniyus/keyword-search-backend/pkg/jwt"
	"github.com/miniyus/keyword-search-backend/pkg/worker"
	"go.uber.org/zap"
)

const ApiPrefix = "/api"

func Api(c container.Container) {
	var jobDispatcher worker.Dispatcher
	c.Resolve(&jobDispatcher)

	var zapLogger *zap.SugaredLogger
	c.Resolve(&zapLogger)

	var tokenGenerator jwt.Generator
	c.Resolve(&tokenGenerator)

	apiRouter := router.New(c.App(), ApiPrefix, "api")

	apiRouter.Route(
		api_auth.Prefix,
		api_auth.Register(api_auth.New(
			c.Database(),
			tokenGenerator,
			zapLogger,
		)),
	).Name("api.auth")

	apiRouter.Route(
		groups.Prefix,
		groups.Register(groups.New(c.Database(), zapLogger)),
		auth.Middlewares(permission.HasPermission())...,
	).Name("api.groups")

	apiRouter.Route(
		users.Prefix,
		users.Register(users.New(c.Database(), zapLogger)),
		auth.Middlewares()...,
	).Name("api.users")

	apiRouter.Route(
		hosts.Prefix,
		hosts.Register(hosts.New(c.Database(), zapLogger)),
		auth.Middlewares()...,
	).Name("api.hosts")

	apiRouter.Route(
		search.Prefix,
		search.Register(search.New(c.Database(), zapLogger)),
		auth.Middlewares()...,
	).Name("api.search")

	apiRouter.Route(
		host_search.Prefix,
		host_search.Register(host_search.New(c.Database(), zapLogger)),
		auth.Middlewares()...,
	).Name("api.hosts.search")

	apiRouter.Route(
		short_url.Prefix,
		short_url.Register(short_url.New(
			c.Database(),
			resolver.MakeRedisClient(c),
			zapLogger,
		)),
		auth.Middlewares()...,
	).Name("api.short_url")

}
