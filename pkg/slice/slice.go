package slice

func Map[T interface{}](s []T, fn func(v T, i int) T) []T {
	var mapped []T

	for i, v := range s {
		mapped = append(mapped, fn(v, i))
	}

	return mapped
}

func Filter[T interface{}](s []T, fn func(v T, i int) bool) []T {
	var filtered []T

	for i, v := range s {
		f := fn(v, i)
		if f {
			filtered = append(filtered, v)
		}
	}

	return filtered
}

func Except[T interface{}](s []T, fn func(v T, i int) bool) []T {
	var excepted []T

	for i, v := range s {
		f := fn(v, i)
		if !f {
			excepted = append(excepted, v)
		}
	}

	return excepted
}

func Chunk[T interface{}](s []T, chunkSize int, fn ...func(v []T, i int)) [][]T {
	chunkSlice := make([]T, 0)
	chunkedSlice := make([][]T, 0)

	chunkedIndex := 0
	for _, v := range s {
		chunkSlice = append(chunkSlice, v)

		if chunkSize == len(chunkSlice) {
			if len(fn) != 0 {
				fn[0](chunkSlice, chunkedIndex)
			}

			chunkedSlice = append(chunkedSlice, chunkSlice)
			chunkedIndex++

			chunkSlice = make([]T, 0)
		}
	}

	if len(chunkSlice) != 0 {
		if len(fn) != 0 {
			fn[0](chunkSlice, chunkedIndex)
			chunkedSlice = append(chunkedSlice, chunkSlice)
		}
	}

	return chunkedSlice
}

func For[T interface{}](s []T, fn func(v T, i int)) []T {
	for i, v := range s {
		fn(v, i)
	}

	return s
}

func Add[T interface{}](s []T, v T) []T {
	return append(s, v)
}

func Remove[T interface{}](s []T, index int) []T {
	return append(s[:index], s[index+1:]...)
}

func Concat[T interface{}](s []T, s2 []T) []T {
	if len(s2) == 0 {
		return s
	}

	for _, v := range s2 {
		s = append(s, v)
	}

	return s
}
