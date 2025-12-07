package glinq

import (
	"testing"
	"time"
)

func TestToSlice(t *testing.T) {
	slice := []int{1, 2, 3}
	stream := From(slice)
	result := stream.ToSlice()

	if len(result) != 3 {
		t.Errorf("expected length 3, got %d", len(result))
	}
}

func TestFirstWithElements(t *testing.T) {
	slice := []int{1, 2, 3}
	stream := From(slice)
	value, ok := stream.First()

	if !ok {
		t.Errorf("expected ok=true, got false")
	}
	if value != 1 {
		t.Errorf("expected 1, got %d", value)
	}
}

func TestFirstEmpty(t *testing.T) {
	stream := Empty[int]()
	value, ok := stream.First()

	if ok {
		t.Errorf("expected ok=false, got true")
	}
	if value != 0 {
		t.Errorf("expected 0, got %d", value)
	}
}

func TestCount(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	stream := From(slice)
	count := stream.Count()

	if count != 5 {
		t.Errorf("expected count 5, got %d", count)
	}
}

func TestCountEmpty(t *testing.T) {
	stream := Empty[int]()
	count := stream.Count()

	if count != 0 {
		t.Errorf("expected count 0, got %d", count)
	}
}

func TestCountWithFilter(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	stream := From(slice).Where(func(x int) bool { return x > 2 })
	count := stream.Count()

	if count != 3 {
		t.Errorf("expected count 3, got %d", count)
	}
}

func TestAnyTrue(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	stream := From(slice)
	result := stream.AnyMatch(func(x int) bool { return x == 3 })

	if !result {
		t.Errorf("expected true, got false")
	}
}

func TestAnyFalse(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	stream := From(slice)
	result := stream.AnyMatch(func(x int) bool { return x == 10 })

	if result {
		t.Errorf("expected false, got true")
	}
}

func TestAnyEmpty(t *testing.T) {
	stream := Empty[int]()
	result := stream.AnyMatch(func(x int) bool { return x == 1 })

	if result {
		t.Errorf("expected false, got true")
	}
}

func TestAny_WithElements(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	stream := From(slice)
	result := stream.Any()

	if !result {
		t.Errorf("expected true, got false")
	}
}

func TestAny_Empty(t *testing.T) {
	stream := Empty[int]()
	result := stream.Any()

	if result {
		t.Errorf("expected false, got true")
	}
}

func TestAny_Optimization(t *testing.T) {
	// Test with known size (should be O(1))
	stream := From([]int{1, 2, 3})
	start := time.Now()
	result := stream.Any()
	duration := time.Since(start)

	if !result {
		t.Errorf("expected true, got false")
	}
	if duration > time.Microsecond*100 {
		t.Errorf("Any with known size should be O(1), took %v", duration)
	}

	// Test with unknown size (should iterate until first element)
	stream2 := From([]int{1, 2, 3}).Where(func(x int) bool { return x > 0 })
	start2 := time.Now()
	result2 := stream2.Any()
	duration2 := time.Since(start2)

	if !result2 {
		t.Errorf("expected true, got false")
	}
	// Should still be fast for first element
	if duration2 > time.Millisecond {
		t.Errorf("Any with unknown size should be fast, took %v", duration2)
	}
}

func TestAllTrue(t *testing.T) {
	slice := []int{2, 4, 6, 8}
	stream := From(slice)
	result := stream.All(func(x int) bool { return x%2 == 0 })

	if !result {
		t.Errorf("expected true, got false")
	}
}

func TestAllFalse(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	stream := From(slice)
	result := stream.All(func(x int) bool { return x > 2 })

	if result {
		t.Errorf("expected false, got true")
	}
}

func TestAllEmpty(t *testing.T) {
	stream := Empty[int]()
	result := stream.All(func(x int) bool { return x == 1 })

	if !result {
		t.Errorf("expected true for empty stream, got false")
	}
}

func TestForEach(t *testing.T) {
	slice := []int{1, 2, 3}
	stream := From(slice)
	var result []int

	stream.ForEach(func(x int) {
		result = append(result, x*2)
	})

	expected := []int{2, 4, 6}
	if len(result) != len(expected) {
		t.Errorf("expected length %d, got %d", len(expected), len(result))
	}

	for i, v := range result {
		if v != expected[i] {
			t.Errorf("expected %d, got %d at index %d", expected[i], v, i)
		}
	}
}

func TestAggregate(t *testing.T) {
	t.Run("Sum with Aggregate", func(t *testing.T) {
		slice := []int{1, 2, 3, 4, 5}
		result := From(slice).Aggregate(0, func(acc, x int) int { return acc + x })

		if result != 15 {
			t.Errorf("expected 15, got %d", result)
		}
	})

	t.Run("Product with Aggregate", func(t *testing.T) {
		slice := []int{1, 2, 3}
		result := From(slice).Aggregate(1, func(acc, x int) int { return acc * x })

		if result != 6 {
			t.Errorf("expected 6, got %d", result)
		}
	})

	t.Run("String concatenation with Aggregate", func(t *testing.T) {
		slice := []string{"a", "b", "c"}
		result := From(slice).Aggregate("", func(acc, x string) string { return acc + x })

		if result != "abc" {
			t.Errorf("expected 'abc', got '%s'", result)
		}
	})

	t.Run("Aggregate with empty stream", func(t *testing.T) {
		stream := Empty[int]()
		result := stream.Aggregate(10, func(acc, x int) int { return acc + x })

		if result != 10 {
			t.Errorf("expected 10 (seed value), got %d", result)
		}
	})

	t.Run("Aggregate with Where filter", func(t *testing.T) {
		slice := []int{1, 2, 3, 4, 5}
		result := From(slice).
			Where(func(x int) bool { return x%2 == 0 }).
			Aggregate(0, func(acc, x int) int { return acc + x })

		if result != 6 {
			t.Errorf("expected 6 (sum of even numbers), got %d", result)
		}
	})

	t.Run("Aggregate with custom type", func(t *testing.T) {
		type Point struct {
			X, Y int
		}

		points := []Point{{1, 2}, {3, 4}, {5, 6}}
		result := From(points).Aggregate(
			Point{0, 0},
			func(acc, p Point) Point {
				return Point{acc.X + p.X, acc.Y + p.Y}
			},
		)

		if result.X != 9 || result.Y != 12 {
			t.Errorf("expected Point{X:9, Y:12}, got Point{X:%d, Y:%d}", result.X, result.Y)
		}
	})
}

func TestElementAt(t *testing.T) {
	t.Run("ElementAt with valid index", func(t *testing.T) {
		slice := []int{10, 20, 30, 40}
		stream := From(slice)
		value, ok := stream.ElementAt(2)

		if !ok {
			t.Errorf("expected ok=true, got false")
		}
		if value != 30 {
			t.Errorf("expected 30, got %d", value)
		}
	})

	t.Run("ElementAt with first element", func(t *testing.T) {
		slice := []int{10, 20, 30}
		stream := From(slice)
		value, ok := stream.ElementAt(0)

		if !ok {
			t.Errorf("expected ok=true, got false")
		}
		if value != 10 {
			t.Errorf("expected 10, got %d", value)
		}
	})

	t.Run("ElementAt with out of range index", func(t *testing.T) {
		slice := []int{10, 20, 30}
		stream := From(slice)
		value, ok := stream.ElementAt(10)

		if ok {
			t.Errorf("expected ok=false, got true")
		}
		if value != 0 {
			t.Errorf("expected 0, got %d", value)
		}
	})

	t.Run("ElementAt with negative index", func(t *testing.T) {
		slice := []int{10, 20, 30}
		stream := From(slice)
		value, ok := stream.ElementAt(-1)

		if ok {
			t.Errorf("expected ok=false, got true")
		}
		if value != 0 {
			t.Errorf("expected 0, got %d", value)
		}
	})

	t.Run("ElementAt with empty stream", func(t *testing.T) {
		stream := Empty[int]()
		value, ok := stream.ElementAt(0)

		if ok {
			t.Errorf("expected ok=false, got true")
		}
		if value != 0 {
			t.Errorf("expected 0, got %d", value)
		}
	})

	t.Run("ElementAt with filtered stream", func(t *testing.T) {
		slice := []int{1, 2, 3, 4, 5, 6, 7, 8}
		stream := From(slice).Where(func(x int) bool { return x%2 == 0 })
		value, ok := stream.ElementAt(1)

		if !ok {
			t.Errorf("expected ok=true, got false")
		}
		if value != 4 {
			t.Errorf("expected 4, got %d", value)
		}
	})
}

func TestElementAtOrDefault(t *testing.T) {
	t.Run("ElementAtOrDefault with valid index", func(t *testing.T) {
		slice := []int{10, 20, 30}
		stream := From(slice)
		value := stream.ElementAtOrDefault(1, 999)

		if value != 20 {
			t.Errorf("expected 20, got %d", value)
		}
	})

	t.Run("ElementAtOrDefault with out of range index", func(t *testing.T) {
		slice := []int{10, 20, 30}
		stream := From(slice)
		value := stream.ElementAtOrDefault(10, 999)

		if value != 999 {
			t.Errorf("expected 999, got %d", value)
		}
	})

	t.Run("ElementAtOrDefault with negative index", func(t *testing.T) {
		slice := []int{10, 20, 30}
		stream := From(slice)
		value := stream.ElementAtOrDefault(-1, 999)

		if value != 999 {
			t.Errorf("expected 999, got %d", value)
		}
	})

	t.Run("ElementAtOrDefault with empty stream", func(t *testing.T) {
		stream := Empty[int]()
		value := stream.ElementAtOrDefault(0, 999)

		if value != 999 {
			t.Errorf("expected 999, got %d", value)
		}
	})

	t.Run("ElementAtOrDefault with string default", func(t *testing.T) {
		slice := []string{"a", "b", "c"}
		stream := From(slice)
		value := stream.ElementAtOrDefault(5, "default")

		if value != "default" {
			t.Errorf("expected 'default', got '%s'", value)
		}
	})
}

func TestContains(t *testing.T) {
	t.Run("Contains with existing element", func(t *testing.T) {
		slice := []int{1, 2, 3, 4, 5}
		stream := From(slice)
		result := stream.Contains(3)

		if !result {
			t.Errorf("expected true, got false")
		}
	})

	t.Run("Contains with non-existing element", func(t *testing.T) {
		slice := []int{1, 2, 3, 4, 5}
		stream := From(slice)
		result := stream.Contains(10)

		if result {
			t.Errorf("expected false, got true")
		}
	})

	t.Run("Contains with empty stream", func(t *testing.T) {
		stream := Empty[int]()
		result := stream.Contains(1)

		if result {
			t.Errorf("expected false, got true")
		}
	})

	t.Run("Contains with strings", func(t *testing.T) {
		slice := []string{"apple", "banana", "cherry"}
		stream := From(slice)
		result := stream.Contains("banana")

		if !result {
			t.Errorf("expected true, got false")
		}
	})

	t.Run("Contains with filtered stream", func(t *testing.T) {
		slice := []int{1, 2, 3, 4, 5, 6, 7, 8}
		stream := From(slice).Where(func(x int) bool { return x%2 == 0 })
		result := stream.Contains(4)

		if !result {
			t.Errorf("expected true, got false")
		}
	})

	t.Run("Contains with struct", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}
		people := []Person{{Name: "Alice", Age: 30}, {Name: "Bob", Age: 25}}
		stream := From(people)
		result := stream.Contains(Person{Name: "Alice", Age: 30})

		if !result {
			t.Errorf("expected true, got false")
		}
	})
}

func TestContainsBy(t *testing.T) {
	t.Run("ContainsBy with ID field", func(t *testing.T) {
		type Person struct {
			ID   int
			Name string
		}
		people := []Person{{ID: 1, Name: "Alice"}, {ID: 2, Name: "Bob"}, {ID: 3, Name: "Charlie"}}
		stream := From(people)
		result := stream.ContainsBy(2, func(p Person) any { return p.ID })

		if !result {
			t.Errorf("expected true, got false")
		}
	})

	t.Run("ContainsBy with non-existing key", func(t *testing.T) {
		type Person struct {
			ID   int
			Name string
		}
		people := []Person{{ID: 1, Name: "Alice"}, {ID: 2, Name: "Bob"}}
		stream := From(people)
		result := stream.ContainsBy(999, func(p Person) any { return p.ID })

		if result {
			t.Errorf("expected false, got true")
		}
	})

	t.Run("ContainsBy with string key", func(t *testing.T) {
		type Person struct {
			ID   int
			Name string
		}
		people := []Person{{ID: 1, Name: "Alice"}, {ID: 2, Name: "Bob"}}
		stream := From(people)
		result := stream.ContainsBy("Alice", func(p Person) any { return p.Name })

		if !result {
			t.Errorf("expected true, got false")
		}
	})

	t.Run("ContainsBy with composite key", func(t *testing.T) {
		type Item struct {
			Category string
			Name     string
		}
		items := []Item{
			{Category: "Food", Name: "Apple"},
			{Category: "Drink", Name: "Water"},
		}
		stream := From(items)
		result := stream.ContainsBy("Food-Apple", func(i Item) any {
			return i.Category + "-" + i.Name
		})

		if !result {
			t.Errorf("expected true, got false")
		}
	})

	t.Run("ContainsBy with empty stream", func(t *testing.T) {
		type Person struct {
			ID   int
			Name string
		}
		stream := Empty[Person]()
		result := stream.ContainsBy(1, func(p Person) any { return p.ID })

		if result {
			t.Errorf("expected false, got true")
		}
	})

	t.Run("ContainsBy with filtered stream", func(t *testing.T) {
		type Person struct {
			ID   int
			Name string
			Age  int
		}
		people := []Person{
			{ID: 1, Name: "Alice", Age: 30},
			{ID: 2, Name: "Bob", Age: 25},
			{ID: 3, Name: "Charlie", Age: 35},
		}
		stream := From(people).Where(func(p Person) bool { return p.Age > 25 })
		result := stream.ContainsBy(3, func(p Person) any { return p.ID })

		if !result {
			t.Errorf("expected true, got false")
		}
	})

	t.Run("ContainsBy with int slice key", func(t *testing.T) {
		type Data struct {
			Values []int
		}
		data := []Data{
			{Values: []int{1, 2, 3}},
			{Values: []int{4, 5, 6}},
		}
		stream := From(data)
		result := stream.ContainsBy([]int{1, 2, 3}, func(d Data) any { return d.Values })

		if !result {
			t.Errorf("expected true, got false")
		}
	})
}
