package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/keyword-search-backend/api_error"
	configure "github.com/miniyus/keyword-search-backend/config"
	"github.com/miniyus/keyword-search-backend/database"
	"github.com/miniyus/keyword-search-backend/pkg/IOContainer"
	"gorm.io/gorm"
	"log"
	"reflect"
	"strconv"
)

type Config struct {
	FiberConfig fiber.Config
}

type Register func(app Application)

type RouterGroup func(router Router, app Application)

type RegisterFiber func(fiber *fiber.App)

type MiddlewareRegister func(fiber *fiber.App, app Application)

type RegisterContainer func(container IOContainer.Container)

type Application interface {
	IOContainer.Container
	IsProduction() bool
	register(fn Register)
	Middleware(fn MiddlewareRegister)
	Route(prefix string, fn RouterGroup, name ...string)
	Stats()
	Run()
	RegisterContainer(fn RegisterContainer)
	RegisterFiber(fn RegisterFiber)
}

type app struct {
	IOContainer.Container
	fiber  *fiber.App
	config *configure.Configs
	db     *gorm.DB
	isRun  bool
}

// New
// fiber app wrapper
func New(configs ...*configure.Configs) Application {
	var config *configure.Configs

	if len(configs) == 0 {
		config = configure.GetConfigs()
	} else {
		config = configs[0]
	}

	fiberConfig := config.App
	dbConfig := config.Database

	fiberConfig.ErrorHandler = api_error.OverrideDefaultErrorHandler(config)
	if fiber.IsChild() {
		dbConfig.AutoMigrate = false
	}

	return &app{
		Container: IOContainer.NewContainer(),
		fiber:     fiber.New(fiberConfig),
		config:    config,
		db:        database.DB(dbConfig),
		isRun:     false,
	}
}

// RegisterFiber
// fiber application을 클로저를 통해 제어
func (a *app) RegisterFiber(fn RegisterFiber) {
	if fn == nil {
		return
	}

	fn(a.fiber)
}

// RegisterContainer
// 컨테이너 구조체를 클로저를 통해 container 제어를 할 수 있는 함수
func (a *app) RegisterContainer(fn RegisterContainer) {
	if fn == nil {
		return
	}

	fn(a.Container)
}

// Config
// get configuration
func (a *app) Config() *configure.Configs {
	return a.config
}

// DB
// get database connection
func (a *app) DB() *gorm.DB {
	return a.db
}

func (a *app) register(fn Register) {
	fn(a)
}

// Middleware
// add middleware from closure
func (a *app) Middleware(fn MiddlewareRegister) {
	a.register(func(this Application) {
		if _, ok := this.(*app); ok {
			fn(this.(*app).fiber, this)
		}

	})
}

// Route
// register route group
func (a *app) Route(prefix string, fn RouterGroup, name ...string) {
	a.register(func(this Application) {
		if _, ok := this.(*app); ok {
			r := NewRouter(this.(*app).fiber, prefix, name...)
			fn(r, a)
		}
	})
}

// Stats
// Debug 용도,
// 현재 생성된 route list
// 컨테이너가 가지고 있는 정보 콘솔 로그로 보여준다.
func (a *app) Stats() {
	if a.IsProduction() {
		log.Printf("'AppEnv' is %s", configure.PRD)
		return
	}

	log.Println("[Container Info]")
	log.Printf("ENV: %s", a.Config().AppEnv)
	log.Printf("Locale: %s", a.Config().Locale)
	log.Printf("Time Zone: %s", a.Config().TimeZone)

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

	port := a.config.AppPort
	err := a.fiber.Listen(":" + strconv.Itoa(port))

	if err != nil {
		log.Fatalf("error start fiber app: %v", err)
	}

	a.isRun = true
}

func (a *app) IsProduction() bool {
	return a.config.AppEnv == configure.PRD
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
