package routes

import (
	"github.com/miniyus/go-fiber/internal/app/api_auth"
	"github.com/miniyus/go-fiber/internal/app/hosts"
	"github.com/miniyus/go-fiber/internal/app/users"
	"github.com/miniyus/go-fiber/internal/core/container"
	"github.com/miniyus/go-fiber/internal/core/context"
	"github.com/miniyus/go-fiber/pkg/jwt"
	"go.uber.org/zap"
)

const Prefix = "/api"

func SetRoutes(container container.Container) {
	api := container.App().Group(Prefix)

	var tokenGenerator jwt.Generator

	container.Resolve(&tokenGenerator)
	logger := container.Get(context.Logger).(*zap.SugaredLogger)

	api_auth.Register(api, api_auth.New(container.Database(), tokenGenerator, logger))

	users.Register(api, users.New(container.Database()))

	hosts.Register(api, hosts.New(container.Database()))
}
