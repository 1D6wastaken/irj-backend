package collections

func MapFilter[K comparable, V any](m map[K]V, f func(k K, v V) bool) map[K]V {
	newMap := make(map[K]V, len(m))

	for k, v := range m {
		if f(k, v) {
			newMap[k] = v
		}
	}

	return newMap
}
