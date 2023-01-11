package utils

import "github.com/miniyus/keyword-search-backend/pkg"

type Collection[T interface{}] interface {
	Items() []T
	Add(item T)
	Map(fn func(v T, i int) T) []T
	Filter(fn func(v T, i int) bool) []T
	Except(fn func(v T, i int) bool) []T
	Chunk(chunkSize int, fn func(v []T, i int)) []T
	For(fn func(v T, i int)) []T
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

func (b *BaseCollection[T]) Add(item T) {
	b.items = append(b.items, item)
}

func (b *BaseCollection[T]) Map(fn func(v T, i int) T) []T {
	return pkg.Map(b.items, fn)
}

func (b *BaseCollection[T]) Filter(fn func(v T, i int) bool) []T {
	return pkg.Filter(b.items, fn)
}

func (b *BaseCollection[T]) Except(fn func(v T, i int) bool) []T {
	return pkg.Except(b.items, fn)
}

func (b *BaseCollection[T]) Chunk(chunkSize int, fn func(v []T, i int)) []T {
	return pkg.Chunk(b.items, chunkSize, fn)
}

func (b *BaseCollection[T]) For(fn func(v T, i int)) []T {
	return pkg.For(b.items, fn)
}

func (b *BaseCollection[T]) Remove(index int) {
	b.items = pkg.Remove(b.items, index)
}

func (b *BaseCollection[T]) Concat(items []T) {
	b.items = pkg.Concat(b.items, items)
}
