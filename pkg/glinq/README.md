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

- **Stream Creation**: `From`, `Empty`, `Range`, `FromEnumerable`
- **Filtering**: `Where`, `DistinctBy`, `Take`, `Skip`
- **Transformation**: `Select`, `SelectWithIndex`, `SelectMany`
- **Ordering**: `OrderBy`, `OrderByDescending`, `Reverse`
- **Grouping**: `GroupBy`
- **Terminal**: `ToSlice`, `First`, `Count`, `Any`, `All`, `Aggregate`, `ForEach`

See [package documentation](https://github.com/CreateLab/GLinq) for full API reference.
