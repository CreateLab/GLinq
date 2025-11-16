package glinq

// KeyValue represents a key-value pair.
type KeyValue[K comparable, V any] struct {
	Key   K
	Value V
}

// FromMap creates a Stream from a map.
func FromMap[K comparable, V any](m map[K]V) Stream[KeyValue[K, V]] {
	// Convert map to slice of pairs
	pairs := make([]KeyValue[K, V], 0, len(m)) // pre-allocate capacity
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

// Keys extracts only keys from Enumerable[KeyValue].
func Keys[K comparable, V any](enum Enumerable[KeyValue[K, V]]) Stream[K] {
	return &stream[K]{
		source: func() (K, bool) {
			kv, ok := enum.Next()
			if !ok {
				var zero K
				return zero, false
			}
			return kv.Key, true
		},
	}
}

// Values extracts only values from Enumerable[KeyValue].
func Values[K comparable, V any](enum Enumerable[KeyValue[K, V]]) Stream[V] {
	return &stream[V]{
		source: func() (V, bool) {
			kv, ok := enum.Next()
			if !ok {
				var zero V
				return zero, false
			}
			return kv.Value, true
		},
	}
}

// ToMap materializes Enumerable[KeyValue] back into a map.
func ToMap[K comparable, V any](enum Enumerable[KeyValue[K, V]]) map[K]V {
	result := make(map[K]V)
	for {
		kv, ok := enum.Next()
		if !ok {
			break
		}
		result[kv.Key] = kv.Value
	}
	return result
}

// GroupBy groups elements by a key selector and returns a Stream of KeyValue pairs.
// Each KeyValue contains a key and a slice of elements that have that key.
// This is a function (not a method) because in Go methods cannot have their own type parameters.
//
// Example:
//
//	type Person struct { Age int; Name string }
//	people := []Person{{25, "Alice"}, {30, "Bob"}, {25, "Charlie"}}
//	grouped := GroupBy(
//	    From(people),
//	    func(p Person) int { return p.Age },
//	).ToSlice()
//	// []KeyValue[int, []Person]{
//	//   {Key: 25, Value: []Person{{25, "Alice"}, {25, "Charlie"}}},
//	//   {Key: 30, Value: []Person{{30, "Bob"}}},
//	// }
func GroupBy[T any, K comparable](enum Enumerable[T], keySelector func(T) K) Stream[KeyValue[K, []T]] {
	// Materialize groups into a map
	groups := make(map[K][]T)
	for {
		elem, ok := enum.Next()
		if !ok {
			break
		}
		key := keySelector(elem)
		groups[key] = append(groups[key], elem)
	}

	// Convert map to slice of KeyValue pairs
	pairs := make([]KeyValue[K, []T], 0, len(groups))
	for key, values := range groups {
		pairs = append(pairs, KeyValue[K, []T]{Key: key, Value: values})
	}

	index := 0
	return &stream[KeyValue[K, []T]]{
		source: func() (KeyValue[K, []T], bool) {
			if index < len(pairs) {
				result := pairs[index]
				index++
				return result, true
			}
			return KeyValue[K, []T]{}, false
		},
	}
}
