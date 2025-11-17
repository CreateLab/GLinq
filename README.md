# glinq

LINQ-like API for Go with support for lazy evaluation.

## Features

- **Lazy Evaluation**: All intermediate operations are executed only when the result is materialized
- **Type Safe**: Full support for generics (Go 1.18+)
- **Composable**: Operations can be easily combined into chains
- **Zero Dependencies**: No external dependencies required
- **Extensible**: Works with any type implementing `Enumerable` interface

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

## Basic Examples

### Filtering and Transformation

```go
// Filter even numbers
evens := glinq.From([]int{1, 2, 3, 4, 5}).
    Where(func(x int) bool { return x%2 == 0 }).
    ToSlice()
// [2, 4]

// Transform to strings
strings := glinq.Select(
    glinq.From([]int{1, 2, 3}),
    func(x int) string { return fmt.Sprintf("num_%d", x) },
).ToSlice()
// []string{"num_1", "num_2", "num_3"}
```

### Sorting and Ordering

```go
// Sort ascending
sorted := glinq.From([]int{5, 2, 8, 1, 9}).
    OrderBy(func(a, b int) int { return a - b }).
    ToSlice()
// [1, 2, 5, 8, 9]

// Get top 3 smallest
top3 := glinq.From([]int{5, 2, 8, 1, 9, 3}).
    TakeOrdered(3).
    ToSlice()
// [1, 2, 3]
```

### Removing Duplicates

```go
// Remove duplicates
unique := glinq.Distinct(glinq.From([]int{1, 2, 2, 3, 3, 4})).ToSlice()
// [1, 2, 3, 4]

// Remove duplicates by key
type Person struct { ID int; Name string }
people := []Person{{1, "Alice"}, {1, "Alice2"}, {2, "Bob"}}
uniquePeople := glinq.From(people).
    DistinctBy(func(p Person) any { return p.ID }).
    ToSlice()
```

### Set Operations

```go
set1 := glinq.From([]int{1, 2, 3})
set2 := glinq.From([]int{3, 4, 5})

union := glinq.Union(set1, set2).ToSlice()        // [1, 2, 3, 4, 5]
intersect := glinq.Intersect(set1, set2).ToSlice() // [3]
except := glinq.Except(set1, set2).ToSlice()      // [1, 2]
```

## Documentation

📚 **[Full Documentation & Wiki](docs/WIKI.md)** - Complete API reference, examples, and guides

### Quick Links

- [Getting Started](docs/WIKI.md#getting-started)
- [API Reference](docs/WIKI.md#api-reference)
- [Architecture](docs/WIKI.md#architecture)
- [Examples](docs/WIKI.md#examples)
- [Best Practices](docs/WIKI.md#best-practices)

## Requirements

- Go 1.18+ (for generics support)

## Comparison with Similar Libraries

| Feature | glinq | [samber/lo](https://github.com/samber/lo) | [thoas/go-funk](https://github.com/thoas/go-funk) |
|---------|-------|-----------|----------|
| **Evaluation** | Lazy (deferred) | Eager (immediate) | Eager (immediate) |
| **API Style** | Fluent/Chainable | Functional | Functional |
| **Type Safety** | Full (generics) | Full (generics) | Runtime (reflection) |
| **Performance** | Single pass, no intermediate arrays | Creates intermediate arrays | Slower due to reflection |
| **Memory Usage** | Minimal (lazy) | Higher (eager) | Higher (eager + reflection) |
| **Extensibility** | Interface-based (Enumerable/Stream) | None | None |
| **Dependencies** | Zero | Zero | Zero |

## Testing

```bash
go test ./...
```

## Running Examples

```bash
go run examples/basic/main.go
```

## License

MIT License - see [LICENSE](LICENSE) file for details.
