package glinq

// KeyValue представляет пару ключ-значение.
type KeyValue[K comparable, V any] struct {
	Key   K
	Value V
}

// FromMap создает Stream из map.
func FromMap[K comparable, V any](m map[K]V) *Stream[KeyValue[K, V]] {
	// Преобразуем map в слайс пар
	var pairs []KeyValue[K, V]
	for key, value := range m {
		pairs = append(pairs, KeyValue[K, V]{Key: key, Value: value})
	}

	index := 0
	return &Stream[KeyValue[K, V]]{
		source: func() (KeyValue[K, V], bool) {
			if index < len(pairs) {
				result := pairs[index]
				index++
				return result, true
			}
			return KeyValue[K, V]{}, false
		},
	}
}

// Keys извлекает только ключи из Stream[KeyValue].
func Keys[K comparable, V any](stream *Stream[KeyValue[K, V]]) *Stream[K] {
	oldSource := stream.source
	return &Stream[K]{
		source: func() (K, bool) {
			kv, ok := oldSource()
			if !ok {
				var zero K
				return zero, false
			}
			return kv.Key, true
		},
	}
}

// Values извлекает только значения из Stream[KeyValue].
func Values[K comparable, V any](stream *Stream[KeyValue[K, V]]) *Stream[V] {
	oldSource := stream.source
	return &Stream[V]{
		source: func() (V, bool) {
			kv, ok := oldSource()
			if !ok {
				var zero V
				return zero, false
			}
			return kv.Value, true
		},
	}
}

// ToMap материализует Stream[KeyValue] обратно в map.
func ToMap[K comparable, V any](stream *Stream[KeyValue[K, V]]) map[K]V {
	result := make(map[K]V)
	for {
		kv, ok := stream.source()
		if !ok {
			break
		}
		result[kv.Key] = kv.Value
	}
	return result
}
