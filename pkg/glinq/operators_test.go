package glinq

import (
	"fmt"
	"reflect"
	"testing"
)

func TestWhere(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5, 6}
	stream := From(slice).Where(func(x int) bool { return x > 3 })
	result := stream.ToSlice()

	expected := []int{4, 5, 6}
	if len(result) != len(expected) {
		t.Errorf("expected length %d, got %d", len(expected), len(result))
	}

	for i, v := range result {
		if v != expected[i] {
			t.Errorf("expected %d, got %d at index %d", expected[i], v, i)
		}
	}
}

func TestTake(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	stream := From(slice).Take(3)
	result := stream.ToSlice()

	expected := []int{1, 2, 3}
	if len(result) != len(expected) {
		t.Errorf("expected length %d, got %d", len(expected), len(result))
	}

	for i, v := range result {
		if v != expected[i] {
			t.Errorf("expected %d, got %d at index %d", expected[i], v, i)
		}
	}
}

func TestTakeMoreThanAvailable(t *testing.T) {
	slice := []int{1, 2, 3}
	stream := From(slice).Take(10)
	result := stream.ToSlice()

	if len(result) != 3 {
		t.Errorf("expected length 3, got %d", len(result))
	}
}

func TestSkip(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	stream := From(slice).Skip(2)
	result := stream.ToSlice()

	expected := []int{3, 4, 5}
	if len(result) != len(expected) {
		t.Errorf("expected length %d, got %d", len(expected), len(result))
	}

	for i, v := range result {
		if v != expected[i] {
			t.Errorf("expected %d, got %d at index %d", expected[i], v, i)
		}
	}
}
func TestSelect(t *testing.T) {
	t.Run("Transform same type", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5}
		result := From(input).
			Select(func(x int) int { return x * 2 }).
			ToSlice()

		expected := []int{2, 4, 6, 8, 10}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("Chaining with Where", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5}
		result := From(input).
			Where(func(x int) bool { return x > 2 }).
			Select(func(x int) int { return x * 10 }).
			ToSlice()

		expected := []int{30, 40, 50}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})
}

func TestMap(t *testing.T) {
	t.Run("Transform int to string", func(t *testing.T) {
		input := []int{1, 2, 3}
		result := Select(
			From(input),
			func(x int) string { return fmt.Sprintf("num_%d", x) },
		).ToSlice()

		expected := []string{"num_1", "num_2", "num_3"}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("Transform to struct", func(t *testing.T) {
		type User struct {
			ID   int
			Name string
		}

		input := []int{1, 2, 3}
		result := Select(
			From(input),
			func(id int) User {
				return User{ID: id, Name: fmt.Sprintf("User%d", id)}
			},
		).ToSlice()

		if len(result) != 3 {
			t.Errorf("Expected 3 users, got %d", len(result))
		}
		if result[0].ID != 1 || result[0].Name != "User1" {
			t.Errorf("Unexpected user: %+v", result[0])
		}
	})

	t.Run("Chaining with Where and Select", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5}
		result := Select(
			From(input).Where(func(x int) bool { return x > 2 }),
			func(x int) string { return fmt.Sprintf("Number: %d", x) },
		).ToSlice()

		expected := []string{"Number: 3", "Number: 4", "Number: 5"}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})
}

func TestLazyEvaluation(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	called := 0

	stream := From(slice).Where(func(x int) bool {
		called++
		return x > 2
	})

	// Predicate should not be called until ToSlice() is called
	if called != 0 {
		t.Errorf("expected predicate not called, but was called %d times", called)
	}

	result := stream.ToSlice()

	// After ToSlice, predicate should have been called for all elements
	if called != 5 {
		t.Errorf("expected predicate called 5 times, but was called %d times", called)
	}

	if len(result) != 3 {
		t.Errorf("expected length 3, got %d", len(result))
	}
}

func TestSelectWithIndex(t *testing.T) {
	t.Run("Transform same type with index", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5}
		result := From(input).
			SelectWithIndex(func(x int, idx int) int { return x * idx }).
			ToSlice()

		expected := []int{0, 2, 6, 12, 20}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("Chaining with Where and SelectWithIndex", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5}
		result := From(input).
			Where(func(x int) bool { return x > 2 }).
			SelectWithIndex(func(x int, idx int) int { return x + idx*10 }).
			ToSlice()

		expected := []int{3, 14, 25}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})
}

func TestSelectWithIndexFunction(t *testing.T) {
	t.Run("Transform int to string with index", func(t *testing.T) {
		input := []int{1, 2, 3}
		result := SelectWithIndex(
			From(input),
			func(x int, idx int) string { return fmt.Sprintf("num_%d_at_%d", x, idx) },
		).ToSlice()

		expected := []string{"num_1_at_0", "num_2_at_1", "num_3_at_2"}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("Transform to struct with index", func(t *testing.T) {
		type User struct {
			ID   int
			Name string
		}

		input := []int{1, 2, 3}
		result := SelectWithIndex(
			From(input),
			func(id int, idx int) User {
				return User{ID: id, Name: fmt.Sprintf("User%d_Index%d", id, idx)}
			},
		).ToSlice()

		if len(result) != 3 {
			t.Errorf("Expected 3 users, got %d", len(result))
		}
		if result[0].ID != 1 || result[0].Name != "User1_Index0" {
			t.Errorf("Unexpected user: %+v", result[0])
		}
		if result[1].ID != 2 || result[1].Name != "User2_Index1" {
			t.Errorf("Unexpected user: %+v", result[1])
		}
	})

	t.Run("Chaining with Where and SelectWithIndex function", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5}
		result := SelectWithIndex(
			From(input).Where(func(x int) bool { return x > 2 }),
			func(x int, idx int) string { return fmt.Sprintf("Number_%d_Index_%d", x, idx) },
		).ToSlice()

		expected := []string{"Number_3_Index_0", "Number_4_Index_1", "Number_5_Index_2"}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})
}

func TestReverse(t *testing.T) {
	t.Run("Reverse integers", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5}
		result := From(input).Reverse().ToSlice()

		expected := []int{5, 4, 3, 2, 1}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("Reverse with chaining", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5}
		result := From(input).
			Where(func(x int) bool { return x > 2 }).
			Reverse().
			ToSlice()

		expected := []int{5, 4, 3}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("Reverse empty stream", func(t *testing.T) {
		result := Empty[int]().Reverse().ToSlice()

		if len(result) != 0 {
			t.Errorf("Expected empty slice, got %v", result)
		}
	})

	t.Run("Reverse single element", func(t *testing.T) {
		input := []int{42}
		result := From(input).Reverse().ToSlice()

		expected := []int{42}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})
}

func TestSelectMany(t *testing.T) {
	t.Run("Flatten slices", func(t *testing.T) {
		input := [][]int{{1, 2}, {3, 4}, {5}}
		result := SelectMany(
			From(input),
			func(slice []int) Enumerable[int] { return From(slice) },
		).ToSlice()

		expected := []int{1, 2, 3, 4, 5}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("Flatten with empty slices", func(t *testing.T) {
		input := [][]int{{1, 2}, {}, {3, 4}, {}}
		result := SelectMany(
			From(input),
			func(slice []int) Enumerable[int] { return From(slice) },
		).ToSlice()

		expected := []int{1, 2, 3, 4}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("Flatten strings", func(t *testing.T) {
		input := []string{"hello", "world"}
		result := SelectMany(
			From(input),
			func(s string) Enumerable[rune] {
				runes := make([]rune, len(s))
				for i, r := range s {
					runes[i] = r
				}
				return From(runes)
			},
		).ToSlice()

		expected := []rune{'h', 'e', 'l', 'l', 'o', 'w', 'o', 'r', 'l', 'd'}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("SelectMany with chaining", func(t *testing.T) {
		input := [][]int{{1, 2, 3}, {4, 5}, {6, 7, 8}}
		result := SelectMany(
			From(input),
			func(slice []int) Enumerable[int] { return From(slice) },
		).Where(func(x int) bool { return x%2 == 0 }).ToSlice()

		expected := []int{2, 4, 6, 8}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})
}
