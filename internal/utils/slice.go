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
