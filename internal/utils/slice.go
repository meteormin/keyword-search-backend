package utils

func Map[T any](s []T, fn func(v T, i int) T) []T {
	var mapped []T

	for i, v := range s {
		mapped = append(mapped, fn(v, i))
	}

	return mapped
}

func Filter[T any](s []T, fn func(v T, i int) bool) []T {
	var filtered []T

	for i, v := range s {
		f := fn(v, i)
		if f {
			filtered = append(s, v)
		}
	}

	return filtered
}

func Chunk[T any](s []T, chunkSize int, fn func(v []T, i int)) []T {
	chunkSlice := make([]T, 0)

	for i, v := range s {
		chunkSlice = append(chunkSlice, v)

		if chunkSize == len(chunkSlice) {
			fn(chunkSlice, i)
			chunkSlice = make([]T, 0)
		}
	}

	return s
}
