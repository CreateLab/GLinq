package glinq

// ToSlice materializes Stream into a slice.
func (s *stream[T]) ToSlice() []T {
	var result []T
	for {
		value, ok := s.source()
		if !ok {
			break
		}
		result = append(result, value)
	}
	return result
}

// First returns the first element and true, or zero value and false if Stream is empty.
func (s *stream[T]) First() (T, bool) {
	return s.source()
}

// Count returns the number of elements in Stream.
func (s *stream[T]) Count() int {
	count := 0
	for {
		_, ok := s.source()
		if !ok {
			break
		}
		count++
	}
	return count
}

// Any checks if there is at least one element satisfying the predicate.
func (s *stream[T]) Any(predicate func(T) bool) bool {
	for {
		value, ok := s.source()
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
	for {
		value, ok := s.source()
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
	for {
		value, ok := s.source()
		if !ok {
			break
		}
		action(value)
	}
}

// ToMapBy materializes Stream[T] into a map using selectors for key and value.
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
	s Stream[T],
	keySelector func(T) K,
	valueSelector func(T) V,
) map[K]V {
	result := make(map[K]V)
	for {
		item, ok := s.Next()
		if !ok {
			break
		}
		key := keySelector(item)
		value := valueSelector(item)
		result[key] = value
	}
	return result
}
