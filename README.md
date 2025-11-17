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

## Performance Characteristics

glinq is optimized for performance with zero-copy defaults, following C# LINQ's approach.

### Stream Creation Performance

#### `From()` - Zero-Copy (O(1))

`From()` creates a stream instantly without copying data - it holds a reference to the original slice:

```go
data := []int{1, 2, 3, /* ...million elements */ }
stream := From(data) // O(1) - instant, no copying!
```

**Characteristics:**
- **Time Complexity:** O(1) - constant time creation
- **Space Complexity:** O(1) - no additional memory allocation
- **Behavior:** Modifications to the original slice are visible during iteration
- **Use Case:** Default choice for maximum performance (safe in 99% of cases)

#### `FromSafe()` - Defensive Copy (O(n))

`FromSafe()` creates an isolated snapshot by copying the entire slice:

```go
data := []int{1, 2, 3, 4, 5}
stream := FromSafe(data) // O(n) - copies all elements
data[0] = 999            // Won't affect stream
```

**Characteristics:**
- **Time Complexity:** O(n) - linear time for copying
- **Space Complexity:** O(n) - full slice copy in memory
- **Behavior:** Completely isolated from original slice modifications
- **Use Case:** When you need protection from concurrent modifications or want isolation

#### `FromMap()` - Keys Only Copy (O(n) keys, O(1) values)

`FromMap()` copies only keys and reads values on-demand from the map:

```go
m := map[string]int{"a": 1, "b": 2, /* ...thousands of entries */ }
stream := FromMap(m) // Copies keys only, reads values on-demand
```

**Characteristics:**
- **Time Complexity:** O(n) for key copying, O(1) per value read
- **Space Complexity:** O(n) for keys only (not values)
- **Behavior:** Values are read from map during iteration (modifications visible)
- **Use Case:** Optimal for large maps with expensive-to-copy value types

#### `FromMapSafe()` - Full Snapshot (O(n))

`FromMapSafe()` creates a complete snapshot of all key-value pairs:

```go
m := map[string]int{"a": 1, "b": 2}
stream := FromMapSafe(m) // Full snapshot
m["a"] = 999             // Won't affect stream
```

**Characteristics:**
- **Time Complexity:** O(n) - copies all key-value pairs
- **Space Complexity:** O(n) - full map snapshot in memory
- **Behavior:** Completely isolated from original map modifications
- **Use Case:** When you need complete isolation from map changes

### Benchmark Comparison

Run benchmarks to see the performance difference:

```bash
go test -bench=BenchmarkFrom -benchmem ./pkg/glinq/...
```

**Expected Results:**
- `From()`: Constant time regardless of slice size
- `FromSafe()`: Linear time scaling with slice size
- `FromMap()`: Faster than `FromMapSafe()` for large maps with expensive value types

### Recommendations

1. **Use `From()` by default** - Maximum performance, safe in most cases
2. **Use `FromSafe()`** - Only when you explicitly need isolation
3. **Use `FromMap()`** - Default for maps, especially with large or expensive value types
4. **Use `FromMapSafe()`** - Only when you need complete map isolation

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

## Performance Benchmarks

### TL;DR  
**glinq is 10-1000x faster than go-funk, 2-5x faster than go-linq, with 10x better memory efficiency than samber/lo for complex chains.**

### Quick Comparison
| Scenario | glinq | samber/lo | go-linq | go-funk |
|----------|-------|-----------|---------|---------|
| Filter+Map+First (1M items) | ✅ **0.4μs** | 1.4ms | 1.8μs | 337ms |
| Complex Chain (1M items) | ✅ **0.6μs** | 2.5ms | 0.9μs | 478ms |
| Memory Usage | ✅ **280B** | 16MB | 992B | 188MB |
| Complex Chain Memory | ✅ **600B** | 24MB | 552B | 214MB |

### Library Recommendations

#### 🟢 glinq
**Best for:** Complex query chains, large datasets, memory-sensitive applications  
**Strengths:** 10x less memory than samber/lo in chains, minimal allocations, lazy evaluation  
**Example:** `From(data).Where(...).Select(...).Take(10)`


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
