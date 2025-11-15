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

### Transformation (Select)

```go
numbers := []int{1, 2, 3}
squared := glinq.From(numbers).
    Select(func(x int) int { return x * x }).
    ToSlice()
// [1, 4, 9]
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

## Supported Operations

### Creators

- `From[T any](slice []T) Stream[T]` - create Stream from slice
- `Empty[T any]() Stream[T]` - create empty Stream
- `Range(start, count int) Stream[int]` - create Stream of integers
- `FromMap[K, V](m map[K]V) Stream[KeyValue[K, V]]` - create Stream from map

### Operators

- `Where(predicate func(T) bool) Stream[T]` - filter by condition
- `Select[R any](mapper func(T) R) Stream[R]` - transform elements
- `Take(n int) Stream[T]` - take first n elements
- `Skip(n int) Stream[T]` - skip first n elements

### Terminal Operations

- `ToSlice() []T` - convert Stream to slice
- `First() (T, bool)` - get first element
- `Count() int` - count number of elements
- `Any(predicate func(T) bool) bool` - check if any element exists
- `All(predicate func(T) bool) bool` - check if all elements satisfy condition
- `ForEach(action func(T))` - execute action for each element

### Helper Functions

- `Keys[K, V](stream Stream[KeyValue[K, V]]) Stream[K]` - extract keys
- `Values[K, V](stream Stream[KeyValue[K, V]]) Stream[V]` - extract values
- `ToMap[K, V](stream Stream[KeyValue[K, V]]) map[K]V` - convert to map

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

## Author

your-username
