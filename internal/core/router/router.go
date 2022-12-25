package router

import (
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type Register = func(router fiber.Router)

type Router interface {
	Route(prefix string, callback Register, middleware ...fiber.Handler) fiber.Router
	GetRoutes() []fiber.Router
}

type Wrapper struct {
	Router     fiber.Router
	name       string
	GroupCount int
	routes     []fiber.Router
}

func New(router fiber.Router, name ...string) Router {
	routeName := ""
	if len(name) > 0 {
		routeName = name[0]
	}

	router.Name(routeName)

	return &Wrapper{
		Router:     router,
		name:       routeName,
		GroupCount: 1,
		routes:     make([]fiber.Router, 0),
	}
}

func (r *Wrapper) Route(prefix string, callback Register, middleware ...fiber.Handler) fiber.Router {
	grp := r.Router.Group(prefix, middleware...)
	callback(grp)

	r.GroupCount += 1
	r.routes = append(r.routes, grp)

	return grp.Name(r.name + strconv.Itoa(r.GroupCount))
}

func (r *Wrapper) GetRoutes() []fiber.Router {
	return r.routes
}
