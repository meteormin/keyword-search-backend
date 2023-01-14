package slice_test

import (
	"github.com/miniyus/keyword-search-backend/pkg/slice"
	"log"
	"testing"
)

func TestChunk(t *testing.T) {
	testData := make([]int, 12)
	slice.Chunk(testData, 2, func(v []int, i int) {
		if len(v) != 2 {
			t.Error(len(v), i)
		}
	})
}

func TestConcat(t *testing.T) {
	testData := make([]int, 3)
	rs := slice.Concat(testData, []int{4, 5})
	if len(rs) != 5 {
		t.Error(len(rs), rs)
	}
}

func TestFor(t *testing.T) {
	testData := make([]int, 10)
	slice.For(testData, func(v int, i int) {
		if i >= 10 {
			t.Error(i)
		}
	})
}

func TestMap(t *testing.T) {
	testData := make([]int, 10)
	rs := slice.Map(testData, func(v int, i int) int {
		return i + 1
	})

	log.Print(rs)
}

func TestExcept(t *testing.T) {
	testData := []int{1, 2, 3, 4, 5}
	rs := slice.Except(testData, func(v int, i int) bool {
		return v == 1
	})

	if rs[0] == 1 {
		t.Error(rs[0])
	}
}

func TestFilter(t *testing.T) {
	testData := []int{1, 2, 3, 4, 5}
	rs := slice.Filter(testData, func(v int, i int) bool {
		return v == 1
	})

	if rs[0] != 1 {
		t.Error(rs[0])
	}
}

func TestAdd(t *testing.T) {
	testData := []int{1, 2, 3, 4}
	rs := slice.Add(testData, 5)

	if rs[len(rs)-1] == 4 {
		t.Error(rs[len(rs)-1])
	}
}

func TestRemove(t *testing.T) {
	testData := []int{1, 2, 3, 4, 5}
	rs := slice.Remove(testData, 0)

	if rs[0] == 1 {
		t.Error(rs[0])
	}
}
