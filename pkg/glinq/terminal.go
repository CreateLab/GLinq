package glinq

// ToSlice материализует Stream в слайс.
func (s *Stream[T]) ToSlice() []T {
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

// First возвращает первый элемент и true, или zero value и false если Stream пустой.
func (s *Stream[T]) First() (T, bool) {
	return s.source()
}

// Count возвращает количество элементов в Stream.
func (s *Stream[T]) Count() int {
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

// Any проверяет, есть ли хотя бы один элемент, удовлетворяющий предикату.
func (s *Stream[T]) Any(predicate func(T) bool) bool {
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

// All проверяет, все ли элементы удовлетворяют предикату.
func (s *Stream[T]) All(predicate func(T) bool) bool {
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

// ForEach выполняет действие для каждого элемента в Stream.
func (s *Stream[T]) ForEach(action func(T)) {
	for {
		value, ok := s.source()
		if !ok {
			break
		}
		action(value)
	}
}

// ToMapBy материализует Stream[T] в map, используя селекторы для ключа и значения.
//
// Пример:
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
	stream *Stream[T],
	keySelector func(T) K,
	valueSelector func(T) V,
) map[K]V {
	result := make(map[K]V)
	for {
		item, ok := stream.source()
		if !ok {
			break
		}
		key := keySelector(item)
		value := valueSelector(item)
		result[key] = value
	}
	return result
}
