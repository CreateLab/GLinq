package glinq

import "sort"

// orderBy is a common sorting function used by OrderBy and OrderByDescending.
func (s *stream[T]) orderBy(ascending bool, comparator func(T, T) int) Stream[T] {
	var sorted []T
	var initialized bool
	index := 0

	return &stream[T]{
		source: func() (T, bool) {
			if !initialized {
				sorted = s.ToSlice()

				sort.Slice(sorted, func(i, j int) bool {
					cmp := comparator(sorted[i], sorted[j])
					if ascending {
						return cmp < 0
					}
					return cmp > 0
				})

				initialized = true
				index = 0
			}

			if index < len(sorted) {
				result := sorted[index]
				index++
				return result, true
			}

			var zero T
			return zero, false
		},
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
