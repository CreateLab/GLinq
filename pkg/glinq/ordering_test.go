package glinq

import (
	"reflect"
	"testing"
)

func TestOrderBy(t *testing.T) {
	t.Run("Sort integers ascending", func(t *testing.T) {
		input := []int{5, 2, 8, 1, 9, 3}
		result := From(input).
			OrderBy(func(a, b int) int { return a - b }).
			ToSlice()

		expected := []int{1, 2, 3, 5, 8, 9}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("Sort structs by field", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}

		people := []Person{
			{Name: "Alice", Age: 30},
			{Name: "Bob", Age: 25},
			{Name: "Charlie", Age: 35},
		}

		result := From(people).
			OrderBy(func(a, b Person) int { return a.Age - b.Age }).
			ToSlice()

		if result[0].Age != 25 || result[1].Age != 30 || result[2].Age != 35 {
			t.Errorf("Unexpected order: %+v", result)
		}
	})

	t.Run("OrderBy with chain", func(t *testing.T) {
		input := []int{10, 3, 7, 1, 9, 2, 8, 4, 6, 5}
		result := From(input).
			Where(func(x int) bool { return x > 3 }).
			OrderBy(func(a, b int) int { return b - a }). // descending
			Take(3).
			ToSlice()

		expected := []int{10, 9, 8}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})
}

func TestOrderByDescending(t *testing.T) {
	input := []int{5, 2, 8, 1, 9, 3}
	result := From(input).
		OrderByDescending(func(a, b int) int { return a - b }).
		ToSlice()

	expected := []int{9, 8, 5, 3, 2, 1}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}
