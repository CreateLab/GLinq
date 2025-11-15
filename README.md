# glinq

LINQ-like API for Go with support for lazy evaluation.

glinq provides a functional approach to working with slices and maps in Go, inspired by Microsoft LINQ.
All operations are executed lazily and do not start until a terminal operation is called.

## Features

- **Lazy Evaluation**: All intermediate operations are executed only when the result is materialized
- **Type Safe**: Full support for generics (Go 1.18+)
- **Composable**: Operations can be easily combined into chains
- **Zero Dependencies**: No external dependencies required
- **Map Support**: Built-in support for working with maps

## Installation

```bash
go get github.com/CreateLab/glinq
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/CreateLab/glinq/pkg/glinq"
)

func main() {
    numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
    result := glinq.From(numbers).
        Where(func(x int) bool { return x > 5 }).
        Select(func(x int) int { return x * 2 }).
        ToSlice()
    
    fmt.Println(result) // [12, 14, 16, 18, 20]
}
```

## Usage Examples

### Filtering (Where)

```go
numbers := []int{1, 2, 3, 4, 5}
evens := glinq.From(numbers).
    Where(func(x int) bool { return x%2 == 0 }).
    ToSlice()
// [2, 4]
```

### Transformation

glinq provides **two ways** to transform elements:

#### Select (Method) - Same Type Transformation

The `Select` **method** transforms elements to the same type and supports method chaining:

```go
numbers := []int{1, 2, 3}
squared := glinq.From(numbers).
    Select(func(x int) int { return x * x }).
    ToSlice()
// [1, 4, 9]
```

#### Select (Function) - Different Type Transformation

The `Select` **function** transforms elements to a different type. It's a standalone function (not a method) because Go methods cannot have their own type parameters:

```go
numbers := []int{1, 2, 3}
strings := glinq.Select(
    glinq.From(numbers),
    func(x int) string { return fmt.Sprintf("num_%d", x) },
).ToSlice()
// []string{"num_1", "num_2", "num_3"}
```

#### SelectWithIndex (Method) - Same Type Transformation with Index

The `SelectWithIndex` **method** transforms elements to the same type, providing the element index to the mapper function:

```go
numbers := []int{1, 2, 3}
result := glinq.From(numbers).
    SelectWithIndex(func(x int, idx int) int { return x * idx }).
    ToSlice()
// []int{0, 2, 6}
```

#### SelectWithIndex (Function) - Different Type Transformation with Index

The `SelectWithIndex` **function** transforms elements to a different type, providing the element index to the mapper function:

```go
numbers := []int{1, 2, 3}
strings := glinq.SelectWithIndex(
    glinq.From(numbers),
    func(x int, idx int) string { return fmt.Sprintf("num_%d_at_%d", x, idx) },
).ToSlice()
// []string{"num_1_at_0", "num_2_at_1", "num_3_at_2"}
```

### Limiting Elements (Take and Skip)

```go
numbers := []int{1, 2, 3, 4, 5}

// Take first 3 elements
first3 := glinq.From(numbers).Take(3).ToSlice()
// [1, 2, 3]

// Skip first 2 elements
rest := glinq.From(numbers).Skip(2).ToSlice()
// [3, 4, 5]
```

### Working with Maps

```go
data := map[string]int{
    "apple":  5,
    "banana": 3,
    "orange": 8,
}

filtered := glinq.FromMap(data).
    Where(func(kv glinq.KeyValue[string, int]) bool {
        return kv.Value > 4
    }).
    ToMap()
// map[apple:5 orange:8]
```

### Condition Checking (Any and All)

```go
numbers := []int{1, 2, 3, 4, 5}

hasEven := glinq.From(numbers).Any(func(x int) bool { 
    return x%2 == 0 
})
// true

allPositive := glinq.From(numbers).All(func(x int) bool { 
    return x > 0 
})
// true
```

### Counting Elements (Count)

```go
numbers := []int{1, 2, 3, 4, 5}
count := glinq.From(numbers).
    Where(func(x int) bool { return x > 2 }).
    Count()
// 3
```

### Executing Action for Each Element (ForEach)

```go
numbers := []int{1, 2, 3}
glinq.From(numbers).ForEach(func(x int) {
    fmt.Println(x)
})
// 1
// 2
// 3
```

### Getting First Element (First)

```go
numbers := []int{1, 2, 3}
first, ok := glinq.From(numbers).First()
// first = 1, ok = true
```

### Splitting into Chunks (Chunk)

```go
numbers := []int{1, 2, 3, 4, 5, 6, 7}
chunks := glinq.From(numbers).Chunk(3)
// [][]int{{1, 2, 3}, {4, 5, 6}, {7}}
```

### Getting Last Element (Last)

```go
numbers := []int{1, 2, 3, 4, 5}
last, ok := glinq.From(numbers).Last()
// last = 5, ok = true
```

### Summing Elements (Sum)

```go
numbers := []int{1, 2, 3, 4, 5}
sum := glinq.Sum(glinq.From(numbers))
// 15
```

### Finding Minimum (Min)

glinq provides **two ways** to find minimum:

#### Min (Function) - For Ordered Types

The `Min` **function** works with ordered types (int, uint, float, string):

```go
numbers := []int{5, 2, 8, 1, 9}
min, ok := glinq.Min(glinq.From(numbers))
// min = 1, ok = true
```

#### Min (Method) - With Comparator

The `Min` **method** works with any type using a comparator function:

```go
type Person struct { Age int; Name string }
people := []Person{{Age: 30, Name: "Alice"}, {Age: 25, Name: "Bob"}}
youngest, ok := glinq.From(people).Min(func(a, b Person) int {
    return a.Age - b.Age
})
// youngest = Person{Age: 25, Name: "Bob"}, ok = true
```

### Aggregating Elements (Aggregate)

The `Aggregate` method applies an accumulator function over the Stream. The seed parameter is the initial accumulator value:

```go
numbers := []int{1, 2, 3, 4, 5}
sum := glinq.From(numbers).Aggregate(0, func(acc, x int) int { return acc + x })
// 15

numbers := []int{2, 3, 4}
product := glinq.From(numbers).Aggregate(1, func(acc, x int) int { return acc * x })
// 24

words := []string{"Hello", " ", "World", "!"}
concatenated := glinq.From(words).Aggregate("", func(acc, x string) string { return acc + x })
// "Hello World!"
```

### Finding Maximum (Max)

glinq provides **two ways** to find maximum:

#### Max (Function) - For Ordered Types

The `Max` **function** works with ordered types (int, uint, float, string):

```go
numbers := []int{5, 2, 8, 1, 9}
max, ok := glinq.Max(glinq.From(numbers))
// max = 9, ok = true
```

#### Max (Method) - With Comparator

The `Max` **method** works with any type using a comparator function:

```go
type Person struct { Age int; Name string }
people := []Person{{Age: 30, Name: "Alice"}, {Age: 25, Name: "Bob"}}
oldest, ok := glinq.From(people).Max(func(a, b Person) int {
    return a.Age - b.Age
})
// oldest = Person{Age: 30, Name: "Alice"}, ok = true
```

### Lazy Evaluation Demonstration

```go
// Thanks to lazy evaluation, the filter is applied only to necessary elements
result := glinq.Range(1, 1000000).
    Where(func(x int) bool { return x%2 == 0 }).
    Take(3).
    ToSlice()
// [2, 4, 6]
// Only ~6 elements processed, not a million!
```

## API Reference

glinq provides both **methods** (on `Stream[T]` interface) and **standalone functions**. Methods support method chaining, while functions are used when type parameters are needed.

### Creator Functions

These functions create a new `Stream[T]`:

- `From[T any](slice []T) Stream[T]` - create Stream from slice
- `Empty[T any]() Stream[T]` - create empty Stream
- `Range(start, count int) Stream[int]` - create Stream of integers
- `FromMap[K, V](m map[K]V) Stream[KeyValue[K, V]]` - create Stream from map

### Stream Methods (Operators)

These methods transform the Stream and return a new `Stream[T]`:

- `Where(predicate func(T) bool) Stream[T]` - filter by condition
- `Select(mapper func(T) T) Stream[T]` - transform elements to the same type
- `SelectWithIndex(mapper func(T, int) T) Stream[T]` - transform elements to the same type with index
- `Take(n int) Stream[T]` - take first n elements
- `Skip(n int) Stream[T]` - skip first n elements

### Stream Methods (Terminal Operations)

These methods materialize the Stream:

- `ToSlice() []T` - convert Stream to slice
- `Chunk(size int) [][]T` - split Stream into chunks of specified size
- `First() (T, bool)` - get first element
- `Last() (T, bool)` - get last element
- `Count() int` - count number of elements
- `Any(predicate func(T) bool) bool` - check if any element exists
- `All(predicate func(T) bool) bool` - check if all elements satisfy condition
- `ForEach(action func(T))` - execute action for each element
- `Min(comparator func(T, T) int) (T, bool)` - find minimum element using comparator (works with any type)
- `Max(comparator func(T, T) int) (T, bool)` - find maximum element using comparator (works with any type)
- `Aggregate(seed T, accumulator func(T, T) T) T` - apply accumulator function over Stream

### Transformation Functions

These standalone functions transform Stream to different types:

- `Select[T, R any](s Stream[T], mapper func(T) R) Stream[R]` - transform elements to a different type (function version)
- `SelectWithIndex[T, R any](s Stream[T], mapper func(T, int) R) Stream[R]` - transform elements to a different type with index (function version)

### Map Helper Functions

These functions work with `Stream[KeyValue[K, V]]`:

- `Keys[K, V](s Stream[KeyValue[K, V]]) Stream[K]` - extract keys from KeyValue pairs
- `Values[K, V](s Stream[KeyValue[K, V]]) Stream[V]` - extract values from KeyValue pairs
- `ToMap[K, V](s Stream[KeyValue[K, V]]) map[K]V` - convert Stream[KeyValue] to map

### Terminal Functions

These standalone functions materialize Streams:

- `ToMapBy[T, K, V](s Stream[T], keySelector func(T) K, valueSelector func(T) V) map[K]V` - convert Stream[T] to map using selectors

### Numeric Functions

These functions work with numeric and ordered types:

- `Sum[T Numeric](s Stream[T]) T` - calculate sum of all elements (works with int, uint, float types)
- `Min[T Ordered](s Stream[T]) (T, bool)` - find minimum element for ordered types (works with int, uint, float, string types)
- `Max[T Ordered](s Stream[T]) (T, bool)` - find maximum element for ordered types (works with int, uint, float, string types)

**Note**: For custom types or complex comparisons, use the `Min` and `Max` methods with comparator functions instead.

## Requirements

- Go 1.18+ (for generics support)

## Testing

```bash
go test ./...
```

## Running Examples

```bash
go run examples/basic/main.go
```

## License

MIT
