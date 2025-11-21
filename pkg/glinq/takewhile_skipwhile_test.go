package glinq

import (
	"reflect"
	"testing"
)

// TestTakeWhile tests the TakeWhile operation
func TestTakeWhile(t *testing.T) {
	t.Run("TakeWhile with predicate", func(t *testing.T) {
		numbers := []int{1, 2, 3, 4, 5, 6}
		result := From(numbers).
			TakeWhile(func(x int) bool { return x < 4 }).
			ToSlice()

		expected := []int{1, 2, 3}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("TakeWhile stops at first false", func(t *testing.T) {
		numbers := []int{2, 4, 6, 3, 8, 10}
		result := From(numbers).
			TakeWhile(func(x int) bool { return x%2 == 0 }).
			ToSlice()

		expected := []int{2, 4, 6}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("TakeWhile with all elements matching", func(t *testing.T) {
		numbers := []int{1, 2, 3}
		result := From(numbers).
			TakeWhile(func(x int) bool { return x < 10 }).
			ToSlice()

		expected := []int{1, 2, 3}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("TakeWhile with no elements matching", func(t *testing.T) {
		numbers := []int{5, 6, 7}
		result := From(numbers).
			TakeWhile(func(x int) bool { return x < 3 }).
			ToSlice()

		if len(result) != 0 {
			t.Errorf("Expected empty slice, got %v", result)
		}
	})

	t.Run("TakeWhile with empty stream", func(t *testing.T) {
		result := Empty[int]().
			TakeWhile(func(x int) bool { return x < 10 }).
			ToSlice()

		if len(result) != 0 {
			t.Errorf("Expected empty slice, got %v", result)
		}
	})

	t.Run("TakeWhile preserves lazy evaluation", func(t *testing.T) {
		called := 0
		numbers := []int{1, 2, 3, 4, 5}

		stream := From(numbers).
			TakeWhile(func(x int) bool {
				called++
				return x < 3
			})

		// Should not be called until materialization
		if called != 0 {
			t.Errorf("Expected predicate not called, but was called %d times", called)
		}

		result := stream.ToSlice()

		// Should be called only for elements 1, 2, and 3 (stops at 3)
		if called != 3 {
			t.Errorf("Expected predicate called 3 times, but was called %d times", called)
		}

		expected := []int{1, 2}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("TakeWhile loses size information", func(t *testing.T) {
		s := From([]int{1, 2, 3, 4, 5})
		result := s.TakeWhile(func(x int) bool { return x < 4 })

		size, ok := result.Size()
		if ok {
			t.Errorf("Expected unknown size for TakeWhile, got size=%d, ok=%v", size, ok)
		}
	})
}

// TestSkipWhile tests the SkipWhile operation
func TestSkipWhile(t *testing.T) {
	t.Run("SkipWhile with predicate", func(t *testing.T) {
		numbers := []int{1, 2, 3, 4, 5, 6}
		result := From(numbers).
			SkipWhile(func(x int) bool { return x < 4 }).
			ToSlice()

		expected := []int{4, 5, 6}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("SkipWhile skips until first false", func(t *testing.T) {
		numbers := []int{2, 4, 6, 3, 8, 10}
		result := From(numbers).
			SkipWhile(func(x int) bool { return x%2 == 0 }).
			ToSlice()

		expected := []int{3, 8, 10}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("SkipWhile with all elements matching", func(t *testing.T) {
		numbers := []int{1, 2, 3}
		result := From(numbers).
			SkipWhile(func(x int) bool { return x < 10 }).
			ToSlice()

		if len(result) != 0 {
			t.Errorf("Expected empty slice, got %v", result)
		}
	})

	t.Run("SkipWhile with no elements matching", func(t *testing.T) {
		numbers := []int{5, 6, 7}
		result := From(numbers).
			SkipWhile(func(x int) bool { return x < 3 }).
			ToSlice()

		expected := []int{5, 6, 7}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("SkipWhile with empty stream", func(t *testing.T) {
		result := Empty[int]().
			SkipWhile(func(x int) bool { return x < 10 }).
			ToSlice()

		if len(result) != 0 {
			t.Errorf("Expected empty slice, got %v", result)
		}
	})

	t.Run("SkipWhile preserves lazy evaluation", func(t *testing.T) {
		called := 0
		numbers := []int{1, 2, 3, 4, 5}

		stream := From(numbers).
			SkipWhile(func(x int) bool {
				called++
				return x < 3
			})

		// Should not be called until materialization
		if called != 0 {
			t.Errorf("Expected predicate not called, but was called %d times", called)
		}

		result := stream.ToSlice()

		// Should be called for elements 1, 2, and 3 (stops skipping at 3)
		if called != 3 {
			t.Errorf("Expected predicate called 3 times, but was called %d times", called)
		}

		expected := []int{3, 4, 5}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("SkipWhile loses size information", func(t *testing.T) {
		s := From([]int{1, 2, 3, 4, 5})
		result := s.SkipWhile(func(x int) bool { return x < 4 })

		size, ok := result.Size()
		if ok {
			t.Errorf("Expected unknown size for SkipWhile, got size=%d, ok=%v", size, ok)
		}
	})
}

// TestTakeWhile_Chain tests chaining TakeWhile with other operations
func TestTakeWhile_Chain(t *testing.T) {
	t.Run("TakeWhile with Where", func(t *testing.T) {
		numbers := []int{1, 2, 3, 4, 5, 6, 7, 8}
		result := From(numbers).
			TakeWhile(func(x int) bool { return x < 6 }).
			Where(func(x int) bool { return x%2 == 0 }).
			ToSlice()

		expected := []int{2, 4}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("TakeWhile with Select", func(t *testing.T) {
		numbers := []int{1, 2, 3, 4, 5}
		result := From(numbers).
			TakeWhile(func(x int) bool { return x < 4 }).
			Select(func(x int) int { return x * 2 }).
			ToSlice()

		expected := []int{2, 4, 6}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("TakeWhile with Take", func(t *testing.T) {
		numbers := []int{1, 2, 3, 4, 5, 6, 7, 8}
		result := From(numbers).
			TakeWhile(func(x int) bool { return x < 10 }).
			Take(3).
			ToSlice()

		expected := []int{1, 2, 3}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})
}

// TestSkipWhile_Chain tests chaining SkipWhile with other operations
func TestSkipWhile_Chain(t *testing.T) {
	t.Run("SkipWhile with Where", func(t *testing.T) {
		numbers := []int{1, 2, 3, 4, 5, 6, 7, 8}
		result := From(numbers).
			SkipWhile(func(x int) bool { return x < 4 }).
			Where(func(x int) bool { return x%2 == 0 }).
			ToSlice()

		expected := []int{4, 6, 8}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("SkipWhile with Select", func(t *testing.T) {
		numbers := []int{1, 2, 3, 4, 5}
		result := From(numbers).
			SkipWhile(func(x int) bool { return x < 3 }).
			Select(func(x int) int { return x * 2 }).
			ToSlice()

		expected := []int{6, 8, 10}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("SkipWhile with Take", func(t *testing.T) {
		numbers := []int{1, 2, 3, 4, 5, 6, 7, 8}
		result := From(numbers).
			SkipWhile(func(x int) bool { return x < 4 }).
			Take(2).
			ToSlice()

		expected := []int{4, 5}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})
}

// TestTakeWhile_SkipWhile_Combined tests combining TakeWhile and SkipWhile
func TestTakeWhile_SkipWhile_Combined(t *testing.T) {
	t.Run("SkipWhile then TakeWhile", func(t *testing.T) {
		numbers := []int{1, 2, 3, 4, 5, 6, 7, 8}
		result := From(numbers).
			SkipWhile(func(x int) bool { return x < 3 }).
			TakeWhile(func(x int) bool { return x < 6 }).
			ToSlice()

		expected := []int{3, 4, 5}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("TakeWhile then SkipWhile", func(t *testing.T) {
		numbers := []int{1, 2, 3, 4, 5, 6, 7, 8}
		result := From(numbers).
			TakeWhile(func(x int) bool { return x < 6 }).
			SkipWhile(func(x int) bool { return x < 3 }).
			ToSlice()

		expected := []int{3, 4, 5}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})
}
