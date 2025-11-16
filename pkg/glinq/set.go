package glinq

// Concat concatenates the current Stream with another Enumerable (preserving duplicates).
// Elements from the current Stream come first, then elements from other.
func (s *stream[T]) Concat(other Enumerable[T]) Stream[T] {
	firstExhausted := false

	return &stream[T]{
		source: func() (T, bool) {
			if !firstExhausted {
				val, ok := s.source()
				if ok {
					return val, true
				}
				firstExhausted = true
			}

			return other.Next()
		},
	}
}

// Union returns the union of two Enumerables (all unique elements from both).
// T must be comparable, otherwise code will not compile.
// This is a function (not a method) because methods cannot have their own type constraints.
//
// Example:
//
//	set1 := []int{1, 2, 3}
//	set2 := []int{3, 4, 5}
//	union := Union(From(set1), From(set2)).ToSlice()
//	// [1, 2, 3, 4, 5]
func Union[T comparable](e1, e2 Enumerable[T]) Stream[T] {
	seen := make(map[T]bool)
	var current Enumerable[T] = e1
	secondStarted := false

	return &stream[T]{
		source: func() (T, bool) {
			for {
				val, ok := current.Next()

				// Switch to second enumerable
				if !ok {
					if secondStarted {
						var zero T
						return zero, false
					}
					current = e2
					secondStarted = true
					continue
				}

				// Return only unique elements
				if !seen[val] {
					seen[val] = true
					return val, true
				}
			}
		},
	}
}

// Intersect returns the intersection of two Enumerables.
// T must be comparable, otherwise code will not compile.
// This is a function (not a method) because methods cannot have their own type constraints.
func Intersect[T comparable](e1, e2 Enumerable[T]) Stream[T] {
	// Materialize e2 into a set
	otherSet := make(map[T]bool)
	for {
		val, ok := e2.Next()
		if !ok {
			break
		}
		otherSet[val] = true
	}

	seen := make(map[T]bool)
	return &stream[T]{
		source: func() (T, bool) {
			for {
				val, ok := e1.Next()
				if !ok {
					var zero T
					return zero, false
				}

				if otherSet[val] && !seen[val] {
					seen[val] = true
					return val, true
				}
			}
		},
	}
}

// Except returns the difference of Enumerables (elements from current that are not in other).
// T must be comparable, otherwise code will not compile.
// This is a function (not a method) because methods cannot have their own type constraints.
func Except[T comparable](e1, e2 Enumerable[T]) Stream[T] {
	// Materialize e2 into a set
	otherSet := make(map[T]bool)
	for {
		val, ok := e2.Next()
		if !ok {
			break
		}
		otherSet[val] = true
	}

	seen := make(map[T]bool)
	return &stream[T]{
		source: func() (T, bool) {
			for {
				val, ok := e1.Next()
				if !ok {
					var zero T
					return zero, false
				}

				if !otherSet[val] && !seen[val] {
					seen[val] = true
					return val, true
				}
			}
		},
	}
}
