// Package glinq provides a LINQ-like API for Go with support for lazy evaluation.
//
// glinq allows working with slices and maps using a functional programming style.
// All operations (Where, Select, Take, Skip) are executed lazily and do not start until a terminal
// operation (ToSlice, First, Count, Any, All, ForEach) is called.
//
// Thread Safety:
// Stream operations are NOT thread-safe. A Stream should not be used concurrently
// by multiple goroutines without external synchronization. However, each Stream
// operation returns a new Stream instance, so you can safely use different Stream
// instances in different goroutines. Modifying the underlying data structure
// (slice or map) while iterating may lead to undefined behavior.
// For concurrent modifications, use FromSafe() or FromMapSafe() to create isolated snapshots.
//
// Example usage:
//
//	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
//	result := glinq.From(numbers).
//		Where(func(x int) bool { return x > 5 }).
//		Select(func(x int) int { return x * 2 }).
//		ToSlice()
//	// result: [12, 14, 16, 18, 20]
//
// Supported operations:
//
// Creators (create Stream):
//   - From: from a slice
//   - Empty: empty Stream
//   - Range: stream of integers
//   - FromMap: from a map (returns KeyValue pairs)
//
// Operators (transform Stream):
//   - Where: filter by predicate
//   - Select: transform elements
//   - Take: first n elements
//   - TakeWhile: take elements while predicate returns true
//   - Skip: skip first n elements
//   - SkipWhile: skip elements while predicate returns true
//   - Reverse: reverse order of elements (materializes stream)
//   - SelectMany: flatten sequences (function, not method)
//   - GroupBy: group elements by key (function, returns KeyValue pairs)
//
// Terminal operations (materialize result):
//   - ToSlice: convert to slice
//   - First: first element
//   - Count: number of elements
//   - Any: check if any element exists
//   - All: check if all elements satisfy condition
//   - ForEach: execute action for each element
//
// Helper functions for working with KeyValue:
//   - Keys: extract keys
//   - Values: extract values
//   - ToMap: convert to map
//   - GroupBy: group elements by key selector
package glinq
