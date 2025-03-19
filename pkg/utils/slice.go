package utils

// FindElementByProperty searches for an element in a slice based on a given match function.
func FindElementByProperty[T any](items []T, matchFunc func(T) bool) bool {
	for _, item := range items {
		if matchFunc(item) {
			return true
		}
	}

	return false
}
