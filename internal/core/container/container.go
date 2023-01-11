package container

import (
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/keyword-search-backend/config"
	"github.com/miniyus/keyword-search-backend/internal/core/context"
	"gorm.io/gorm"
	"log"
	"reflect"
	"strconv"
)

// Container
// 일종의 IoC 컨테이너, fiber.App 객체를 담고 있다.
type Container interface {
	App() *fiber.App
	Config() *config.Configs
	Database() *gorm.DB
	Instances() map[context.Key]interface{}
	Bindings() map[reflect.Type]interface{}
	Singleton(key context.Key, instance interface{})
	Get(key context.Key) interface{}
	Bind(keyType interface{}, resolver interface{})
	Resolve(keyType interface{}) interface{}
	Run()
	Stats()
}

// Wrapper
// Container 구현체
type Wrapper struct {
	app       *fiber.App
	database  *gorm.DB
	config    *config.Configs
	instances map[context.Key]interface{}
	bindings  map[reflect.Type]interface{}
}

// New
// IoC 컨테이너 새성 함수
func New(app *fiber.App, db *gorm.DB, config *config.Configs) Container {
	return &Wrapper{
		app:       app,
		database:  db,
		config:    config,
		instances: make(map[context.Key]interface{}),
		bindings:  make(map[reflect.Type]interface{}),
	}
}

// App
// fiber.App 객체 getter 역할
func (w *Wrapper) App() *fiber.App {
	return w.app
}

// Config
// config.Configs getter
func (w *Wrapper) Config() *config.Configs {
	return w.config
}

// Database
// *gorm.DB getter
func (w *Wrapper) Database() *gorm.DB {
	return w.database
}

// Singleton
// 특정 객체를 singleton 패턴으로 컨테이너에 저장하는 메서드
func (w *Wrapper) Singleton(key context.Key, instance interface{}) {
	w.instances[key] = instance
}

// Instances
// 저장된 singleton 객체 슬라이스를 리턴한다.
func (w *Wrapper) Instances() map[context.Key]interface{} {
	return w.instances
}

// Bind
// 구조체 혹은 인터페이스를 키 값으로 저장할 수 있음
// resolver 파라미터는 콜백 함수 혹은 객체를 직접 바인딩하여, 특정 인터 페이스와 구현체를 매치할 수 있다.
func (w *Wrapper) Bind(keyType interface{}, resolver interface{}) {
	w.bindings[reflect.TypeOf(keyType)] = resolver
}

// Bindings
// 저장된 bind 객체 슬라이스를 리턴한다.
func (w *Wrapper) Bindings() map[reflect.Type]interface{} {
	return w.bindings
}

// Resolve
// Bind 메서드에서 바인딩한 구조체 혹은 인터페이스를 가지고 객체를 생성하여 주거나 저장된 객체를 리턴 해준다.
// Bind 메서드에서 resolver 파라미터가 함수일 경우, 해당 함수를 실행 시켜 새로운 객체를 생성하여 리턴 한다.
// Bind 메서드에서 resolver 파라미터가 구조체의 포인터일 경우 reflect를 활용하여 객체를 생성하여 리턴한다.
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

// call
// Bind 의 resolver 파라미터가 함수인 경우 reflect를 활용하여 함수를 실행 시켜 준다.
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

// Get
// Singleton 으로 주입한 객체를 주입했을 때 사용한 key 값으로 가져온다.
func (w *Wrapper) Get(key context.Key) interface{} {
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

// Stats
// Debug 용도,
// 현재 생성된 route list
// 컨테이너가 가지고 있는 정보 콘솔 로그로 보여준다.
func (w *Wrapper) Stats() {
	if w.Config().AppEnv == config.PRD {
		log.Printf("'AppEnv' is %s", config.PRD)
		return
	}
	log.Println("[Container Info]")
	log.Printf("ENV: %s", w.Config().AppEnv)
	log.Printf("Locale: %s", w.Config().Locale)
	log.Printf("Time Zone: %s", w.Config().TimeZone)
	log.Printf("Injected Instances: %#v", w.Instances())
	log.Printf("Bindings: %#v", w.Bindings())

	log.Println("[Fiber App Info]")
	log.Printf("Handlers Count: %d", w.App().HandlersCount())
	log.Println("[Router]")
	for _, r := range w.App().GetRoutes() {
		log.Printf(
			"[%s] '%s' | '%s' , Params: %s",
			r.Method, r.Name, r.Path, r.Params,
		)
	}

}
