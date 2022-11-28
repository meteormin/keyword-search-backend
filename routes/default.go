package routes

import (
	"github.com/miniyus/go-fiber/core/auth"
	"github.com/miniyus/go-fiber/core/container"
	"github.com/miniyus/go-fiber/internal/api_auth"
	"github.com/miniyus/go-fiber/internal/test_api"
	"github.com/miniyus/go-fiber/internal/users"
	"github.com/miniyus/go-fiber/pkg/jwt"
)

var Prefix = "/api"

func SetRoutes(container container.Container) {
	api := container.App().Group(Prefix)

	var tokenGenerator jwt.Generator

	container.Resolve(&tokenGenerator)
	api_auth.Register(api, api_auth.Factory(container.Database(), tokenGenerator))

	test_api.Register(api, test_api.Factory(container.Database()))

	// need jwt token
	container.App().Use(auth.JwtMiddleware(container.Config().Auth.Jwt))

	users.Register(api, users.Factory(container.Database()))

}
