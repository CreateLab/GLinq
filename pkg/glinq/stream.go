package glinq

// Stream represents a lazy sequence of elements.
// Elements are provided through a source function that returns (element, has_element).
type Stream[T any] interface {
	// Next returns the next element and true, or zero value and false if there are no more elements.
	Next() (T, bool)
	// Where filters elements by predicate.
	Where(predicate func(T) bool) Stream[T]
	// Select transforms elements to the same type.
	Select(mapper func(T) T) Stream[T]
	// Take takes the first n elements from Stream.
	Take(n int) Stream[T]
	// Skip skips the first n elements from Stream.
	Skip(n int) Stream[T]
	// ToSlice materializes Stream into a slice.
	ToSlice() []T
	// First returns the first element and true, or zero value and false if Stream is empty.
	First() (T, bool)
	// Count returns the number of elements in Stream.
	Count() int
	// Any checks if there is at least one element satisfying the predicate.
	Any(predicate func(T) bool) bool
	// All checks if all elements satisfy the predicate.
	All(predicate func(T) bool) bool
	// ForEach executes an action for each element in Stream.
	ForEach(action func(T))
	// Chunk splits Stream into chunks of specified size and returns slice of slices.
	Chunk(size int) [][]T
	// Last returns the last element and true, or zero value and false if Stream is empty.
	Last() (T, bool)
	// Min returns the minimum element using comparator function.
	// Comparator should return negative value if first < second, 0 if equal, positive if first > second.
	Min(comparator func(T, T) int) (T, bool)
	// Max returns the maximum element using comparator function.
	// Comparator should return negative value if first < second, 0 if equal, positive if first > second.
	Max(comparator func(T, T) int) (T, bool)
	// SelectWithIndex transforms elements to the same type, providing index to mapper function.
	SelectWithIndex(mapper func(T, int) T) Stream[T]
	// Aggregate applies an accumulator function over the Stream.
	// The seed parameter is the initial accumulator value.
	// The accumulator function is invoked for each element.
	Aggregate(seed T, accumulator func(T, T) T) T
}

// stream represents the internal implementation of Stream.
type stream[T any] struct {
	source func() (T, bool)
}

// From creates a Stream from a slice.
// The slice is copied, so changes to the original slice do not affect the Stream.
func From[T any](slice []T) Stream[T] {
	// Copy the slice to avoid issues with changes to the original slice
	data := make([]T, len(slice))
	copy(data, slice)

	index := 0
	return &stream[T]{
		source: func() (T, bool) {
			if index < len(data) {
				result := data[index]
				index++
				return result, true
			}
			var zero T
			return zero, false
		},
	}
}

// Empty creates an empty Stream that contains no elements.
func Empty[T any]() Stream[T] {
	return &stream[T]{
		source: func() (T, bool) {
			var zero T
			return zero, false
		},
	}
}

// Range creates a Stream of integers from start to start+count-1.
func Range(start, count int) Stream[int] {
	index := 0
	return &stream[int]{
		source: func() (int, bool) {
			if index < count {
				result := start + index
				index++
				return result, true
			}
			return 0, false
		},
	}
}

// Next returns the next element from stream.
func (s *stream[T]) Next() (T, bool) {
	return s.source()
}
