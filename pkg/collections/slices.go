package collections

import (
	"slices"
)

func OrEmptySlice[T any](slice []T) []T {
	if slice == nil {
		return []T{}
	}

	return slice
}

func Map[T, U any](data []T, mapFn func(T) U) []U {
	out := make([]U, 0, len(data))

	for _, datum := range data {
		out = append(out, mapFn(datum))
	}

	return out
}

func FlatMap[T, U any](data []T, mapFn func(T) []U) []U {
	out := make([]U, 0)

	for _, datum := range data {
		out = append(out, mapFn(datum)...)
	}

	return out
}

func Filter[T any](data []T, filterFn func(T) bool) []T {
	out := make([]T, 0)

	for _, datum := range data {
		if filterFn(datum) {
			out = append(out, datum)
		}
	}

	return out
}

func FilterMap[T, U any](data []T, filterMapFn func(T) (U, bool)) []U {
	out := make([]U, 0)

	for _, datum := range data {
		if o, ok := filterMapFn(datum); ok {
			out = append(out, o)
		}
	}

	return out
}

//nolint:nonamedreturns
func Find[T any](data []T, findFn func(T) bool) (t T, ok bool) {
	for _, datum := range data {
		if findFn(datum) {
			return datum, true
		}
	}

	return
}

func Some[T any](data []T, someFn func(T) bool) bool {
	for _, datum := range data {
		if someFn(datum) {
			return true
		}
	}

	return false
}

func Every[T any](data []T, someFn func(T) bool) bool {
	for _, datum := range data {
		if !someFn(datum) {
			return false
		}
	}

	return true
}

func Reduce[T, A any](data []T, reduceFn func(A, T) A, accumulator A) A {
	for _, datum := range data {
		accumulator = reduceFn(accumulator, datum)
	}

	return accumulator
}

func Sort[T any](data []T, sortFn func(a, b T) int) []T {
	d := slices.Clone(data)

	slices.SortStableFunc(d, sortFn)

	return d
}

func InterfaceToStringSlice(input interface{}) []string {
	if input == nil {
		return []string{}
	}

	// Tentative de conversion en []interface{}
	if slice, ok := input.([]interface{}); ok {
		result := make([]string, 0, len(slice))
		for _, v := range slice {
			if str, ok := v.(string); ok {
				result = append(result, str)
			}
		}

		return result
	}

	// Tentative de cast direct en []string
	if s, ok := input.([]string); ok {
		return s
	}

	return []string{}
}

func InterfaceToInt32Slice(input interface{}) []int32 {
	if input == nil {
		return []int32{}
	}

	// Tentative de conversion en []interface{}
	if slice, ok := input.([]interface{}); ok {
		result := make([]int32, 0, len(slice))
		for _, v := range slice {
			if str, ok := v.(int32); ok {
				result = append(result, str)
			}
		}

		return result
	}

	// Tentative de cast direct en []int32
	if s, ok := input.([]int32); ok {
		return s
	}

	return []int32{}
}
