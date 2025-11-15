package glinq

// Stream представляет ленивую последовательность элементов.
// Элементы передаются через функцию source, которая возвращает (элемент, есть_ли_элемент).
type Stream[T any] struct {
	source func() (T, bool)
}

// From создает Stream из слайса.
// Слайс копируется, так что изменения к исходному слайсу не влияют на Stream.
func From[T any](slice []T) *Stream[T] {
	// Копируем слайс, чтобы избежать проблем с изменениями исходного слайса
	data := make([]T, len(slice))
	copy(data, slice)

	index := 0
	return &Stream[T]{
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

// Empty создает пустой Stream, который не содержит элементов.
func Empty[T any]() *Stream[T] {
	return &Stream[T]{
		source: func() (T, bool) {
			var zero T
			return zero, false
		},
	}
}

// Range создает Stream целых чисел от start до start+count-1.
func Range(start, count int) *Stream[int] {
	index := 0
	return &Stream[int]{
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
