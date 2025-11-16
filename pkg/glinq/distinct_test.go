package glinq

import (
	"reflect"
	"testing"
)

func TestDistinct(t *testing.T) {
	t.Run("Remove duplicates from integers", func(t *testing.T) {
		input := []int{1, 2, 2, 3, 3, 4, 1}
		result := Distinct(From(input)).ToSlice()

		expected := []int{1, 2, 3, 4}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("Distinct with structs", func(t *testing.T) {
		type Point struct{ X, Y int }

		input := []Point{{1, 2}, {1, 2}, {3, 4}}
		result := Distinct(From(input)).ToSlice()

		if len(result) != 2 {
			t.Errorf("Expected 2 unique points, got %d", len(result))
		}
	})

	t.Run("Distinct in chain", func(t *testing.T) {
		input := []int{1, 2, 2, 3, 3, 4, 5, 5}
		result := Distinct(From(input)).
			Where(func(x int) bool { return x > 2 }).
			ToSlice()

		expected := []int{3, 4, 5}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})
}

func TestDistinctBy(t *testing.T) {
	t.Run("Distinct by ID field", func(t *testing.T) {
		type Person struct {
			ID   int
			Name string
		}

		people := []Person{
			{1, "Alice"},
			{1, "Alice Duplicate"},
			{2, "Bob"},
			{2, "Bob Duplicate"},
		}

		result := From(people).
			DistinctBy(func(p Person) any { return p.ID }).
			ToSlice()

		if len(result) != 2 {
			t.Errorf("Expected 2 unique people, got %d", len(result))
		}
	})

	t.Run("Distinct by composite key", func(t *testing.T) {
		type Item struct{ Category, Name string }

		items := []Item{
			{"Food", "Apple"},
			{"Food", "Apple"},
			{"Drink", "Water"},
		}

		result := From(items).
			DistinctBy(func(i Item) any {
				return i.Category + "-" + i.Name
			}).
			ToSlice()

		if len(result) != 2 {
			t.Errorf("Expected 2 unique items, got %d", len(result))
		}
	})
}
