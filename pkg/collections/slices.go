package collections

import (
	"slices"
)

func OrNil[T any](slice []T) []T {
	if len(slice) == 0 {
		return nil
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

func FilterMap[T, U any](data []T, filterMapFn func(T) (U, bool)) []U {
	out := make([]U, 0)

	for _, datum := range data {
		if o, ok := filterMapFn(datum); ok {
			out = append(out, o)
		}
	}

	return out
}

func Sort[T any](data []T, sortFn func(a, b T) int) []T {
	d := slices.Clone(data)

	slices.SortStableFunc(d, sortFn)

	return d
}

func InterfaceToStringSlice(input any) []string {
	if input == nil {
		return []string{}
	}

	// Tentative de conversion en []any
	if slice, ok := input.([]any); ok {
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

func InterfaceToInt32Slice(input any) []int32 {
	if input == nil {
		return []int32{}
	}

	// Tentative de conversion en []any
	if slice, ok := input.([]any); ok {
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
