package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/keyword-search-backend/api_error"
	configure "github.com/miniyus/keyword-search-backend/config"
	"github.com/miniyus/keyword-search-backend/database"
	"gorm.io/gorm"
	"log"
	"strconv"
)

type RouterGroup func(router Router, app Application)

type MiddlewareRegister func(fiber *fiber.App, app Application)

type Application interface {
	IsProduction() bool
	Middleware(fn MiddlewareRegister)
	Route(prefix string, fn RouterGroup, name ...string)
	Fiber() *fiber.App
	Config() *configure.Configs
	DB() *gorm.DB
	Stats()
	Run()
}

type app struct {
	fiber  *fiber.App
	config *configure.Configs
	db     *gorm.DB
}

func (a *app) Fiber() *fiber.App {
	return a.fiber
}

func (a *app) Config() *configure.Configs {
	return a.config
}

func (a *app) DB() *gorm.DB {
	return a.db
}

func (a *app) Middleware(fn MiddlewareRegister) {
	fn(a.Fiber(), a)
}

func (a *app) Route(prefix string, fn RouterGroup, name ...string) {
	r := NewRouter(a.Fiber(), prefix, name...)

	fn(r, a)
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
	log.Printf("Handlers Count: %d", a.Fiber().HandlersCount())
	log.Println("[Router]")

	for _, r := range a.Fiber().GetRoutes() {
		log.Printf(
			"[%s] '%s' | '%s' , Params: %s",
			r.Method, r.Name, r.Path, r.Params,
		)
	}

}
func (a *app) Run() {
	port := a.config.AppPort
	err := a.fiber.Listen(":" + strconv.Itoa(port))

	if err != nil {
		log.Fatalf("error start fiber app: %v", err)
	}
}

func (a *app) IsProduction() bool {
	return a.Config().AppEnv == configure.PRD
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
		fiber.New(config.App),
		config,
		database.DB(config.Database),
	}
}
