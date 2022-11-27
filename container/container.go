package container

import (
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/go-fiber/config"
	"github.com/miniyus/go-fiber/database"
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
	InjectMap(instances map[string]interface{})
	Inject(key string, instance interface{})
	Get(key string) interface{}
	Bind(keyType interface{}, resolver interface{})
	Resolve(keyType interface{}) interface{}
}

type Wrapper struct {
	app       *fiber.App
	config    *config.Configs
	database  *gorm.DB
	instances map[string]interface{}
	bindings  map[reflect.Type]interface{}
}

func NewContainer(app *fiber.App, config *config.Configs) *Wrapper {
	return &Wrapper{
		app:       app,
		config:    config,
		database:  database.DB(config.Database),
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

func (w *Wrapper) InjectMap(instances map[string]interface{}) {
	w.instances = instances
}

func (w *Wrapper) Inject(key string, instance interface{}) {
	w.instances[key] = instance
}

func (w *Wrapper) Instances() map[string]interface{} {
	return w.instances
}

func (w *Wrapper) Bind(keyType interface{}, resolver interface{}) {
	resolverType := reflect.TypeOf(resolver)
	if resolverType.Kind() == reflect.Func {
		reflectedFunction := reflect.TypeOf(resolver)
		argumentsCount := reflectedFunction.NumIn()
		arguments := make([]reflect.Value, argumentsCount)
		values := reflect.ValueOf(resolver).Call(arguments)

		w.bindings[reflect.TypeOf(keyType)] = values[0].Interface()
	} else {
		w.bindings[reflect.TypeOf(keyType)] = resolver
	}
}

func (w *Wrapper) Resolve(keyType interface{}) interface{} {
	receiverType := reflect.TypeOf(keyType)

	if receiverType.Kind() == reflect.Ptr {
		bind := w.bindings[reflect.TypeOf(keyType)]
		reflect.ValueOf(keyType).Elem().Set(reflect.ValueOf(bind))

		return bind
	}

	return nil
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
