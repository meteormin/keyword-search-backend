package routes

import (
	"github.com/miniyus/keyword-search-backend/internal/app/api_auth"
	"github.com/miniyus/keyword-search-backend/internal/app/groups"
	"github.com/miniyus/keyword-search-backend/internal/app/host_search"
	"github.com/miniyus/keyword-search-backend/internal/app/hosts"
	"github.com/miniyus/keyword-search-backend/internal/app/search"
	"github.com/miniyus/keyword-search-backend/internal/app/short_url"
	"github.com/miniyus/keyword-search-backend/internal/app/users"
	"github.com/miniyus/keyword-search-backend/internal/core/auth"
	"github.com/miniyus/keyword-search-backend/internal/core/container"
	"github.com/miniyus/keyword-search-backend/internal/core/permission"
	"github.com/miniyus/keyword-search-backend/internal/core/register/resolver"
	"github.com/miniyus/keyword-search-backend/internal/core/router"
)

const ApiPrefix = "/api"

func Api(c container.Container) {
	apiRouter := router.New(c.App(), ApiPrefix, "api")

	apiRouter.Route(
		api_auth.Prefix,
		api_auth.Register(api_auth.New(
			c.Database(),
			resolver.TokenGenerator(c),
			resolver.Logger(c),
		)),
	).Name("api.auth")

	apiRouter.Route(
		groups.Prefix,
		groups.Register(groups.New(c.Database(), resolver.Logger(c))),
		auth.Middlewares(permission.HasPermission())...,
	).Name("api.groups")

	apiRouter.Route(
		users.Prefix,
		users.Register(users.New(c.Database(), resolver.Logger(c))),
		auth.Middlewares()...,
	).Name("api.users")

	apiRouter.Route(
		hosts.Prefix,
		hosts.Register(hosts.New(c.Database(), resolver.Logger(c))),
		auth.Middlewares()...,
	).Name("api.hosts")

	apiRouter.Route(
		search.Prefix,
		search.Register(search.New(c.Database(), resolver.Logger(c))),
		auth.Middlewares()...,
	).Name("api.search")

	apiRouter.Route(
		host_search.Prefix,
		host_search.Register(host_search.New(c.Database(), resolver.Logger(c))),
		auth.Middlewares()...,
	).Name("api.hosts.search")

	apiRouter.Route(
		short_url.Prefix,
		short_url.Register(short_url.New(c.Database(), resolver.Logger(c))),
		auth.Middlewares()...,
	).Name("api.short_url")

}
