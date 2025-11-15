package glinq

// KeyValue represents a key-value pair.
type KeyValue[K comparable, V any] struct {
	Key   K
	Value V
}

// FromMap creates a Stream from a map.
func FromMap[K comparable, V any](m map[K]V) Stream[KeyValue[K, V]] {
	// Convert map to slice of pairs
	pairs := make([]KeyValue[K, V], 0, len(m)) // предварительная аллокация
	for key, value := range m {
		pairs = append(pairs, KeyValue[K, V]{Key: key, Value: value})
	}

	index := 0
	return &stream[KeyValue[K, V]]{
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

// Keys extracts only keys from Stream[KeyValue].
func Keys[K comparable, V any](s Stream[KeyValue[K, V]]) Stream[K] {
	return &stream[K]{
		source: func() (K, bool) {
			kv, ok := s.Next()
			if !ok {
				var zero K
				return zero, false
			}
			return kv.Key, true
		},
	}
}

// Values extracts only values from Stream[KeyValue].
func Values[K comparable, V any](s Stream[KeyValue[K, V]]) Stream[V] {
	return &stream[V]{
		source: func() (V, bool) {
			kv, ok := s.Next()
			if !ok {
				var zero V
				return zero, false
			}
			return kv.Value, true
		},
	}
}

// ToMap materializes Stream[KeyValue] back into a map.
func ToMap[K comparable, V any](s Stream[KeyValue[K, V]]) map[K]V {
	result := make(map[K]V)
	for {
		kv, ok := s.Next()
		if !ok {
			break
		}
		result[kv.Key] = kv.Value
	}
	return result
}
