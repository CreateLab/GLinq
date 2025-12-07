package glinq

// Enumerable is the minimal interface for iterable collections.
// Any type that can provide a sequence of elements.
type Enumerable[T any] interface {
	// Next returns the next element and true, or zero value and false if there are no more elements.
	Next() (T, bool)
}

// Sizable extends Enumerable with optional size information.
// This is an optional interface - not all Enumerables need to implement it.
// Size information is used for performance optimization of terminal operations.
type Sizable[T any] interface {
	Enumerable[T]

	// Size returns the known size of the collection and true,
	// or 0 and false if size is unknown.
	// This is a hint for optimization - implementations should return
	// (0, false) rather than computing size expensively.
	Size() (int, bool)
}

// Stream extends Enumerable and adds operators for working with sequences.
type Stream[T any] interface {
	Enumerable[T] // Embed Enumerable
	Sizable[T]    // Embed Sizable for size information
	// Where filters elements by predicate.
	Where(predicate func(T) bool) Stream[T]
	// Select transforms elements to the same type.
	Select(mapper func(T) T) Stream[T]
	// Take takes the first n elements from Stream.
	Take(n int) Stream[T]
	// TakeWhile takes elements while the predicate returns true.
	// Stops at the first element where predicate returns false.
	//
	// SIZE: Loses size (unknown how many elements satisfy predicate).
	//
	// Example:
	//
	//	numbers := []int{1, 2, 3, 4, 5, 6}
	//	result := From(numbers).TakeWhile(func(x int) bool { return x < 4 }).ToSlice()
	//	// [1, 2, 3]
	TakeWhile(predicate func(T) bool) Stream[T]
	// Skip skips the first n elements from Stream.
	Skip(n int) Stream[T]
	// SkipWhile skips elements while the predicate returns true.
	// Starts returning elements at the first element where predicate returns false.
	//
	// SIZE: Loses size (unknown how many elements to skip).
	//
	// Example:
	//
	//	numbers := []int{1, 2, 3, 4, 5, 6}
	//	result := From(numbers).SkipWhile(func(x int) bool { return x < 4 }).ToSlice()
	//	// [4, 5, 6]
	SkipWhile(predicate func(T) bool) Stream[T]
	// ToSlice materializes Stream into a slice.
	ToSlice() []T
	// First returns the first element and true, or zero value and false if Stream is empty.
	First() (T, bool)
	// Count returns the number of elements in Stream.
	Count() int
	// Any checks if there is at least one element in the Stream.
	// OPTIMIZATION: Returns O(1) if size is known, otherwise iterates until first element.
	Any() bool
	// AnyMatch checks if there is at least one element satisfying the predicate.
	AnyMatch(predicate func(T) bool) bool
	// All checks if all elements satisfy the predicate.
	All(predicate func(T) bool) bool
	// ForEach executes an action for each element in Stream.
	ForEach(action func(T))
	// Chunk splits Stream into chunks of specified size and returns slice of slices.
	Chunk(size int) [][]T
	// Last returns the last element and true, or zero value and false if Stream is empty.
	Last() (T, bool)
	// ElementAt returns the element at the specified index and true, or zero value and false if index is out of range.
	// Index is zero-based. Negative indices are treated as out of range.
	//
	// Example:
	//   numbers := []int{10, 20, 30, 40}
	//   value, ok := From(numbers).ElementAt(2)
	//   // value = 30, ok = true
	ElementAt(index int) (T, bool)
	// ElementAtOrDefault returns the element at the specified index, or the default value if index is out of range.
	// Index is zero-based. Negative indices return the default value.
	//
	// Example:
	//   numbers := []int{10, 20, 30}
	//   value := From(numbers).ElementAtOrDefault(5, 0)
	//   // value = 0 (default value)
	ElementAtOrDefault(index int, defaultValue T) T
	// Contains checks if the Stream contains the specified element.
	// Uses reflect.DeepEqual for comparison, which works with all types including non-comparable ones.
	//
	// Example:
	//   numbers := []int{1, 2, 3, 4, 5}
	//   hasThree := From(numbers).Contains(3)
	//   // hasThree = true
	Contains(value T) bool
	// ContainsBy checks if the Stream contains an element matching the specified key.
	// keySelector extracts a key from each element, and the result is compared with the target key.
	// This allows custom comparison logic, similar to DistinctBy.
	//
	// Example:
	//   type Person struct { ID int; Name string }
	//   people := []Person{{ID: 1, Name: "Alice"}, {ID: 2, Name: "Bob"}}
	//   hasPersonWithID1 := From(people).ContainsBy(1, func(p Person) any { return p.ID })
	//   // hasPersonWithID1 = true
	ContainsBy(targetKey any, keySelector func(T) any) bool
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
	// OrderBy sorts elements using a comparator function.
	// Comparator should return: negative value if a < b,
	// 0 if a == b, positive if a > b.
	// NOTE: OrderBy materializes the entire stream for sorting (partially lazy).
	//
	// Example:
	//   sorted := From([]int{5, 2, 8}).
	//       OrderBy(func(a, b int) int { return a - b }).
	//       ToSlice()
	//   // [2, 5, 8]
	OrderBy(comparator func(T, T) int) Stream[T]
	// OrderByDescending sorts elements in reverse order.
	// This is a shortcut for OrderBy with inverted comparator.
	//
	// Example:
	//   sorted := From([]int{5, 2, 8}).
	//       OrderByDescending(func(a, b int) int { return a - b }).
	//       ToSlice()
	//   // [8, 5, 2]
	OrderByDescending(comparator func(T, T) int) Stream[T]
	// DistinctBy removes duplicates by key extracted by keySelector.
	// keySelector should return a comparable value.
	// RUNTIME REQUIREMENT: returned value must be comparable,
	// otherwise panic will occur.
	//
	// Example:
	//   type Person struct { ID int; Name string }
	//   unique := From(people).
	//       DistinctBy(func(p Person) any { return p.ID }).
	//       ToSlice()
	DistinctBy(keySelector func(T) any) Stream[T]
	// Concat concatenates the current Stream with another Enumerable, preserving duplicates.
	// Elements from the current Stream come first, then elements from other.
	//
	// Example:
	//   result := From([]int{1, 2}).
	//       Concat(From([]int{2, 3})).
	//       ToSlice()
	//   // [1, 2, 2, 3]
	Concat(other Enumerable[T]) Stream[T]
	// Reverse reverses the order of elements in the Stream.
	// NOTE: Reverse materializes the entire stream (partially lazy).
	//
	// Example:
	//   reversed := From([]int{1, 2, 3, 4}).
	//       Reverse().
	//       ToSlice()
	//   // [4, 3, 2, 1]
	Reverse() Stream[T]
}

// stream represents the internal implementation of Stream.
type stream[T any] struct {
	sourceFactory   func() func() (T, bool)
	currentIterator func() (T, bool) // For Enumerable.Next()
	size            int              // -1 if unknown, actual size if known
}

// From creates a Stream from a slice.
//
// PERFORMANCE: The stream holds a reference to the original slice (zero-copy approach).
// This matches C# LINQ behavior for maximum performance.
//
// WARNING: Modifying the slice during iteration may produce unexpected results.
// If you need protection from modifications, use FromSafe() instead.
//
// Example:
//
//	numbers := []int{1, 2, 3, 4, 5}
//	stream := From(numbers)
//	// Efficient - no copying!
func From[T any](slice []T) Stream[T] {
	return &stream[T]{
		sourceFactory: func() func() (T, bool) {
			index := 0 // Fresh index for each iterator
			return func() (T, bool) {
				if index >= len(slice) {
					var zero T
					return zero, false
				}
				result := slice[index]
				index++
				return result, true
			}
		},
		size: len(slice),
	}
}

// FromSafe creates a Stream from a slice with defensive copying.
//
// SAFETY: The slice is copied, so modifications to the original slice
// will not affect the Stream. Use this when you need isolation.
//
// PERFORMANCE: Copying large slices can be expensive.
// If performance is critical and you control the slice lifecycle, use From() instead.
//
// Example:
//
//	numbers := []int{1, 2, 3, 4, 5}
//	stream := FromSafe(numbers)
//	numbers[0] = 999 // Won't affect the stream
func FromSafe[T any](slice []T) Stream[T] {
	// Copy the slice to avoid issues with changes to the original slice
	data := make([]T, len(slice))
	copy(data, slice)

	return &stream[T]{
		sourceFactory: func() func() (T, bool) {
			index := 0 // Fresh index for each iterator
			return func() (T, bool) {
				if index >= len(data) {
					var zero T
					return zero, false
				}
				result := data[index]
				index++
				return result, true
			}
		},
		size: len(data),
	}
}

// Empty creates an empty Stream that contains no elements.
func Empty[T any]() Stream[T] {
	return &stream[T]{
		sourceFactory: func() func() (T, bool) {
			return func() (T, bool) {
				var zero T
				return zero, false
			}
		},
		size: 0,
	}
}

// Range creates a Stream of integers from start to start+count-1.
// If count is negative, returns an empty Stream.
func Range(start, count int) Stream[int] {
	if count < 0 {
		return Empty[int]()
	}

	return &stream[int]{
		sourceFactory: func() func() (int, bool) {
			index := 0 // Fresh index for each iterator
			return func() (int, bool) {
				if index >= count {
					return 0, false
				}
				result := start + index
				index++
				return result, true
			}
		},
		size: count,
	}
}

// Next implements Enumerable
func (s *stream[T]) Next() (T, bool) {
	if s.currentIterator == nil {
		s.currentIterator = s.sourceFactory()
	}
	return s.currentIterator()
}

// Size implements Sizable
func (s *stream[T]) Size() (int, bool) {
	if s.size == -1 {
		return 0, false
	}
	return s.size, true
}

// FromEnumerable creates a Stream from any Enumerable.
// SIZE: Preserves size if source is Sizable, otherwise unknown.
func FromEnumerable[T any](enum Enumerable[T]) Stream[T] {
	size := -1
	if sizable, ok := enum.(Sizable[T]); ok {
		if s, known := sizable.Size(); known {
			size = s
		}
	}
	return &stream[T]{
		sourceFactory: func() func() (T, bool) {
			return enum.Next
		},
		size: size,
	}
}
