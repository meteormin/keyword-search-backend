package utils_test

import (
	"github.com/miniyus/keyword-search-backend/utils"
	"log"
	"testing"
)

var testData = []int{
	1, 2, 3,
}

func TestBaseCollection_Items(t *testing.T) {
	var collection = utils.NewCollection(testData)

	items := collection.Items()

	for i, n := range items {
		log.Print(i, n)
		if n != testData[i] {
			t.Errorf("not match! %d:%d", i, n)
		}
	}
}

func TestBaseCollection_Add(t *testing.T) {
	var collection = utils.NewCollection(testData)

	collection.Add(4)
	i := collection.Items()

	log.Print(i[len(i)-1])

	if 4 != i[len(i)-1] {
		t.Error("result must be 4")
	}

}

func TestBaseCollection_Chunk(t *testing.T) {
	var collection = utils.NewCollection(testData)

	chunked := collection.Chunk(1, func(n []int, i int) {
		log.Print(n, i)
	})

	if len(chunked) != len(testData) {
		t.Error("failed chunk!")
	}
}

func TestBaseCollection_Concat(t *testing.T) {
	var collection = utils.NewCollection(testData)
	collection.Concat([]int{4, 5, 6})

	resultData := []int{1, 2, 3, 4, 5, 6}

	for i, n := range collection.Items() {
		log.Print(i, n)
		if n != resultData[i] {
			t.Errorf("not match!! %d:%d", i, n)
		}
	}
}

func TestBaseCollection_Except(t *testing.T) {
	var collection = utils.NewCollection(testData)
	result := collection.Except(func(v int, i int) bool {
		return v == 1
	})

	for _, n := range result {
		if n == 1 {
			t.Error("FAIL!")
		}
	}
}

func TestBaseCollection_Filter(t *testing.T) {
	var collection = utils.NewCollection(testData)
	resultCollection := collection.Filter(func(v int, i int) bool {
		return v == 1
	})

	for _, n := range resultCollection {
		if n != 1 {
			t.Error("FAIL!")
		}
	}
}

func TestBaseCollection_For(t *testing.T) {
	var collection = utils.NewCollection(testData)
	collection.For(func(v int, i int) {
		if v != testData[i] {
			t.Error("FAIL!")
		}
	})
}

func TestBaseCollection_Map(t *testing.T) {
	var collection = utils.NewCollection(testData)
	result := collection.Map(func(v int, i int) int {
		return i + 1
	})

	for i, n := range result {
		log.Print(i, n)
		if n != testData[i] {
			t.Error("Fail")
		}
	}
}

func TestBaseCollection_Remove(t *testing.T) {
	var collection = utils.NewCollection(testData)
	collection.Remove(0)
	collection.For(func(v int, i int) {
		log.Print(i, v)
		if v == 1 {
			t.Error("not removed")
		}
	})
}
