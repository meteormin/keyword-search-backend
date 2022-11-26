package router

import (
	"github.com/miniyus/go-fiber/container"
	"github.com/miniyus/go-fiber/internal/test_api"
	"github.com/miniyus/go-fiber/internal/users"
)

func SetRoutes(container container.Container) {
	api := container.App().Group("/api")

	test_api.SetRoutes(api, test_api.Factory(container.Database()))

	users.SetRoutes(api, users.Factory(container.Database()))
}
