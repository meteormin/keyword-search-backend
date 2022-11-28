package routes

import (
	api_auth2 "github.com/miniyus/go-fiber/internal/app/api_auth"
	test_api2 "github.com/miniyus/go-fiber/internal/app/test_api"
	users2 "github.com/miniyus/go-fiber/internal/app/users"
	"github.com/miniyus/go-fiber/internal/core/container"
	"github.com/miniyus/go-fiber/pkg/jwt"
)

var Prefix = "/api"

func SetRoutes(container container.Container) {
	api := container.App().Group(Prefix)

	var tokenGenerator jwt.Generator

	container.Resolve(&tokenGenerator)
	api_auth2.Register(api, api_auth2.Factory(container.Database(), tokenGenerator))

	test_api2.Register(api, test_api2.Factory(container.Database()))

	// need jwt token

	users2.Register(api, users2.Factory(container.Database()))

}
