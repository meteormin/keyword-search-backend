package utils

import (
	"github.com/miniyus/keyword-search-backend/pkg/slice"
)

type Collection[T interface{}] interface {
	Items() []T
	Count() int
	Add(item T)
	Map(fn func(v T, i int) T) Collection[T]
	Filter(fn func(v T, i int) bool) Collection[T]
	Except(fn func(v T, i int) bool) Collection[T]
	Chunk(chunkSize int, fn func(v []T, i int)) [][]T
	For(fn func(v T, i int))
	Remove(index int)
	Concat(items []T)
}

type BaseCollection[T interface{}] struct {
	items []T
}

func NewCollection[T interface{}](items []T) Collection[T] {
	return &BaseCollection[T]{
		items: items,
	}
}

func (b *BaseCollection[T]) Items() []T {
	return b.items
}

func (b *BaseCollection[T]) Count() int {
	return len(b.items)
}

func (b *BaseCollection[T]) Add(item T) {
	b.items = slice.Add(b.items, item)
}

func (b *BaseCollection[T]) Map(fn func(v T, i int) T) Collection[T] {
	items := slice.Map(b.items, fn)
	return NewCollection(items)
}

func (b *BaseCollection[T]) Filter(fn func(v T, i int) bool) Collection[T] {
	filtered := slice.Filter(b.items, fn)
	return NewCollection(filtered)
}

func (b *BaseCollection[T]) Except(fn func(v T, i int) bool) Collection[T] {
	excepts := slice.Except(b.items, fn)
	return NewCollection(excepts)
}

func (b *BaseCollection[T]) Chunk(chunkSize int, fn func(v []T, i int)) [][]T {
	return slice.Chunk(b.items, chunkSize, fn)
}

func (b *BaseCollection[T]) For(fn func(v T, i int)) {
	slice.For(b.items, fn)
}

func (b *BaseCollection[T]) Remove(index int) {
	b.items = slice.Remove(b.items, index)
}

func (b *BaseCollection[T]) Concat(items []T) {
	b.items = slice.Concat(b.items, items)
}
