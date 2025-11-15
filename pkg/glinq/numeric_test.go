package glinq

import (
	"testing"
)

func TestLast(t *testing.T) {
	t.Run("Last with elements", func(t *testing.T) {
		slice := []int{1, 2, 3, 4, 5}
		stream := From(slice)
		value, ok := stream.Last()

		if !ok {
			t.Errorf("expected ok=true, got false")
		}
		if value != 5 {
			t.Errorf("expected 5, got %d", value)
		}
	})

	t.Run("Last empty", func(t *testing.T) {
		stream := Empty[int]()
		value, ok := stream.Last()

		if ok {
			t.Errorf("expected ok=false, got true")
		}
		if value != 0 {
			t.Errorf("expected 0, got %d", value)
		}
	})

	t.Run("Last single element", func(t *testing.T) {
		slice := []int{42}
		stream := From(slice)
		value, ok := stream.Last()

		if !ok {
			t.Errorf("expected ok=true, got false")
		}
		if value != 42 {
			t.Errorf("expected 42, got %d", value)
		}
	})

	t.Run("Last with filtered stream", func(t *testing.T) {
		slice := []int{1, 2, 3, 4, 5}
		stream := From(slice).Where(func(x int) bool { return x > 3 })
		value, ok := stream.Last()

		if !ok {
			t.Errorf("expected ok=true, got false")
		}
		if value != 5 {
			t.Errorf("expected 5, got %d", value)
		}
	})
}

func TestSum(t *testing.T) {
	t.Run("Sum integers", func(t *testing.T) {
		slice := []int{1, 2, 3, 4, 5}
		result := Sum(From(slice))

		if result != 15 {
			t.Errorf("expected 15, got %d", result)
		}
	})

	t.Run("Sum empty", func(t *testing.T) {
		slice := []int{}
		result := Sum(From(slice))

		if result != 0 {
			t.Errorf("expected 0, got %d", result)
		}
	})

	t.Run("Sum floats", func(t *testing.T) {
		slice := []float64{1.5, 2.5, 3.5}
		result := Sum(From(slice))

		if result != 7.5 {
			t.Errorf("expected 7.5, got %f", result)
		}
	})

	t.Run("Sum with filtered stream", func(t *testing.T) {
		slice := []int{1, 2, 3, 4, 5}
		result := Sum(From(slice).Where(func(x int) bool { return x%2 == 0 }))

		if result != 6 {
			t.Errorf("expected 6, got %d", result)
		}
	})

	t.Run("Sum uint", func(t *testing.T) {
		slice := []uint{10, 20, 30}
		result := Sum(From(slice))

		if result != 60 {
			t.Errorf("expected 60, got %d", result)
		}
	})
}

func TestMin(t *testing.T) {
	t.Run("Min integers", func(t *testing.T) {
		slice := []int{5, 2, 8, 1, 9}
		value, ok := Min(From(slice))

		if !ok {
			t.Errorf("expected ok=true, got false")
		}
		if value != 1 {
			t.Errorf("expected 1, got %d", value)
		}
	})

	t.Run("Min empty", func(t *testing.T) {
		slice := []int{}
		value, ok := Min(From(slice))

		if ok {
			t.Errorf("expected ok=false, got true")
		}
		if value != 0 {
			t.Errorf("expected 0, got %d", value)
		}
	})

	t.Run("Min single element", func(t *testing.T) {
		slice := []int{42}
		value, ok := Min(From(slice))

		if !ok {
			t.Errorf("expected ok=true, got false")
		}
		if value != 42 {
			t.Errorf("expected 42, got %d", value)
		}
	})

	t.Run("Min floats", func(t *testing.T) {
		slice := []float64{5.5, 2.2, 8.8, 1.1}
		value, ok := Min(From(slice))

		if !ok {
			t.Errorf("expected ok=true, got false")
		}
		if value != 1.1 {
			t.Errorf("expected 1.1, got %f", value)
		}
	})

	t.Run("Min strings", func(t *testing.T) {
		slice := []string{"zebra", "apple", "banana"}
		value, ok := Min(From(slice))

		if !ok {
			t.Errorf("expected ok=true, got false")
		}
		if value != "apple" {
			t.Errorf("expected 'apple', got %s", value)
		}
	})

	t.Run("Min with filtered stream", func(t *testing.T) {
		slice := []int{5, 2, 8, 1, 9}
		value, ok := Min(From(slice).Where(func(x int) bool { return x > 3 }))

		if !ok {
			t.Errorf("expected ok=true, got false")
		}
		if value != 5 {
			t.Errorf("expected 5, got %d", value)
		}
	})
}

func TestMax(t *testing.T) {
	t.Run("Max integers", func(t *testing.T) {
		slice := []int{5, 2, 8, 1, 9}
		value, ok := Max(From(slice))

		if !ok {
			t.Errorf("expected ok=true, got false")
		}
		if value != 9 {
			t.Errorf("expected 9, got %d", value)
		}
	})

	t.Run("Max empty", func(t *testing.T) {
		slice := []int{}
		value, ok := Max(From(slice))

		if ok {
			t.Errorf("expected ok=false, got true")
		}
		if value != 0 {
			t.Errorf("expected 0, got %d", value)
		}
	})

	t.Run("Max single element", func(t *testing.T) {
		slice := []int{42}
		value, ok := Max(From(slice))

		if !ok {
			t.Errorf("expected ok=true, got false")
		}
		if value != 42 {
			t.Errorf("expected 42, got %d", value)
		}
	})

	t.Run("Max floats", func(t *testing.T) {
		slice := []float64{5.5, 2.2, 8.8, 1.1}
		value, ok := Max(From(slice))

		if !ok {
			t.Errorf("expected ok=true, got false")
		}
		if value != 8.8 {
			t.Errorf("expected 8.8, got %f", value)
		}
	})

	t.Run("Max strings", func(t *testing.T) {
		slice := []string{"zebra", "apple", "banana"}
		value, ok := Max(From(slice))

		if !ok {
			t.Errorf("expected ok=true, got false")
		}
		if value != "zebra" {
			t.Errorf("expected 'zebra', got %s", value)
		}
	})

	t.Run("Max with filtered stream", func(t *testing.T) {
		slice := []int{5, 2, 8, 1, 9}
		value, ok := Max(From(slice).Where(func(x int) bool { return x < 8 }))

		if !ok {
			t.Errorf("expected ok=true, got false")
		}
		if value != 5 {
			t.Errorf("expected 5, got %d", value)
		}
	})
}
