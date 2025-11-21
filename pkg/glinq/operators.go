package glinq

import "sort"

// Where filters elements by predicate.
//
// SIZE: Loses size (unknown how many elements pass filter).
func (s *stream[T]) Where(predicate func(T) bool) Stream[T] {
	return &stream[T]{
		sourceFactory: func() func() (T, bool) {
			source := s.sourceFactory() // Get fresh source
			return func() (T, bool) {
				for {
					value, ok := source()
					if !ok {
						var zero T
						return zero, false
					}
					if predicate(value) {
						return value, true
					}
				}
			}
		},
		size: -1, // LOSE: unknown how many pass filter
	}
}

// Select transforms elements to the same type.
// Supports method chaining.
//
// SIZE: Preserves size (1-to-1 transformation).
//
// Example:
//
//	doubled := From([]int{1, 2, 3}).
//	    Select(func(x int) int { return x * 2 }).
//	    ToSlice()
//	// []int{2, 4, 6}
func (s *stream[T]) Select(mapper func(T) T) Stream[T] {
	return &stream[T]{
		sourceFactory: func() func() (T, bool) {
			source := s.sourceFactory() // Get fresh source
			return func() (T, bool) {
				value, ok := source()
				if !ok {
					var zero T
					return zero, false
				}
				return mapper(value), true
			}
		},
		size: s.size, // PRESERVE: 1-to-1 transformation
	}
}

// SelectWithIndex transforms elements to the same type, providing index to mapper function.
// Supports method chaining.
//
// SIZE: Preserves size (1-to-1 transformation).
//
// Example:
//
//	doubled := From([]int{1, 2, 3}).
//	    SelectWithIndex(func(x int, idx int) int { return x * idx }).
//	    ToSlice()
//	// []int{0, 2, 6}
func (s *stream[T]) SelectWithIndex(mapper func(T, int) T) Stream[T] {
	return &stream[T]{
		sourceFactory: func() func() (T, bool) {
			source := s.sourceFactory() // Get fresh source
			index := 0                  // Fresh counter
			return func() (T, bool) {
				value, ok := source()
				if !ok {
					var zero T
					return zero, false
				}
				result := mapper(value, index)
				index++
				return result, true
			}
		},
		size: s.size, // PRESERVE: 1-to-1 transformation
	}
}

// Map transforms elements to a different type.
// This is a function (not a method) because in Go methods cannot have their own type parameters.
//
// SIZE: Preserves size if source is Sizable (1-to-1 transformation).
//
// Example:
//
//	strings := Select(
//	    From([]int{1, 2, 3}),
//	    func(x int) string { return fmt.Sprintf("num_%d", x) },
//	).ToSlice()
//	// []string{"num_1", "num_2", "num_3"}
func Select[T, R any](enum Enumerable[T], mapper func(T) R) Stream[R] {
	size := -1
	if sizable, ok := enum.(Sizable[T]); ok {
		if s, known := sizable.Size(); known {
			size = s
		}
	}
	return &stream[R]{
		sourceFactory: func() func() (R, bool) {
			return func() (R, bool) {
				value, ok := enum.Next()
				if !ok {
					var zero R
					return zero, false
				}
				return mapper(value), true
			}
		},
		size: size, // PRESERVE if possible
	}
}

// SelectWithIndex transforms elements to a different type, providing index to mapper function.
// This is a function (not a method) because in Go methods cannot have their own type parameters.
//
// SIZE: Preserves size if source is Sizable (1-to-1 transformation).
//
// Example:
//
//	strings := SelectWithIndex(
//	    From([]int{1, 2, 3}),
//	    func(x int, idx int) string { return fmt.Sprintf("num_%d_at_%d", x, idx) },
//	).ToSlice()
//	// []string{"num_1_at_0", "num_2_at_1", "num_3_at_2"}
func SelectWithIndex[T, R any](enum Enumerable[T], mapper func(T, int) R) Stream[R] {
	size := -1
	if sizable, ok := enum.(Sizable[T]); ok {
		if s, known := sizable.Size(); known {
			size = s
		}
	}
	return &stream[R]{
		sourceFactory: func() func() (R, bool) {
			index := 0 // Fresh counter
			return func() (R, bool) {
				value, ok := enum.Next()
				if !ok {
					var zero R
					return zero, false
				}
				result := mapper(value, index)
				index++
				return result, true
			}
		},
		size: size, // PRESERVE if possible
	}
}

// Take takes the first n elements from Stream.
//
// SIZE: Calculated as min(sourceSize, n) if source size known, else n.
// If n is negative, returns an empty Stream.
func (s *stream[T]) Take(n int) Stream[T] {
	if n < 0 {
		return Empty[T]()
	}

	var newSize int
	if s.size != -1 {
		if s.size < n {
			newSize = s.size
		} else {
			newSize = n
		}
	} else {
		newSize = n // Don't know source size, but result won't exceed n
	}

	return &stream[T]{
		sourceFactory: func() func() (T, bool) {
			source := s.sourceFactory() // Get fresh source
			count := 0                  // Fresh counter
			return func() (T, bool) {
				if count >= n {
					var zero T
					return zero, false
				}
				value, ok := source()
				if !ok {
					var zero T
					return zero, false
				}
				count++
				return value, true
			}
		},
		size: newSize,
	}
}

// Skip skips the first n elements from Stream.
//
// SIZE: Calculated as max(0, sourceSize - n) if source size known, else unknown.
// If n is negative, treats it as 0 (no skipping).
func (s *stream[T]) Skip(n int) Stream[T] {
	if n < 0 {
		n = 0
	}

	// Early exit optimization: if skipping more than or equal to known size
	if s.size != -1 && n >= s.size {
		return Empty[T]()
	}

	var newSize = -1
	if s.size != -1 {
		size := s.size - n
		if size < 0 {
			size = 0
		}
		newSize = size
	}
	// else: -1 (unknown)

	return &stream[T]{
		sourceFactory: func() func() (T, bool) {
			source := s.sourceFactory() // Get fresh source
			skipped := 0                // Fresh counter
			return func() (T, bool) {
				for skipped < n {
					_, ok := source()
					if !ok {
						var zero T
						return zero, false
					}
					skipped++
				}
				return source()
			}
		},
		size: newSize,
	}
}

// TakeWhile takes elements while the predicate returns true.
// Stops at the first element where predicate returns false.
//
// SIZE: Loses size (unknown how many elements satisfy predicate).
func (s *stream[T]) TakeWhile(predicate func(T) bool) Stream[T] {
	return &stream[T]{
		sourceFactory: func() func() (T, bool) {
			source := s.sourceFactory() // Get fresh source
			stopped := false            // Flag to stop iteration

			return func() (T, bool) {
				if stopped {
					var zero T
					return zero, false
				}

				value, ok := source()
				if !ok {
					var zero T
					return zero, false
				}

				if !predicate(value) {
					stopped = true
					var zero T
					return zero, false
				}

				return value, true
			}
		},
		size: -1, // LOSE: unknown how many satisfy predicate
	}
}

// SkipWhile skips elements while the predicate returns true.
// Starts returning elements at the first element where predicate returns false.
//
// SIZE: Loses size (unknown how many elements to skip).
func (s *stream[T]) SkipWhile(predicate func(T) bool) Stream[T] {
	return &stream[T]{
		sourceFactory: func() func() (T, bool) {
			source := s.sourceFactory() // Get fresh source
			skipping := true            // Flag to continue skipping

			return func() (T, bool) {
				for skipping {
					value, ok := source()
					if !ok {
						var zero T
						return zero, false
					}

					if !predicate(value) {
						skipping = false
						return value, true
					}
				}

				// After skipping phase, return remaining elements
				return source()
			}
		},
		size: -1, // LOSE: unknown how many to skip
	}
}

// buildHeap is a helper function for building a heap.
func buildHeap[T any](items []T, less func(a, b T) bool) {
	n := len(items)
	for i := n/2 - 1; i >= 0; i-- {
		heapifyDown(items, i, n, less)
	}
}

// heapifyDown performs heapify down operation (sift down).
func heapifyDown[T any](items []T, i, n int, less func(a, b T) bool) {
	for {
		largest := i
		left := 2*i + 1
		right := 2*i + 2

		if left < n && less(items[largest], items[left]) {
			largest = left
		}
		if right < n && less(items[largest], items[right]) {
			largest = right
		}

		if largest == i {
			break
		}

		items[i], items[largest] = items[largest], items[i]
		i = largest
	}
}

// TakeOrderedBy returns the first n elements ordered by the less function.
// Uses a heap-based algorithm for efficient processing.
//
// SIZE: Calculated as min(sourceSize, n) if source size known, else n.
//
// Example:
//
//	numbers := []int{5, 2, 8, 1, 9, 3}
//	top3 := TakeOrderedBy(From(numbers), 3, func(a, b int) bool { return a < b })
//	// Returns the 3 smallest elements: [1, 2, 3]
//
//nolint:gocognit
func TakeOrderedBy[T any](enum Enumerable[T], n int, less func(a, b T) bool) Stream[T] {
	if n <= 0 {
		return Empty[T]()
	}

	var size int
	if sizable, ok := enum.(Sizable[T]); ok {
		if s, known := sizable.Size(); known {
			if s < n {
				size = s
			} else {
				size = n
			}
		} else {
			size = n // Don't know source size, but result won't exceed n
		}
	} else {
		size = n // Don't know source size, but result won't exceed n
	}

	return &stream[T]{
		sourceFactory: func() func() (T, bool) {
			var heap []T

			// Collect first n elements
			for i := 0; i < n; i++ {
				val, ok := enum.Next()
				if !ok {
					break
				}
				heap = append(heap, val)
			}

			if len(heap) == 0 {
				var zero T
				return func() (T, bool) {
					return zero, false
				}
			}

			// Build max-heap
			buildHeap(heap, func(a, b T) bool { return !less(a, b) })

			// Process remaining elements
			for {
				val, ok := enum.Next()
				if !ok {
					break
				}

				if less(val, heap[0]) {
					heap[0] = val
					heapifyDown(heap, 0, len(heap), func(a, b T) bool { return !less(a, b) })
				}
			}

			// Sort the result
			sort.Slice(heap, func(i, j int) bool {
				return less(heap[i], heap[j])
			})

			index := 0
			return func() (T, bool) {
				if index >= len(heap) {
					var zero T
					return zero, false
				}
				result := heap[index]
				index++
				return result, true
			}
		},
		size: size,
	}
}

// TakeOrderedDescendingBy returns the first n elements ordered in descending order by the less function.
// This is equivalent to TakeOrderedBy with inverted comparator.
//
// Example:
//
//	numbers := []int{5, 2, 8, 1, 9, 3}
//	top3 := TakeOrderedDescendingBy(From(numbers), 3, func(a, b int) bool { return a < b })
//	// Returns the 3 largest elements: [9, 8, 5]
func TakeOrderedDescendingBy[T any](enum Enumerable[T], n int, less func(a, b T) bool) Stream[T] {
	return TakeOrderedBy(enum, n, func(a, b T) bool { return !less(a, b) })
}

// Reverse reverses the order of elements in the Stream.
// NOTE: Reverse materializes the entire stream (partially lazy).
//
// SIZE: Preserves size (1-to-1 transformation, materializes).
func (s *stream[T]) Reverse() Stream[T] {
	items := s.ToSlice()
	for i, j := 0, len(items)-1; i < j; i, j = i+1, j-1 {
		items[i], items[j] = items[j], items[i]
	}
	return From(items) // From preserves size
}

// SelectMany transforms each element into a sequence and flattens the resulting sequences.
// This is a function (not a method) because in Go methods cannot have their own type parameters.
//
// SIZE: Loses size (1-to-many transformation, unknown result count).
//
// Example:
//
//	numbers := [][]int{{1, 2}, {3, 4}, {5}}
//	flattened := SelectMany(
//	    From(numbers),
//	    func(slice []int) Enumerable[int] { return From(slice) },
//	).ToSlice()
//	// [1, 2, 3, 4, 5]
//
//nolint:gocognit
func SelectMany[T, R any](enum Enumerable[T], selector func(T) Enumerable[R]) Stream[R] {
	return &stream[R]{
		sourceFactory: func() func() (R, bool) {
			var currentEnum Enumerable[R]
			var hasCurrent bool

			return func() (R, bool) {
				for {
					// If we have a current enumerable, try to get next element from it
					if hasCurrent {
						val, ok := currentEnum.Next()
						if ok {
							return val, true
						}
						// Current enumerable exhausted, move to next
						hasCurrent = false
					}

					// Get next element from source
					elem, ok := enum.Next()
					if !ok {
						var zero R
						return zero, false
					}

					// Transform element into enumerable
					currentEnum = selector(elem)
					hasCurrent = true
				}
			}
		},
		size: -1, // LOSE: 1-to-many transformation
	}
}
