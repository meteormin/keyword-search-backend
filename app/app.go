package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/keyword-search-backend/pkg/IOContainer"
	"log"
	"reflect"
	"strconv"
)

type Env string

const (
	PRD   Env = "production"
	DEV   Env = "development"
	LOCAL Env = "local"
)

type Config struct {
	Env         Env
	Port        int
	Locale      string
	TimeZone    string
	FiberConfig fiber.Config
}

var defaultConfig = Config{
	Env:         LOCAL,
	Port:        8000,
	Locale:      "",
	TimeZone:    "Asia/Seoul",
	FiberConfig: fiber.Config{},
}

type Register func(app Application)

type RouterGroup func(router Router, app Application)

type RegisterFiber func(fiber *fiber.App)

type MiddlewareRegister func(fiber *fiber.App, app Application)

type Application interface {
	IOContainer.Container
	IsProduction() bool
	Middleware(fn MiddlewareRegister)
	Route(prefix string, fn RouterGroup, name ...string)
	Status()
	Run()
	Register(fn Register)
	RegisterFiber(fn RegisterFiber)
}

type app struct {
	IOContainer.Container
	fiber  *fiber.App
	config Config
	isRun  bool
}

// New
// fiber app wrapper
func New(cfgs ...Config) Application {
	var fiberConfig fiber.Config
	var cfg Config

	if len(cfgs) == 0 {
		cfg = defaultConfig
		fiberConfig = cfg.FiberConfig
	} else {
		cfg = cfgs[0]
		fiberConfig = cfg.FiberConfig
	}

	return &app{
		Container: IOContainer.NewContainer(),
		config:    cfg,
		fiber:     fiber.New(fiberConfig),
		isRun:     false,
	}
}

func (a *app) Register(fn Register) {
	fn(a)
}

// RegisterFiber
// fiber application을 클로저를 통해 제어
func (a *app) RegisterFiber(fn RegisterFiber) {
	if fn == nil {
		return
	}

	fn(a.fiber)
}

// Middleware
// add middleware from closure
func (a *app) Middleware(fn MiddlewareRegister) {
	a.Register(func(this Application) {
		if _, ok := this.(*app); ok {
			fn(this.(*app).fiber, this)
		}

	})
}

// Route
// register route group
func (a *app) Route(prefix string, fn RouterGroup, name ...string) {
	a.Register(func(this Application) {
		if _, ok := this.(*app); ok {
			r := NewRouter(this.(*app).fiber, prefix, name...)
			fn(r, a)
		}
	})
}

// Status
// Debug 용도,
// 현재 생성된 route list
// 컨테이너가 가지고 있는 정보 콘솔 로그로 보여준다.
func (a *app) Status() {
	if a.IsProduction() {
		log.Printf("'AppEnv' is %s", PRD)
		return
	}

	log.Println("[Container Info]")
	log.Printf("ENV: %s", a.config.Env)
	log.Printf("Locale: %s", a.config.Locale)
	log.Printf("Time Zone: %s", a.config.TimeZone)

	log.Println("[Fiber App Info]")
	log.Printf("Handlers Count: %d", a.fiber.HandlersCount())
	log.Println("[Router]")

	for _, r := range a.fiber.GetRoutes() {
		log.Printf(
			"[%s] '%s' | '%s' , Params: %s",
			r.Method, r.Name, r.Path, r.Params,
		)
	}

}

// Run
// run fiber application
func (a *app) Run() {
	if a.isRun {
		return
	}

	port := a.config.Port
	err := a.fiber.Listen(":" + strconv.Itoa(port))

	if err != nil {
		log.Fatalf("error start fiber app: %v", err)
	}

	a.isRun = true
}

func (a *app) IsProduction() bool {
	return a.config.Env == PRD
}

func (a *app) Instances() map[reflect.Type]interface{} {
	return a.Container.Instances()
}

func (a *app) Bind(keyType interface{}, resolver interface{}) {
	a.Container.Bind(keyType, resolver)
}

func (a *app) Resolve(resolver interface{}) interface{} {
	return a.Container.Resolve(resolver)
}

func (a *app) Singleton(instance interface{}) {
	a.Container.Singleton(instance)
}
