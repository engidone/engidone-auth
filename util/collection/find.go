package collection

func Find[T any](items []T, predicate func(item T) bool) (T, bool) {
	for _, item := range items {
		if predicate(item) {
			return item, true
		}
	}

	var zero T
	return zero, false
}
