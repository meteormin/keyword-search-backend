package container

import (
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/go-fiber/config"
	"github.com/miniyus/go-fiber/database"
	"gorm.io/gorm"
	"log"
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
}

type Wrapper struct {
	app       *fiber.App
	config    *config.Configs
	database  *gorm.DB
	instances map[string]interface{}
}

func NewContainer(app *fiber.App, config *config.Configs) *Wrapper {
	return &Wrapper{
		app:       app,
		config:    config,
		database:  database.DB(config.Database),
		instances: make(map[string]interface{}),
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
