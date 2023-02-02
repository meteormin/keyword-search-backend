package IOContainer

import (
	"fmt"
	"reflect"
)

// Container
// 일종의 IoC 컨테이너, fiber.App 객체를 담고 있다.
type Container interface {
	Instances() map[reflect.Type]interface{}
	Singleton(instance interface{})
	Bind(keyType interface{}, resolver interface{})
	Resolve(keyType interface{}) interface{}
}

// ContainerStruct
// Container 구현체
type ContainerStruct struct {
	instances map[reflect.Type]interface{}
	bindings  map[reflect.Type]interface{}
}

// NewContainer
// IoC 컨테이너 생성 함수
func NewContainer() Container {
	return &ContainerStruct{
		instances: make(map[reflect.Type]interface{}),
	}
}

// Singleton
// 특정 객체를 singleton 패턴으로 컨테이너에 저장하는 메서드
// 클로저를 받을 수도 있다.
func (w *ContainerStruct) Singleton(instance interface{}) {
	reflectInstanceType := reflect.TypeOf(instance)
	if reflectInstanceType.Kind() == reflect.Func {
		instance = w.call(instance)
		reflectInstanceType = reflect.TypeOf(instance)
	}

	w.instances[reflectInstanceType] = instance
}

// Instances
// 저장된 singleton 객체 슬라이스를 리턴한다.
func (w *ContainerStruct) Instances() map[reflect.Type]interface{} {
	return w.instances
}

// Bind
// 구조체 혹은 인터페이스의 타입을 키 값으로 저장한다.
// keyType 파라미터는 주소 값을 전달 해야 한다.
// resolver 파라미터는 콜백 함수(클로저)를 통해 특정 인터 페이스와 구현체를 매치할 수 있다.
func (w *ContainerStruct) Bind(keyType interface{}, resolver interface{}) {
	reflectResolver := reflect.TypeOf(resolver)
	reflectKeyType := reflect.TypeOf(keyType)
	if reflectResolver.Kind() == reflect.Func {
		w.instances[reflectKeyType.Elem()] = resolver
		return
	}

	panic("Can not Bind...")
}

// Resolve get or create instance in container
func (w *ContainerStruct) Resolve(keyType interface{}) interface{} {
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
func (w *ContainerStruct) call(callable interface{}) interface{} {
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
