package routes

import (
	"github.com/miniyus/go-fiber/internal/app/api_auth"
	"github.com/miniyus/go-fiber/internal/app/host_search"
	"github.com/miniyus/go-fiber/internal/app/hosts"
	"github.com/miniyus/go-fiber/internal/app/search"
	"github.com/miniyus/go-fiber/internal/app/short_url"
	"github.com/miniyus/go-fiber/internal/app/users"
	"github.com/miniyus/go-fiber/internal/core/auth"
	"github.com/miniyus/go-fiber/internal/core/container"
	"github.com/miniyus/go-fiber/internal/core/router"
	"github.com/miniyus/go-fiber/pkg/jwt"
	"go.uber.org/zap"
)

const ApiPrefix = "/api"

func Api(container container.Container) {
	apiGroup := container.App().Group(ApiPrefix)

	apiRouter := router.New(apiGroup, "api")

	var tokenGenerator jwt.Generator
	container.Resolve(&tokenGenerator)

	var logger *zap.SugaredLogger
	container.Resolve(&logger)

	apiRouter.Route(
		api_auth.Prefix,
		api_auth.Register(api_auth.New(
			container.Database(),
			tokenGenerator,
			logger,
		)),
	).Name("api.auth")

	apiRouter.Route(
		users.Prefix,
		users.Register(users.New(container.Database(), logger)),
		auth.Middlewares()...,
	).Name("api.users")

	apiRouter.Route(
		hosts.Prefix,
		hosts.Register(hosts.New(container.Database(), logger)),
		auth.Middlewares()...,
	).Name("api.hosts")

	apiRouter.Route(
		search.Prefix,
		search.Register(search.New(container.Database(), logger)),
		auth.Middlewares()...,
	).Name("api.search")

	apiRouter.Route(
		host_search.Prefix,
		host_search.Register(host_search.New(container.Database(), logger)),
		auth.Middlewares()...,
	).Name("api.hosts.search")

	apiRouter.Route(
		short_url.Prefix,
		short_url.Register(short_url.New(container.Database(), logger)),
		auth.Middlewares()...,
	).Name("api.short_url")

}
