package glinq

// Concat concatenates the current Stream with another Enumerable (preserving duplicates).
// Elements from the current Stream come first, then elements from other.
//
// SIZE: Calculated as currentSize + otherSize if both known, else unknown.
func (s *stream[T]) Concat(other Enumerable[T]) Stream[T] {
	var newSize = -1
	if s.size != -1 {
		if sizable, ok := other.(Sizable[T]); ok {
			if otherSize, known := sizable.Size(); known {
				newSize = s.size + otherSize
			}
		}
	}
	// else: -1 (unknown)

	return &stream[T]{
		sourceFactory: func() func() (T, bool) {
			source := s.sourceFactory() // Get fresh source
			firstExhausted := false     // Fresh flag for each iterator

			return func() (T, bool) {
				if !firstExhausted {
					val, ok := source()
					if ok {
						return val, true
					}
					firstExhausted = true
				}

				return other.Next()
			}
		},
		size: newSize, // CALCULATED: currentSize + otherSize if both known
	}
}

// Union returns the union of two Enumerables (all unique elements from both).
// T must be comparable, otherwise code will not compile.
// This is a function (not a method) because methods cannot have their own type constraints.
//
// SIZE: Loses size (unknown how many duplicates exist).
//
// Example:
//
//	set1 := []int{1, 2, 3}
//	set2 := []int{3, 4, 5}
//	union := Union(From(set1), From(set2)).ToSlice()
//	// [1, 2, 3, 4, 5]
//
//nolint:gocognit
func Union[T comparable](e1, e2 Enumerable[T]) Stream[T] {
	return &stream[T]{
		sourceFactory: func() func() (T, bool) {
			seen := make(map[T]bool)
			var current = e1
			secondStarted := false

			return func() (T, bool) {
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
			}
		},
		size: -1, // LOSE: unknown how many duplicates
	}
}

// Intersect returns the intersection of two Enumerables.
// T must be comparable, otherwise code will not compile.
// This is a function (not a method) because methods cannot have their own type constraints.
//
// SIZE: Loses size (unknown result count).
//
//nolint:gocognit
func Intersect[T comparable](e1, e2 Enumerable[T]) Stream[T] {
	return &stream[T]{
		sourceFactory: func() func() (T, bool) {
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
			return func() (T, bool) {
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
			}
		},
		size: -1, // LOSE: unknown result count
	}
}

// Except returns the difference of Enumerables (elements from current that are not in other).
// T must be comparable, otherwise code will not compile.
// This is a function (not a method) because methods cannot have their own type constraints.
//
// SIZE: Loses size (unknown result count).
//
//nolint:gocognit
func Except[T comparable](e1, e2 Enumerable[T]) Stream[T] {
	return &stream[T]{
		sourceFactory: func() func() (T, bool) {
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
			return func() (T, bool) {
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
			}
		},
		size: -1, // LOSE: unknown result count
	}
}

// Zip combines two Enumerables by applying a result selector function to corresponding elements.
// The result stream stops when either source is exhausted.
// This is a function (not a method) because methods cannot have their own type parameters.
//
// SIZE: Calculated as min(e1Size, e2Size) if both sizes are known, else unknown.
//
// Example:
//
//	numbers := []int{1, 2, 3}
//	letters := []string{"a", "b", "c"}
//	zipped := Zip(
//	    From(numbers),
//	    From(letters),
//	    func(n int, s string) string { return fmt.Sprintf("%d:%s", n, s) },
//	).ToSlice()
//	// ["1:a", "2:b", "3:c"]
func Zip[T1, T2, R any](e1 Enumerable[T1], e2 Enumerable[T2], resultSelector func(T1, T2) R) Stream[R] {
	var size = -1
	if sizable1, ok1 := e1.(Sizable[T1]); ok1 {
		if sizable2, ok2 := e2.(Sizable[T2]); ok2 {
			if s1, known1 := sizable1.Size(); known1 {
				if s2, known2 := sizable2.Size(); known2 {
					// Take minimum of both sizes
					if s1 < s2 {
						size = s1
					} else {
						size = s2
					}
				}
			}
		}
	}

	return &stream[R]{
		sourceFactory: func() func() (R, bool) {
			return func() (R, bool) {
				val1, ok1 := e1.Next()
				if !ok1 {
					var zero R
					return zero, false
				}

				val2, ok2 := e2.Next()
				if !ok2 {
					var zero R
					return zero, false
				}

				return resultSelector(val1, val2), true
			}
		},
		size: size, // CALCULATED: min(e1Size, e2Size) if both known
	}
}
