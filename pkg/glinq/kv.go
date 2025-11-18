package glinq

// KeyValue represents a key-value pair.
type KeyValue[K comparable, V any] struct {
	Key   K
	Value V
}

// FromMap creates a Stream from a map.
//
// PERFORMANCE: Only keys are copied (O(n) where n is map size).
// Values are read from the map on-demand during iteration (zero-copy for values).
// This provides better performance for large maps with expensive-to-copy value types.
//
// WARNING: Modifying the map during iteration may produce unexpected results.
// If you need protection from modifications, use FromMapSafe() instead.
//
// Example:
//
//	m := map[string]int{"a": 1, "b": 2}
//	stream := FromMap(m)
//	// Efficient - only keys copied, values read on-demand!
func FromMap[K comparable, V any](m map[K]V) Stream[KeyValue[K, V]] {
	// Copy only keys to preserve iteration order
	keys := make([]K, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}

	return &stream[KeyValue[K, V]]{
		sourceFactory: func() func() (KeyValue[K, V], bool) {
			index := 0 // Fresh index for each iterator
			return func() (KeyValue[K, V], bool) {
				if index >= len(keys) {
					return KeyValue[K, V]{}, false
				}
				key := keys[index]
				// Read value from map on-demand
				value := m[key]
				index++
				return KeyValue[K, V]{Key: key, Value: value}, true
			}
		},
		size: len(m),
	}
}

// FromMapSafe creates a Stream from a map with full snapshot (defensive copying).
//
// SAFETY: Both keys and values are copied into a snapshot, so modifications
// to the original map will not affect the Stream. Use this when you need isolation.
//
// PERFORMANCE: Full snapshot can be expensive for large maps or expensive-to-copy value types.
// If performance is critical and you control the map lifecycle, use FromMap() instead.
//
// Example:
//
//	m := map[string]int{"a": 1, "b": 2}
//	stream := FromMapSafe(m)
//	m["a"] = 999 // Won't affect the stream
func FromMapSafe[K comparable, V any](m map[K]V) Stream[KeyValue[K, V]] {
	// Convert map to slice of pairs (full snapshot)
	pairs := make([]KeyValue[K, V], 0, len(m)) // pre-allocate capacity
	for key, value := range m {
		pairs = append(pairs, KeyValue[K, V]{Key: key, Value: value})
	}

	return &stream[KeyValue[K, V]]{
		sourceFactory: func() func() (KeyValue[K, V], bool) {
			index := 0 // Fresh index for each iterator
			return func() (KeyValue[K, V], bool) {
				if index >= len(pairs) {
					return KeyValue[K, V]{}, false
				}
				result := pairs[index]
				index++
				return result, true
			}
		},
		size: len(pairs),
	}
}

// Keys extracts only keys from Enumerable[KeyValue].
// SIZE: Preserves size if source is Sizable (1-to-1 transformation).
func Keys[K comparable, V any](enum Enumerable[KeyValue[K, V]]) Stream[K] {
	size := -1
	if sizable, ok := enum.(Sizable[KeyValue[K, V]]); ok {
		if s, known := sizable.Size(); known {
			size = s
		}
	}
	return &stream[K]{
		sourceFactory: func() func() (K, bool) {
			return func() (K, bool) {
				kv, ok := enum.Next()
				if !ok {
					var zero K
					return zero, false
				}
				return kv.Key, true
			}
		},
		size: size,
	}
}

// Values extracts only values from Enumerable[KeyValue].
// SIZE: Preserves size if source is Sizable (1-to-1 transformation).
func Values[K comparable, V any](enum Enumerable[KeyValue[K, V]]) Stream[V] {
	size := -1
	if sizable, ok := enum.(Sizable[KeyValue[K, V]]); ok {
		if s, known := sizable.Size(); known {
			size = s
		}
	}
	return &stream[V]{
		sourceFactory: func() func() (V, bool) {
			return func() (V, bool) {
				kv, ok := enum.Next()
				if !ok {
					var zero V
					return zero, false
				}
				return kv.Value, true
			}
		},
		size: size,
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

	// SIZE: Known after materialization
	return &stream[KeyValue[K, []T]]{
		sourceFactory: func() func() (KeyValue[K, []T], bool) {
			index := 0 // Fresh index for each iterator
			return func() (KeyValue[K, []T], bool) {
				if index >= len(pairs) {
					return KeyValue[K, []T]{}, false
				}
				result := pairs[index]
				index++
				return result, true
			}
		},
		size: len(pairs), // Actually we know it after materialization
	}
}
