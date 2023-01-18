package container

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/keyword-search-backend/config"
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
	IsProduction() bool
	Instances() map[reflect.Type]interface{}
	Singleton(instance interface{})
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
	instances map[reflect.Type]interface{}
	bindings  map[reflect.Type]interface{}
}

// New
// IoC 컨테이너 생성 함수
func New(app *fiber.App, db *gorm.DB, config *config.Configs) Container {
	return &Wrapper{
		app:       app,
		database:  db,
		config:    config,
		instances: make(map[reflect.Type]interface{}),
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

// IsProduction
// check app env is production
func (w *Wrapper) IsProduction() bool {
	if w.Config().AppEnv == config.PRD {
		return true
	}

	return false
}

// Singleton
// 특정 객체를 singleton 패턴으로 컨테이너에 저장하는 메서드
// 클로저를 받을 수도 있다.
func (w *Wrapper) Singleton(instance interface{}) {
	reflectInstanceType := reflect.TypeOf(instance)
	if reflectInstanceType.Kind() == reflect.Func {
		instance = w.call(instance)
		reflectInstanceType = reflect.TypeOf(instance)
	}

	w.instances[reflectInstanceType] = instance
}

// Instances
// 저장된 singleton 객체 슬라이스를 리턴한다.
func (w *Wrapper) Instances() map[reflect.Type]interface{} {
	return w.instances
}

// Bind
// 구조체 혹은 인터페이스의 타입을 키 값으로 저장한다.
// keyType 파라미터는 주소 값을 전달 해야 한다.
// resolver 파라미터는 콜백 함수(클로저)를 통해 특정 인터 페이스와 구현체를 매치할 수 있다.
func (w *Wrapper) Bind(keyType interface{}, resolver interface{}) {
	reflectResolver := reflect.TypeOf(resolver)
	reflectKeyType := reflect.TypeOf(keyType)
	if reflectResolver.Kind() == reflect.Func {
		w.instances[reflectKeyType.Elem()] = resolver
		return
	}

	panic("Can not Bind...")
}

// Resolve get or create instance in container
func (w *Wrapper) Resolve(keyType interface{}) interface{} {
	receiverType := reflect.TypeOf(keyType)

	receiver, exists := w.instances[receiverType.Elem()]
	if !exists {
		panic(fmt.Sprintf("Can not find %v", receiverType.Elem()))
	}

	if reflect.TypeOf(receiver).Kind() == reflect.Func {
		receiver = w.call(receiver)
	}

	if receiverType.Kind() == reflect.Ptr {
		reflect.ValueOf(keyType).Elem().Set(reflect.ValueOf(receiver))
	}

	return receiver
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
	if w.IsProduction() {
		log.Printf("'AppEnv' is %s", config.PRD)
		return
	}

	log.Println("[Container Info]")
	log.Printf("ENV: %s", w.Config().AppEnv)
	log.Printf("Locale: %s", w.Config().Locale)
	log.Printf("Time Zone: %s", w.Config().TimeZone)
	log.Printf("Injected Instances: %#v", w.Instances())

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
