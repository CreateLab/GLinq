package glinq

import "sort"

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
func Select[T, R any](enum Enumerable[T], mapper func(T) R) Stream[R] {
	return &stream[R]{
		source: func() (R, bool) {
			value, ok := enum.Next()
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
func SelectWithIndex[T, R any](enum Enumerable[T], mapper func(T, int) R) Stream[R] {
	index := 0
	return &stream[R]{
		source: func() (R, bool) {
			value, ok := enum.Next()
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
// Example:
//
//	numbers := []int{5, 2, 8, 1, 9, 3}
//	top3 := TakeOrderedBy(From(numbers), 3, func(a, b int) bool { return a < b })
//	// Returns the 3 smallest elements: [1, 2, 3]
func TakeOrderedBy[T any](enum Enumerable[T], n int, less func(a, b T) bool) Stream[T] {
	if n <= 0 {
		return Empty[T]()
	}

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
		return Empty[T]()
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

	return From(heap)
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
