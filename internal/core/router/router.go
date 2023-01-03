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
	App() *fiber.App
	Route(prefix string, callback Register, middleware ...fiber.Handler) fiber.Router
	GetRoutes() []*fiber.Router
}

// Wrapper
// route wrapper struct
type Wrapper struct {
	app        *fiber.App
	router     fiber.Router
	name       string
	routes     []*fiber.Router
	GroupCount int
}

// New
// 라우터 생성
func New(app *fiber.App, prefix string, name ...string) Router {
	routeName := ""
	if len(name) > 0 {
		routeName = name[0]
	}

	router := app.Group(prefix).Name(routeName)

	return &Wrapper{
		app:        app,
		router:     router,
		name:       routeName,
		routes:     make([]*fiber.Router, 0),
		GroupCount: 1,
	}
}

// App get fiber app
func (r *Wrapper) App() *fiber.App {
	return r.app
}

// Route
// route 등록 메서드
func (r *Wrapper) Route(prefix string, callback Register, middleware ...fiber.Handler) fiber.Router {
	grp := r.router.Group(prefix, middleware...)
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
