package utils

func Contains[T comparable](s []T, e T) bool {
	return Find(s, e) != -1
}

func Find[T comparable](s []T, e T) int {
	for i := range s {
		if s[i] == e {
			return i
		}
	}

	return -1
}

func FindIndexByFunc[T any](data []T, fn func(T) bool) int {
	for i, datum := range data {
		if fn(datum) {
			return i
		}
	}

	return -1
}

func Some[T any](data []T, some func(t T) bool) bool {
	for _, datum := range data {
		if some(datum) {
			return true
		}
	}

	return false
}

func Filter[T any](data []T, filter func(t T) bool) []T {
	if data == nil {
		return nil
	}

	out := make([]T, 0)

	for _, datum := range data {
		if filter(datum) {
			out = append(out, datum)
		}
	}

	return out
}

func Map[T any, O any](data []T, mapper func(T) O) []O {
	if data == nil {
		return nil
	}

	out := make([]O, len(data))

	for i, datum := range data {
		out[i] = mapper(datum)
	}

	return out
}

func FilterAndMap[T any, O any](data []T, mapper func(T) (O, bool)) []O {
	if data == nil {
		return nil
	}

	out := make([]O, 0)

	for _, datum := range data {
		if newDatum, ok := mapper(datum); ok {
			out = append(out, newDatum)
		}
	}

	return out
}

func SkipWhile[T any](data []T, fn func(T) bool) []T {
	for i, datum := range data {
		if !fn(datum) {
			return data[i:]
		}
	}

	return nil
}

func Reduce[T any, A any](data []T, fn func(acc A, datum T) A, initialValue A) A {
	if data == nil {
		return initialValue
	}

	acc := initialValue
	for _, datum := range data {
		acc = fn(acc, datum)
	}

	return acc
}
