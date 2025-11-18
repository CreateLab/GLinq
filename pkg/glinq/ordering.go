package glinq

import "sort"

// orderBy is a common sorting function used by OrderBy and OrderByDescending.
func (s *stream[T]) orderBy(ascending bool, comparator func(T, T) int) Stream[T] {
	sorted := s.ToSlice()

	sort.Slice(sorted, func(i, j int) bool {
		cmp := comparator(sorted[i], sorted[j])
		if ascending {
			return cmp < 0
		}
		return cmp > 0
	})

	// SIZE: Preserves size (1-to-1 transformation, materializes)
	size := len(sorted)
	return &stream[T]{
		sourceFactory: func() func() (T, bool) {
			index := 0 // Fresh index for each iterator
			return func() (T, bool) {
				if index >= len(sorted) {
					var zero T
					return zero, false
				}
				result := sorted[index]
				index++
				return result, true
			}
		},
		size: &size, // PRESERVE: 1-to-1 transformation
	}
}

// OrderBy sorts elements in ascending order.
func (s *stream[T]) OrderBy(comparator func(T, T) int) Stream[T] {
	return s.orderBy(true, comparator)
}

// OrderByDescending sorts elements in descending order.
func (s *stream[T]) OrderByDescending(comparator func(T, T) int) Stream[T] {
	return s.orderBy(false, comparator)
}
