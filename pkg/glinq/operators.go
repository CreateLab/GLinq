package glinq

// Where filters elements by predicate.
func (s *stream[T]) Where(predicate func(T) bool) Stream[T] {
	oldSource := s.source
	return &stream[T]{
		source: func() (T, bool) {
			for {
				value, ok := oldSource()
				if !ok {
					var zero T
					return zero, false
				}
				if predicate(value) {
					return value, true
				}
			}
		},
	}
}

// Select transforms elements to the same type.
// Supports method chaining.
//
// Example:
//
//	doubled := From([]int{1, 2, 3}).
//	    Select(func(x int) int { return x * 2 }).
//	    ToSlice()
//	// []int{2, 4, 6}
func (s *stream[T]) Select(mapper func(T) T) Stream[T] {
	oldSource := s.source
	return &stream[T]{
		source: func() (T, bool) {
			value, ok := oldSource()
			if !ok {
				var zero T
				return zero, false
			}
			return mapper(value), true
		},
	}
}

// SelectWithIndex transforms elements to the same type, providing index to mapper function.
// Supports method chaining.
//
// Example:
//
//	doubled := From([]int{1, 2, 3}).
//	    SelectWithIndex(func(x int, idx int) int { return x * idx }).
//	    ToSlice()
//	// []int{0, 2, 6}
func (s *stream[T]) SelectWithIndex(mapper func(T, int) T) Stream[T] {
	oldSource := s.source
	index := 0
	return &stream[T]{
		source: func() (T, bool) {
			value, ok := oldSource()
			if !ok {
				var zero T
				return zero, false
			}
			result := mapper(value, index)
			index++
			return result, true
		},
	}
}

// Map transforms elements to a different type.
// This is a function (not a method) because in Go methods cannot have their own type parameters.
//
// Example:
//
//	strings := Select(
//	    From([]int{1, 2, 3}),
//	    func(x int) string { return fmt.Sprintf("num_%d", x) },
//	).ToSlice()
//	// []string{"num_1", "num_2", "num_3"}
func Select[T, R any](s Stream[T], mapper func(T) R) Stream[R] {
	return &stream[R]{
		source: func() (R, bool) {
			value, ok := s.Next()
			if !ok {
				var zero R
				return zero, false
			}
			return mapper(value), true
		},
	}
}

// SelectWithIndex transforms elements to a different type, providing index to mapper function.
// This is a function (not a method) because in Go methods cannot have their own type parameters.
//
// Example:
//
//	strings := SelectWithIndex(
//	    From([]int{1, 2, 3}),
//	    func(x int, idx int) string { return fmt.Sprintf("num_%d_at_%d", x, idx) },
//	).ToSlice()
//	// []string{"num_1_at_0", "num_2_at_1", "num_3_at_2"}
func SelectWithIndex[T, R any](s Stream[T], mapper func(T, int) R) Stream[R] {
	index := 0
	return &stream[R]{
		source: func() (R, bool) {
			value, ok := s.Next()
			if !ok {
				var zero R
				return zero, false
			}
			result := mapper(value, index)
			index++
			return result, true
		},
	}
}

// Take takes the first n elements from Stream.
func (s *stream[T]) Take(n int) Stream[T] {
	oldSource := s.source
	count := 0
	return &stream[T]{
		source: func() (T, bool) {
			if count >= n {
				var zero T
				return zero, false
			}
			value, ok := oldSource()
			if !ok {
				var zero T
				return zero, false
			}
			count++
			return value, true
		},
	}
}

// Skip skips the first n elements from Stream.
func (s *stream[T]) Skip(n int) Stream[T] {
	oldSource := s.source
	skipped := 0
	return &stream[T]{
		source: func() (T, bool) {
			for skipped < n {
				_, ok := oldSource()
				if !ok {
					var zero T
					return zero, false
				}
				skipped++
			}
			return oldSource()
		},
	}
}
