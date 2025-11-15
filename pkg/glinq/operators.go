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

// TakeOrderedBy takes the first n smallest elements from Stream using comparator function.
// Uses a buffer of size n to keep track of the n smallest elements.
// Stream is read lazily, and only the buffer is sorted before being returned.
// This is a function (not a method) because in Go methods cannot have their own type parameters.
//
// Example:
//
//	type Person struct { Age int; Name string }
//	people := []Person{{Age: 30, Name: "Alice"}, {Age: 25, Name: "Bob"}, {Age: 35, Name: "Charlie"}}
//	result := TakeOrderedBy(
//	    From(people),
//	    2,
//	    func(a, b Person) bool { return a.Age < b.Age },
//	).ToSlice()
//	// []Person{{Age: 25, Name: "Bob"}, {Age: 30, Name: "Alice"}}
func TakeOrderedBy[T any](s Stream[T], n int, less func(a, b T) bool) Stream[T] {
	if n <= 0 {
		return Empty[T]()
	}

	var buf []T
	var result []T
	var materialized bool

	return &stream[T]{
		source: func() (T, bool) {
			if !materialized {
				// Materialize stream and build buffer
				for {
					value, ok := s.Next()
					if !ok {
						break
					}

					if len(buf) < n {
						buf = append(buf, value)
					} else {
						// Find maximum element in buffer (according to less)
						maxIdx := 0
						for i := 1; i < len(buf); i++ {
							if less(buf[maxIdx], buf[i]) {
								maxIdx = i
							}
						}
						// Replace if current value is smaller (according to less)
						if less(value, buf[maxIdx]) {
							buf[maxIdx] = value
						}
					}
				}

				// Sort buffer in ascending order using less function
				result = make([]T, len(buf))
				copy(result, buf)
				for i := 0; i < len(result); i++ {
					for j := i + 1; j < len(result); j++ {
						if less(result[j], result[i]) {
							result[i], result[j] = result[j], result[i]
						}
					}
				}

				materialized = true
			}

			if len(result) == 0 {
				var zero T
				return zero, false
			}

			value := result[0]
			result = result[1:]
			return value, true
		},
	}
}

// TakeOrderedDescendingBy takes the first n largest elements from Stream using comparator function.
// Uses a buffer of size n to keep track of the n largest elements.
// Stream is read lazily, and only the buffer is sorted before being returned.
// This is a function (not a method) because in Go methods cannot have their own type parameters.
//
// Example:
//
//	type Person struct { Age int; Name string }
//	people := []Person{{Age: 30, Name: "Alice"}, {Age: 25, Name: "Bob"}, {Age: 35, Name: "Charlie"}}
//	result := TakeOrderedDescendingBy(
//	    From(people),
//	    2,
//	    func(a, b Person) bool { return a.Age < b.Age },
//	).ToSlice()
//	// []Person{{Age: 35, Name: "Charlie"}, {Age: 30, Name: "Alice"}}
func TakeOrderedDescendingBy[T any](s Stream[T], n int, less func(a, b T) bool) Stream[T] {
	if n <= 0 {
		return Empty[T]()
	}

	var buf []T
	var result []T
	var materialized bool

	return &stream[T]{
		source: func() (T, bool) {
			if !materialized {
				// Materialize stream and build buffer
				for {
					value, ok := s.Next()
					if !ok {
						break
					}

					if len(buf) < n {
						buf = append(buf, value)
					} else {
						// Find minimum element in buffer (according to less)
						minIdx := 0
						for i := 1; i < len(buf); i++ {
							if less(buf[i], buf[minIdx]) {
								minIdx = i
							}
						}
						// Replace if current value is larger (according to less)
						if less(buf[minIdx], value) {
							buf[minIdx] = value
						}
					}
				}

				// Sort buffer in descending order using less function
				result = make([]T, len(buf))
				copy(result, buf)
				for i := 0; i < len(result); i++ {
					for j := i + 1; j < len(result); j++ {
						if less(result[i], result[j]) {
							result[i], result[j] = result[j], result[i]
						}
					}
				}

				materialized = true
			}

			if len(result) == 0 {
				var zero T
				return zero, false
			}

			value := result[0]
			result = result[1:]
			return value, true
		},
	}
}
