# glinq

LINQ-like API for Go with lazy evaluation support.

## Quick Start

```go
import "github.com/CreateLab/glinq/pkg/glinq"

// Create stream and chain operations
result := glinq.From([]int{1, 2, 3, 4, 5}).
    Where(func(x int) bool { return x > 2 }).
    Select(func(x int) int { return x * 2 }).
    ToSlice()
// [6, 8, 10]
```

## Features

- **Lazy Evaluation** - Operations execute only when materialized
- **Type Safe** - Full generics support (Go 1.18+)
- **Composable** - Fluent method chaining
- **Zero Dependencies** - No external packages

## Main Operations

- **Stream Creation**: `From`, `Empty`, `Range`, `FromEnumerable`, `FromMap`
- **Filtering**: `Where`, `DistinctBy`, `Take`, `TakeWhile`, `Skip`, `SkipWhile`
- **Transformation**: `Select`, `SelectWithIndex`, `SelectMany`
- **Ordering**: `OrderBy`, `OrderByDescending`, `Reverse`
- **Grouping**: `GroupBy`
- **Combining**: `Zip` - combine two sequences using result selector
- **Terminal**: `ToSlice`, `First`, `Last`, `ElementAt`, `ElementAtOrDefault`, `Contains`, `ContainsBy`, `Count`, `Any`, `AnyMatch`, `All`, `Aggregate`, `ForEach`
- **Size Information**: `Size()` - returns known size for performance optimizations

## Performance Optimizations

glinq automatically tracks size information when possible:

- **`Count()`** - O(1) when size is known
- **`Any()`** - O(1) when size is known  
- **`ToSlice()`** - Preallocates capacity when size is known
- **`Chunk()`** - Preallocates result capacity when size is known

Size is preserved through transformations (`Select`, `Take`, `Skip`) and lost through filters (`Where`, `DistinctBy`).

See [package documentation](https://github.com/CreateLab/GLinq) for full API reference.
