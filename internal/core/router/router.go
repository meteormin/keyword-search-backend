package router

import (
	"github.com/gofiber/fiber/v2"
)

// Register
// sub router 등록을 위해 일관성을 위해 생성한 타입
type Register = func(router fiber.Router)

// Router
// Routes wrapper
type Router interface {
	Route(prefix string, callback Register, middleware ...fiber.Handler) fiber.Router
	GetRoutes() []*fiber.Router
}

// Wrapper
// route wrapper struct
type Wrapper struct {
	Router     fiber.Router
	name       string
	GroupCount int
	routes     []*fiber.Router
}

// New
// 라우터 생성
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
		routes:     make([]*fiber.Router, 0),
	}
}

// Route
// route 등록 메서드
func (r *Wrapper) Route(prefix string, callback Register, middleware ...fiber.Handler) fiber.Router {
	grp := r.Router.Group(prefix, middleware...)
	callback(grp)

	r.GroupCount += 1
	r.routes = append(r.routes, &grp)

	return grp.Name(r.name + "." + prefix)
}

// GetRoutes
// 등록한 route slice
func (r *Wrapper) GetRoutes() []*fiber.Router {
	return r.routes
}
