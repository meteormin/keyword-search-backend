package app

import (
	"github.com/gofiber/fiber/v2"
	"strings"
)

// SubRouter
// sub router 등록을 위해 일관성을 위해 생성한 타입
type SubRouter func(router fiber.Router)

// Router
// Routes wrapper
type Router interface {
	App() *fiber.App
	Route(prefix string, callback SubRouter, middleware ...fiber.Handler) fiber.Router
	GetRoutes() []fiber.Router
}

// routerStruct
// route wrapper struct
type routerStruct struct {
	app        *fiber.App
	router     fiber.Router
	name       string
	routes     []fiber.Router
	GroupCount int
}

// NewRouter
// 라우터 생성
func NewRouter(app *fiber.App, prefix string, name ...string) Router {
	routeName := ""
	if len(name) > 0 {
		routeName = name[0]
	}

	router := app.Group(prefix).Name(routeName)

	return &routerStruct{
		app:        app,
		router:     router,
		name:       routeName,
		routes:     make([]fiber.Router, 0),
		GroupCount: 1,
	}
}

// App get fiber app
func (r *routerStruct) App() *fiber.App {
	return r.app
}

// Route
// route 등록 메서드
func (r *routerStruct) Route(prefix string, callback SubRouter, middleware ...fiber.Handler) fiber.Router {
	grp := r.router.Group(prefix, middleware...)
	callback(grp)

	r.GroupCount += 1
	r.routes = append(r.routes, grp)

	name := strings.Replace(prefix, "/", ".", -1)
	if !strings.HasPrefix(name, ".") {
		name = "." + name
	}

	return grp.Name(name)
}

// GetRoutes
// 등록한 route slice
func (r *routerStruct) GetRoutes() []fiber.Router {
	return r.routes
}
