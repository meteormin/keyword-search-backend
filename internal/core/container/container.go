package container

import (
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/keyword-search-backend/config"
	"gorm.io/gorm"
	"log"
	"reflect"
	"strconv"
)

type Container interface {
	App() *fiber.App
	Config() *config.Configs
	Database() *gorm.DB
	Instances() map[string]interface{}
	Bindings() map[reflect.Type]interface{}
	Singleton(key string, instance interface{})
	Get(key string) interface{}
	Bind(keyType interface{}, resolver interface{})
	Resolve(keyType interface{}) interface{}
	Run()
	Stats()
}

type Wrapper struct {
	app       *fiber.App
	database  *gorm.DB
	config    *config.Configs
	instances map[string]interface{}
	bindings  map[reflect.Type]interface{}
}

func New(app *fiber.App, db *gorm.DB, config *config.Configs) Container {
	return &Wrapper{
		app:       app,
		database:  db,
		config:    config,
		instances: make(map[string]interface{}),
		bindings:  make(map[reflect.Type]interface{}),
	}
}

func (w *Wrapper) App() *fiber.App {
	return w.app
}

func (w *Wrapper) Config() *config.Configs {
	return w.config
}

func (w *Wrapper) Database() *gorm.DB {
	return w.database
}

func (w *Wrapper) Singleton(key string, instance interface{}) {
	w.instances[key] = instance
}

func (w *Wrapper) Instances() map[string]interface{} {
	return w.instances
}

func (w *Wrapper) Bindings() map[reflect.Type]interface{} {
	return w.bindings
}

func (w *Wrapper) Bind(keyType interface{}, resolver interface{}) {
	w.bindings[reflect.TypeOf(keyType)] = resolver
}

func (w *Wrapper) Resolve(keyType interface{}) interface{} {
	receiverType := reflect.TypeOf(keyType)

	bind := w.bindings[reflect.TypeOf(keyType)]

	if reflect.TypeOf(bind).Kind() == reflect.Func {
		bind = w.call(bind)
	}

	if receiverType.Kind() == reflect.Ptr {
		reflect.ValueOf(keyType).Elem().Set(reflect.ValueOf(bind))
	}

	return bind
}

func (w *Wrapper) call(callable interface{}) interface{} {
	resolverType := reflect.TypeOf(callable)
	if resolverType.Kind() == reflect.Func {
		reflectedFunction := reflect.TypeOf(callable)
		argumentsCount := reflectedFunction.NumIn()
		arguments := make([]reflect.Value, argumentsCount)
		values := reflect.ValueOf(callable).Call(arguments)

		return values[0].Interface()
	} else {
		return callable
	}
}

func (w *Wrapper) Get(key string) interface{} {
	return w.instances[key]
}

// Run fiber app
func (w *Wrapper) Run() {
	port := w.config.AppPort
	err := w.app.Listen(":" + strconv.Itoa(port))

	if err != nil {
		log.Fatalf("error start fiber app: %v", err)
	}

}

func (w *Wrapper) Stats() {
	if w.Config().AppEnv == config.PRD {
		log.Printf("'AppEnv' is %s", config.PRD)
		return
	}
	log.Println("[Container Info]")
	log.Printf("[ENV] %s", w.Config().AppEnv)
	log.Printf("[Locale] %s", w.Config().Locale)
	log.Printf("[Time Zone] %s", w.Config().TimeZone)
	log.Printf("[Injected Instances] %#v", w.Instances())
	log.Printf("[Bindings] %#v", w.Bindings())
	log.Println("...")
	log.Println("[Fiber App Info]")
	log.Printf("[Handlers Count] %d", w.App().HandlersCount())
	log.Println("[Router]")
	for _, r := range w.App().GetRoutes() {
		log.Printf(
			"[%s] '%s' | '%s' , Params: %s",
			r.Method, r.Name, r.Path, r.Params,
		)
	}

}
