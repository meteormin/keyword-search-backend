package routes

import (
	"github.com/miniyus/go-fiber/container"
	"github.com/miniyus/go-fiber/internal/auth"
	"github.com/miniyus/go-fiber/internal/test_api"
	"github.com/miniyus/go-fiber/internal/users"
	"github.com/miniyus/go-fiber/pkg/jwt"
)

var Prefix = "/api"

func SetRoutes(container container.Container) {
	api := container.App().Group(Prefix)

	test_api.Register(api, test_api.Factory(container.Database()))

	users.Register(api, users.Factory(container.Database()))

	var tokenGenerator jwt.Generator

	container.Resolve(&tokenGenerator)

	auth.Register(api, auth.Factory(container.Database(), tokenGenerator))
}
