package glinq

// ToSlice materializes Stream into a slice.
// OPTIMIZATION: Preallocates capacity if size is known.
func (s *stream[T]) ToSlice() []T {
	iterator := s.sourceFactory() // Fresh iterator
	var result []T
	// OPTIMIZATION: preallocate if size known
	if s.size != -1 {
		result = make([]T, 0, s.size)
	}
	for {
		value, ok := iterator()
		if !ok {
			break
		}
		result = append(result, value)
	}
	return result
}

// First returns the first element and true, or zero value and false if Stream is empty.
func (s *stream[T]) First() (T, bool) {
	iterator := s.sourceFactory() // Fresh iterator
	return iterator()
}

// Count returns the number of elements in Stream.
// OPTIMIZATION: Returns O(1) if size is known, otherwise O(n).
func (s *stream[T]) Count() int {
	// OPTIMIZATION: O(1) if size known!
	if s.size != -1 {
		return s.size
	}

	// Fallback: iterate and count
	iterator := s.sourceFactory() // Fresh iterator
	count := 0
	for {
		_, ok := iterator()
		if !ok {
			break
		}
		count++
	}
	return count
}

// Any checks if there is at least one element in the Stream.
// OPTIMIZATION: Returns O(1) if size is known, otherwise iterates until first element.
func (s *stream[T]) Any() bool {
	// OPTIMIZATION: O(1) if size known!
	if s.size != -1 {
		return s.size > 0
	}

	// Fallback: iterate until first element
	iterator := s.sourceFactory() // Fresh iterator
	_, ok := iterator()
	return ok
}

// AnyMatch checks if there is at least one element satisfying the predicate.
func (s *stream[T]) AnyMatch(predicate func(T) bool) bool {
	iterator := s.sourceFactory() // Fresh iterator
	for {
		value, ok := iterator()
		if !ok {
			break
		}
		if predicate(value) {
			return true
		}
	}
	return false
}

// All checks if all elements satisfy the predicate.
func (s *stream[T]) All(predicate func(T) bool) bool {
	iterator := s.sourceFactory() // Fresh iterator
	for {
		value, ok := iterator()
		if !ok {
			break
		}
		if !predicate(value) {
			return false
		}
	}
	return true
}

// ForEach executes an action for each element in Stream.
func (s *stream[T]) ForEach(action func(T)) {
	iterator := s.sourceFactory() // Fresh iterator
	for {
		value, ok := iterator()
		if !ok {
			break
		}
		action(value)
	}
}

// Chunk splits Stream into chunks of specified size and returns slice of slices.
// The last chunk may contain fewer elements than the specified size.
//
// SIZE: Calculated as ceil(sourceSize / size) if source size known, else unknown.
//
// Example:
//
//	numbers := []int{1, 2, 3, 4, 5, 6, 7}
//	chunks := From(numbers).Chunk(3)
//	// [][]int{{1, 2, 3}, {4, 5, 6}, {7}}
func (s *stream[T]) Chunk(size int) [][]T {
	if size <= 0 {
		return nil
	}

	iterator := s.sourceFactory() // Fresh iterator
	var result [][]T
	// OPTIMIZATION: preallocate if size known
	if s.size != -1 {
		chunkCount := (s.size + size - 1) / size // ceil division
		result = make([][]T, 0, chunkCount)
	}
	var currentChunk []T

	for {
		value, ok := iterator()
		if !ok {
			// Add the last chunk if it's not empty
			if len(currentChunk) > 0 {
				result = append(result, currentChunk)
			}
			break
		}

		currentChunk = append(currentChunk, value)

		// When chunk reaches the desired size, add it to result and start a new one
		if len(currentChunk) == size {
			result = append(result, currentChunk)
			currentChunk = nil
		}
	}

	return result
}

// Last returns the last element and true, or zero value and false if Stream is empty.
func (s *stream[T]) Last() (T, bool) {
	iterator := s.sourceFactory() // Fresh iterator
	var last T
	var found bool

	for {
		value, ok := iterator()
		if !ok {
			break
		}
		last = value
		found = true
	}

	return last, found
}

// Min returns the minimum element using comparator function.
// Comparator should return negative value if first < second, 0 if equal, positive if first > second.
// Returns zero value and false if Stream is empty.
//
// Example:
//
//	type Person struct { Age int; Name string }
//	people := []Person{{Age: 30, Name: "Alice"}, {Age: 25, Name: "Bob"}}
//	youngest, ok := From(people).Min(func(a, b Person) int {
//	    return a.Age - b.Age
//	})
//	// youngest = Person{Age: 25, Name: "Bob"}, ok = true
func (s *stream[T]) Min(comparator func(T, T) int) (T, bool) {
	iterator := s.sourceFactory() // Fresh iterator
	var minVal T
	var found bool

	for {
		value, ok := iterator()
		if !ok {
			break
		}
		if !found || comparator(value, minVal) < 0 {
			minVal = value
			found = true
		}
	}

	return minVal, found
}

// Max returns the maximum element using comparator function.
// Comparator should return negative value if first < second, 0 if equal, positive if first > second.
// Returns zero value and false if Stream is empty.
//
// Example:
//
//	type Person struct { Age int; Name string }
//	people := []Person{{Age: 30, Name: "Alice"}, {Age: 25, Name: "Bob"}}
//	oldest, ok := From(people).Max(func(a, b Person) int {
//	    return a.Age - b.Age
//	})
//	// oldest = Person{Age: 30, Name: "Alice"}, ok = true
func (s *stream[T]) Max(comparator func(T, T) int) (T, bool) {
	iterator := s.sourceFactory() // Fresh iterator
	var maxVal T
	var found bool

	for {
		value, ok := iterator()
		if !ok {
			break
		}
		if !found || comparator(value, maxVal) > 0 {
			maxVal = value
			found = true
		}
	}

	return maxVal, found
}

// Aggregate applies an accumulator function over the Stream.
// The seed parameter is the initial accumulator value.
// The accumulator function is invoked for each element.
// Returns the final accumulator value.
//
// Example:
//
//	numbers := []int{1, 2, 3, 4, 5}
//	sum := From(numbers).Aggregate(0, func(acc, x int) int { return acc + x })
//	// 15
//
//	numbers := []int{1, 2, 3}
//	product := From(numbers).Aggregate(1, func(acc, x int) int { return acc * x })
//	// 6
func (s *stream[T]) Aggregate(seed T, accumulator func(T, T) T) T {
	iterator := s.sourceFactory() // Fresh iterator
	result := seed
	for {
		value, ok := iterator()
		if !ok {
			break
		}
		result = accumulator(result, value)
	}
	return result
}

// ToMapBy materializes Enumerable[T] into a map using selectors for key and value.
//
// Example:
//
//	type User struct { ID int; Name string }
//	users := []User{{ID: 1, Name: "Alice"}, {ID: 2, Name: "Bob"}}
//	userMap := ToMapBy(
//	    From(users),
//	    func(u User) int { return u.ID },
//	    func(u User) string { return u.Name },
//	)
//	// map[int]string{1: "Alice", 2: "Bob"}
func ToMapBy[T any, K comparable, V any](
	enum Enumerable[T],
	keySelector func(T) K,
	valueSelector func(T) V,
) map[K]V {
	result := make(map[K]V)
	for {
		item, ok := enum.Next()
		if !ok {
			break
		}
		key := keySelector(item)
		value := valueSelector(item)
		result[key] = value
	}
	return result
}
