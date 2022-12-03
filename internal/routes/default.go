package routes

import (
	"github.com/miniyus/go-fiber/internal/app/api_auth"
	"github.com/miniyus/go-fiber/internal/app/hosts"
	"github.com/miniyus/go-fiber/internal/app/users"
	"github.com/miniyus/go-fiber/internal/core/container"
	"github.com/miniyus/go-fiber/pkg/jwt"
)

const Prefix = "/api"

func SetRoutes(container container.Container) {
	api := container.App().Group(Prefix)

	var tokenGenerator jwt.Generator

	container.Resolve(&tokenGenerator)
	api_auth.Register(api, api_auth.New(container.Database(), tokenGenerator))

	// need jwt token
	users.Register(api, users.New(container.Database()))
	hosts.Register(api, hosts.New(container.Database()))
}
