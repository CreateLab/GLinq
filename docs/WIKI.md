# glinq Wiki - Complete Documentation

## Table of Contents

- [Getting Started](#getting-started)
- [Architecture](#architecture)
- [API Reference](#api-reference)
- [Examples](#examples)
- [Best Practices](#best-practices)
- [Performance Considerations](#performance-considerations)

---

## Getting Started

### Installation

```bash
go get github.com/CreateLab/glinq
```

### Basic Usage

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

### Key Concepts

1. **Lazy Evaluation**: Operations don't execute until a terminal operation is called
2. **Method Chaining**: Most operations return `Stream[T]` for fluent chaining
3. **Type Safety**: Full generics support ensures compile-time type checking

---

## Architecture

glinq separates interfaces into two levels:

### Enumerable[T]

The minimal interface for iterable collections. Any type can implement it for integration with glinq:

```go
type Enumerable[T any] interface {
    Next() (T, bool)
}
```

### Stream[T]

Extends `Enumerable` and adds all operators (Where, Select, OrderBy, etc.):

```go
type Stream[T any] interface {
    Enumerable[T]
    Where(predicate func(T) bool) Stream[T]
    Select(mapper func(T) T) Stream[T]
    // ... more methods
}
```

### Why Two Interfaces?

- **Functions accept `Enumerable`**: Makes them universal and work with any iterator
- **Methods on `Stream`**: Provide convenient chaining syntax
- **Custom Sources**: You can implement `Enumerable` for your own data sources

### Creating Custom Enumerables

```go
type MyIterator struct {
    data []int
    pos  int
}

func (m *MyIterator) Next() (int, bool) {
    if m.pos < len(m.data) {
        val := m.data[m.pos]
        m.pos++
        return val, true
    }
    return 0, false
}

// Now you can use it with glinq functions
iter := &MyIterator{data: []int{1, 2, 3}}
result := glinq.Select(iter, func(x int) int { return x * 2 }).ToSlice()
```

---

## API Reference

### Creator Functions

Functions that create a new `Stream[T]`:

#### From

Creates a Stream from a slice:

```go
stream := glinq.From([]int{1, 2, 3})
```

#### Empty

Creates an empty Stream:

```go
empty := glinq.Empty[int]()
```

#### Range

Creates a Stream of integers:

```go
numbers := glinq.Range(1, 10) // [1, 2, 3, ..., 10]
```

#### FromMap

Creates a Stream from a map (returns `KeyValue` pairs):

```go
m := map[string]int{"a": 1, "b": 2}
stream := glinq.FromMap(m)
```

#### FromEnumerable

Converts any `Enumerable` to `Stream`:

```go
stream := glinq.FromEnumerable(myEnumerable)
```

---

### Stream Methods (Operators)

Methods that transform the Stream and return a new `Stream[T]`:

#### Where

Filters elements by predicate:

```go
evens := glinq.From([]int{1, 2, 3, 4, 5}).
    Where(func(x int) bool { return x%2 == 0 }).
    ToSlice()
// [2, 4]
```

#### Select

Transforms elements to the same type:

```go
squared := glinq.From([]int{1, 2, 3}).
    Select(func(x int) int { return x * x }).
    ToSlice()
// [1, 4, 9]
```

#### SelectWithIndex

Transforms elements with index:

```go
result := glinq.From([]int{1, 2, 3}).
    SelectWithIndex(func(x int, idx int) int { return x * idx }).
    ToSlice()
// [0, 2, 6]
```

#### Take

Takes first n elements:

```go
first3 := glinq.From([]int{1, 2, 3, 4, 5}).Take(3).ToSlice()
// [1, 2, 3]
```

#### Skip

Skips first n elements:

```go
rest := glinq.From([]int{1, 2, 3, 4, 5}).Skip(2).ToSlice()
// [3, 4, 5]
```

#### OrderBy

Sorts elements using comparator:

```go
sorted := glinq.From([]int{5, 2, 8, 1, 9}).
    OrderBy(func(a, b int) int { return a - b }).
    ToSlice()
// [1, 2, 5, 8, 9]
```

**Note**: `OrderBy` materializes the entire stream for sorting (partially lazy).

#### OrderByDescending

Sorts elements in reverse order:

```go
sorted := glinq.From([]int{5, 2, 8, 1, 9}).
    OrderByDescending(func(a, b int) int { return a - b }).
    ToSlice()
// [9, 8, 5, 2, 1]
```

#### DistinctBy

Removes duplicates by key selector:

```go
type Person struct { ID int; Name string }
people := []Person{{1, "Alice"}, {1, "Alice2"}, {2, "Bob"}}
unique := glinq.From(people).
    DistinctBy(func(p Person) any { return p.ID }).
    ToSlice()
// [{1, "Alice"}, {2, "Bob"}]
```

#### Concat

Concatenates two streams (preserves duplicates):

```go
result := glinq.From([]int{1, 2}).
    Concat(glinq.From([]int{2, 3})).
    ToSlice()
// [1, 2, 2, 3]
```

---

### Stream Methods (Terminal Operations)

Methods that materialize the Stream:

#### ToSlice

Converts Stream to slice:

```go
result := glinq.From([]int{1, 2, 3}).ToSlice()
```

#### First

Gets first element:

```go
first, ok := glinq.From([]int{1, 2, 3}).First()
// first = 1, ok = true
```

#### Last

Gets last element:

```go
last, ok := glinq.From([]int{1, 2, 3}).Last()
// last = 3, ok = true
```

#### Count

Counts number of elements:

```go
count := glinq.From([]int{1, 2, 3, 4, 5}).
    Where(func(x int) bool { return x > 2 }).
    Count()
// 3
```

#### Any

Checks if any element satisfies predicate:

```go
hasEven := glinq.From([]int{1, 2, 3}).Any(func(x int) bool { return x%2 == 0 })
// true
```

#### All

Checks if all elements satisfy predicate:

```go
allPositive := glinq.From([]int{1, 2, 3}).All(func(x int) bool { return x > 0 })
// true
```

#### ForEach

Executes action for each element:

```go
glinq.From([]int{1, 2, 3}).ForEach(func(x int) {
    fmt.Println(x)
})
```

#### Min / Max

Find minimum/maximum using comparator:

```go
type Person struct { Age int; Name string }
people := []Person{{30, "Alice"}, {25, "Bob"}}
youngest, _ := glinq.From(people).Min(func(a, b Person) int {
    return a.Age - b.Age
})
```

#### Aggregate

Applies accumulator function:

```go
sum := glinq.From([]int{1, 2, 3, 4, 5}).
    Aggregate(0, func(acc, x int) int { return acc + x })
// 15
```

#### Chunk

Splits Stream into chunks:

```go
chunks := glinq.From([]int{1, 2, 3, 4, 5, 6, 7}).Chunk(3)
// [][]int{{1, 2, 3}, {4, 5, 6}, {7}}
```

---

### Transformation Functions

Standalone functions that transform `Enumerable` to different types:

#### Select

Transforms to different type:

```go
strings := glinq.Select(
    glinq.From([]int{1, 2, 3}),
    func(x int) string { return fmt.Sprintf("num_%d", x) },
).ToSlice()
// []string{"num_1", "num_2", "num_3"}
```

#### SelectWithIndex

Transforms to different type with index:

```go
strings := glinq.SelectWithIndex(
    glinq.From([]int{1, 2, 3}),
    func(x int, idx int) string { return fmt.Sprintf("num_%d_at_%d", x, idx) },
).ToSlice()
```

#### Distinct

Removes duplicates (requires `comparable` T):

```go
unique := glinq.Distinct(glinq.From([]int{1, 2, 2, 3, 3, 4})).ToSlice()
// [1, 2, 3, 4]
```

#### Union

Combines two enumerables and removes duplicates:

```go
set1 := glinq.From([]int{1, 2, 3})
set2 := glinq.From([]int{3, 4, 5})
union := glinq.Union(set1, set2).ToSlice()
// [1, 2, 3, 4, 5]
```

#### Intersect

Returns intersection of two enumerables:

```go
intersect := glinq.Intersect(
    glinq.From([]int{1, 2, 3}),
    glinq.From([]int{2, 3, 4}),
).ToSlice()
// [2, 3]
```

#### Except

Returns difference of enumerables:

```go
except := glinq.Except(
    glinq.From([]int{1, 2, 3}),
    glinq.From([]int{2, 3}),
).ToSlice()
// [1]
```

#### TakeOrderedBy

Takes n smallest elements using comparator:

```go
type Person struct { Age int; Name string }
people := []Person{{30, "Alice"}, {25, "Bob"}, {35, "Charlie"}}
youngest := glinq.TakeOrderedBy(
    glinq.From(people),
    2,
    func(a, b Person) bool { return a.Age < b.Age },
).ToSlice()
```

#### TakeOrderedDescendingBy

Takes n largest elements using comparator:

```go
oldest := glinq.TakeOrderedDescendingBy(
    glinq.From(people),
    2,
    func(a, b Person) bool { return a.Age < b.Age },
).ToSlice()
```

---

### Map Helper Functions

Functions that work with `Enumerable[KeyValue[K, V]]`:

#### Keys

Extracts keys from KeyValue pairs:

```go
m := map[string]int{"a": 1, "b": 2}
keys := glinq.Keys(glinq.FromMap(m)).ToSlice()
// []string{"a", "b"}
```

#### Values

Extracts values from KeyValue pairs:

```go
values := glinq.Values(glinq.FromMap(m)).ToSlice()
// []int{1, 2}
```

#### ToMap

Converts Enumerable[KeyValue] to map:

```go
m := map[string]int{"a": 1, "b": 2}
stream := glinq.FromMap(m)
result := glinq.ToMap(stream)
```

#### ToMapBy

Converts Enumerable[T] to map using selectors:

```go
type User struct { ID int; Name string }
users := []User{{1, "Alice"}, {2, "Bob"}}
userMap := glinq.ToMapBy(
    glinq.From(users),
    func(u User) int { return u.ID },
    func(u User) string { return u.Name },
)
// map[int]string{1: "Alice", 2: "Bob"}
```

---

### Numeric Functions

Functions that work with numeric and ordered types:

#### Sum

Calculates sum of all elements:

```go
sum := glinq.Sum(glinq.From([]int{1, 2, 3, 4, 5}))
// 15
```

#### Min

Finds minimum element (for ordered types):

```go
min, ok := glinq.Min(glinq.From([]int{5, 2, 8, 1, 9}))
// min = 1, ok = true
```

#### Max

Finds maximum element (for ordered types):

```go
max, ok := glinq.Max(glinq.From([]int{5, 2, 8, 1, 9}))
// max = 9, ok = true
```

**Note**: For custom types, use `Min` and `Max` methods with comparator functions.

---

## Examples

### Complex Chain

```go
result := glinq.From([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}).
    Skip(2).
    Where(func(x int) bool { return x%2 == 0 }).
    Select(func(x int) int { return x * x }).
    Take(3).
    ToSlice()
// [16, 36, 64]
```

### Working with Structs

```go
type Person struct {
    Name string
    Age  int
}

people := []Person{
    {"Alice", 30},
    {"Bob", 25},
    {"Charlie", 35},
}

// Sort by age
sorted := glinq.From(people).
    OrderBy(func(a, b Person) int { return a.Age - b.Age }).
    ToSlice()

// Get youngest
youngest, _ := glinq.From(people).Min(func(a, b Person) int {
    return a.Age - b.Age
})
```

### Lazy Evaluation

```go
// Only processes ~6 elements, not a million!
result := glinq.Range(1, 1000000).
    Where(func(x int) bool { return x%2 == 0 }).
    Take(3).
    ToSlice()
// [2, 4, 6]
```

### Custom Enumerable

```go
type Fibonacci struct {
    a, b int
}

func (f *Fibonacci) Next() (int, bool) {
    next := f.a
    f.a, f.b = f.b, f.a+f.b
    return next, true
}

fib := &Fibonacci{a: 0, b: 1}
first10 := glinq.FromEnumerable(fib).Take(10).ToSlice()
// [0, 1, 1, 2, 3, 5, 8, 13, 21, 34]
```

---

## Best Practices

### 1. Use Method Chaining When Possible

```go
// Good
result := glinq.From(data).
    Where(predicate).
    Select(mapper).
    ToSlice()

// Less ideal (but works)
s1 := glinq.From(data)
s2 := s1.Where(predicate)
s3 := s2.Select(mapper)
result := s3.ToSlice()
```

### 2. Leverage Lazy Evaluation

```go
// Good: Only processes necessary elements
result := glinq.Range(1, 1000000).
    Where(expensivePredicate).
    Take(10).
    ToSlice()

// Less efficient: Processes all elements
all := glinq.Range(1, 1000000).
    Where(expensivePredicate).
    ToSlice()
result := all[:10]
```

### 3. Use Functions for Type Transformations

```go
// Use Select function for different types
strings := glinq.Select(
    glinq.From([]int{1, 2, 3}),
    func(x int) string { return fmt.Sprintf("%d", x) },
).ToSlice()

// Use Select method for same type
doubled := glinq.From([]int{1, 2, 3}).
    Select(func(x int) int { return x * 2 }).
    ToSlice()
```

### 4. Prefer Comparable Types for Set Operations

```go
// Good: int is comparable
unique := glinq.Distinct(glinq.From([]int{1, 2, 2, 3})).ToSlice()

// For non-comparable types, use DistinctBy
type Person struct { ID int; Name string }
unique := glinq.From(people).
    DistinctBy(func(p Person) any { return p.ID }).
    ToSlice()
```

---

## Performance Considerations

### Lazy Evaluation Benefits

- **Early Termination**: Operations stop when enough elements are collected
- **Memory Efficient**: Only processes what's needed
- **Composable**: Can chain many operations without materializing intermediate results

### When Materialization Happens

Some operations materialize the entire stream:

- `OrderBy` / `OrderByDescending` - sorts entire collection
- `Intersect` / `Except` - materializes second enumerable into a set
- `Chunk` - needs to see all elements to create chunks

### Memory Usage

- **Stream operations**: O(1) memory (lazy)
- **Set operations** (Union, Intersect, Except): O(n) memory for seen elements
- **OrderBy**: O(n) memory (materializes entire stream)

### Tips

1. Use `Take` early in chains to limit processing
2. Avoid `OrderBy` on very large streams if possible
3. Use `DistinctBy` instead of `Distinct` for non-comparable types
4. Consider materializing once if you need to iterate multiple times

---

## Type Constraints

### Comparable

Required for: `Distinct`, `Union`, `Intersect`, `Except`

```go
// Works: int is comparable
glinq.Distinct(glinq.From([]int{1, 2, 2}))

// Doesn't compile: slice is not comparable
glinq.Distinct(glinq.From([][]int{{1}, {2}})) // Error!
```

### Ordered

Required for: `Min`, `Max` functions (not methods)

```go
// Works: int is Ordered
glinq.Min(glinq.From([]int{1, 2, 3}))

// Use method with comparator for custom types
glinq.From(people).Min(func(a, b Person) int { return a.Age - b.Age })
```

### Numeric

Required for: `Sum`

```go
// Works: int is Numeric
glinq.Sum(glinq.From([]int{1, 2, 3}))

// Doesn't work: string is not Numeric
glinq.Sum(glinq.From([]string{"a", "b"})) // Error!
```

---

## Common Patterns

### Filter and Transform

```go
result := glinq.From(data).
    Where(func(x T) bool { /* filter */ }).
    Select(func(x T) T { /* transform */ }).
    ToSlice()
```

### Top N Elements

```go
top5 := glinq.From(data).
    OrderByDescending(comparator).
    Take(5).
    ToSlice()
```

### Grouping by Key

```go
// Using ToMapBy
groups := glinq.ToMapBy(
    glinq.From(items),
    func(item Item) string { return item.Category },
    func(item Item) Item { return item },
)
```

### Flattening

```go
nested := [][]int{{1, 2}, {3, 4}, {5}}
flat := glinq.From(nested).
    Aggregate([]int{}, func(acc []int, x []int) []int {
        return append(acc, x...)
    })
```

---

## Troubleshooting

### "T does not satisfy comparable"

**Problem**: Trying to use `Distinct`, `Union`, etc. with non-comparable types.

**Solution**: Use `DistinctBy` or ensure your type is comparable.

```go
// Error
glinq.Distinct(glinq.From([][]int{{1}}))

// Fix
glinq.From([][]int{{1}}).
    DistinctBy(func(x []int) any { return len(x) })
```

### "T does not satisfy Ordered"

**Problem**: Using `Min`/`Max` functions with non-ordered types.

**Solution**: Use `Min`/`Max` methods with comparator.

```go
// Error
glinq.Min(glinq.From([]Person{{Age: 30}}))

// Fix
glinq.From([]Person{{Age: 30}}).
    Min(func(a, b Person) int { return a.Age - b.Age })
```

### Performance Issues

**Problem**: Slow operations on large datasets.

**Solutions**:
- Use `Take` early to limit processing
- Avoid `OrderBy` on very large streams
- Consider materializing once if iterating multiple times

---

## Contributing

See [CONTRIBUTING.md](../CONTRIBUTING.md) for guidelines.

---

## License

MIT
