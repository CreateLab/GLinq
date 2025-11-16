package glinq

// Distinct removes duplicates from Stream.
// T must be comparable, otherwise code will not compile.
// This is a function (not a method) because methods cannot have their own type constraints.
func Distinct[T comparable](enum Enumerable[T]) Stream[T] {
	seen := make(map[T]bool)

	return &stream[T]{
		source: func() (T, bool) {
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
		},
	}
}

// DistinctBy removes duplicates by key extracted by keySelector.
func (s *stream[T]) DistinctBy(keySelector func(T) any) Stream[T] {
	seen := make(map[any]bool)

	return &stream[T]{
		source: func() (T, bool) {
			for {
				val, ok := s.source()
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
		},
	}
}
