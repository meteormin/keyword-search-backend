package IOContainer_test

import (
	"github.com/miniyus/keyword-search-backend/pkg/IOContainer"
	"log"
	"testing"
)

type AnyType struct {
	Name string
}

var con = IOContainer.NewContainer()

func TestContainerStruct_Bind(t *testing.T) {
	var testData AnyType

	con.Bind(&testData, func() AnyType {
		return AnyType{Name: "test"}
	})
	var resolveData AnyType
	con.Resolve(&resolveData)

	if resolveData.Name != "test" {
		t.Error("error...")
	}

	resolveData.Name = "test2"
}

func TestContainerStruct_Singleton(t *testing.T) {
	testData := &AnyType{Name: "test2"}

	var resolveData *AnyType

	con.Singleton(func() *AnyType {
		return testData
	})

	con.Resolve(&resolveData)
	log.Println(con.Instances())
	if resolveData.Name != "test2" {
		log.Print(resolveData)
		t.Error("error...")
	}
}
