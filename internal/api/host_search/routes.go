package host_search

import (
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/keyword-search-backend/app"
	"github.com/miniyus/keyword-search-backend/config"
	"github.com/miniyus/keyword-search-backend/permission"
	"github.com/miniyus/keyword-search-backend/pkg/worker"
)

const Prefix = "/hosts/:id/search"
const jobDispatcher = config.JobDispatcher

type RegisterParameter struct {
	HasPerm       permission.HasPermissionParameter
	JobDispatcher worker.Dispatcher
}

func Register(handler Handler, parameter RegisterParameter) app.SubRouter {
	return func(router fiber.Router) {
		router.Get("/", handler.GetByHostId).Name("api.hosts.id.search")
		router.Get("/descriptions", handler.GetDescriptionsByHostId).Name("api.hosts.id.search.descriptions")
		router.Post("/",
			permission.HasPermission(parameter.HasPerm),
			config.AddContext(jobDispatcher, parameter.JobDispatcher),
			handler.BatchCreate,
		).Name("api.hosts.batch-create")
	}
}
