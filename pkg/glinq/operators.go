package glinq

// Where фильтрует элементы по предикату.
func (s *Stream[T]) Where(predicate func(T) bool) *Stream[T] {
	oldSource := s.source
	return &Stream[T]{
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

// Select преобразует элементы в тот же тип.
// Поддерживает цепочки вызовов (chaining).
//
// Пример:
//   doubled := From([]int{1, 2, 3}).
//       Select(func(x int) int { return x * 2 }).
//       ToSlice()
//   // []int{2, 4, 6}
func (s *Stream[T]) Select(mapper func(T) T) *Stream[T] {
	oldSource := s.source
	return &Stream[T]{
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

// Map преобразует элементы в другой тип.
// Это функция (не метод), так как в Go методы не могут иметь собственные type parameters.
//
// Пример:
//   strings := Map(
//       From([]int{1, 2, 3}),
//       func(x int) string { return fmt.Sprintf("num_%d", x) },
//   ).ToSlice()
//   // []string{"num_1", "num_2", "num_3"}
func Map[T, R any](s *Stream[T], mapper func(T) R) *Stream[R] {
	oldSource := s.source
	return &Stream[R]{
		source: func() (R, bool) {
			value, ok := oldSource()
			if !ok {
				var zero R
				return zero, false
			}
			return mapper(value), true
		},
	}
}

// Take берет первые n элементов из Stream.
func (s *Stream[T]) Take(n int) *Stream[T] {
	oldSource := s.source
	count := 0
	return &Stream[T]{
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

// Skip пропускает первые n элементов из Stream.
func (s *Stream[T]) Skip(n int) *Stream[T] {
	oldSource := s.source
	skipped := 0
	return &Stream[T]{
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
