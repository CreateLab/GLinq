package glinq

// Distinct removes duplicates from Stream.
// T must be comparable, otherwise code will not compile.
// This is a function (not a method) because methods cannot have their own type constraints.
//
// SIZE: Loses size (unknown how many duplicates exist).
func Distinct[T comparable](enum Enumerable[T]) Stream[T] {
	return &stream[T]{
		sourceFactory: func() func() (T, bool) {
			seen := make(map[T]bool) // Fresh map for each iterator

			return func() (T, bool) {
				for {
					val, ok := enum.Next()
					if !ok {
						var zero T
						return zero, false
					}

					if !seen[val] {
						seen[val] = true
						return val, true
					}
				}
			}
		},
		size: nil, // LOSE: unknown how many duplicates
	}
}

// DistinctBy removes duplicates by key extracted by keySelector.
//
// SIZE: Loses size (unknown how many duplicates exist).
func (s *stream[T]) DistinctBy(keySelector func(T) any) Stream[T] {
	return &stream[T]{
		sourceFactory: func() func() (T, bool) {
			source := s.sourceFactory() // Get fresh source
			seen := make(map[any]bool)  // Fresh map for each iterator

			return func() (T, bool) {
				for {
					val, ok := source()
					if !ok {
						var zero T
						return zero, false
					}

					key := keySelector(val)
					if !seen[key] {
						seen[key] = true
						return val, true
					}
				}
			}
		},
		size: nil, // LOSE: unknown how many duplicates
	}
}
