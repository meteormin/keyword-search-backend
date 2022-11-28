package routes

import (
	"github.com/miniyus/go-fiber/internal/app/api_auth"
	"github.com/miniyus/go-fiber/internal/app/test_api"
	"github.com/miniyus/go-fiber/internal/app/users"
	"github.com/miniyus/go-fiber/internal/core/container"
	"github.com/miniyus/go-fiber/pkg/jwt"
)

const Prefix = "/api"

func SetRoutes(container container.Container) {
	api := container.App().Group(Prefix)

	var tokenGenerator jwt.Generator

	container.Resolve(&tokenGenerator)
	api_auth.Register(api, api_auth.Factory(container.Database(), tokenGenerator))

	test_api.Register(api, test_api.Factory(container.Database()))

	// need jwt token

	users.Register(api, users.Factory(container.Database()))

}
