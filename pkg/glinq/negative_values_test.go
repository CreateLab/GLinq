package glinq

import (
	"testing"
)

// TestTake_NegativeValue tests that Take with negative values returns empty stream
func TestTake_NegativeValue(t *testing.T) {
	t.Run("Take with negative value", func(t *testing.T) {
		s := From([]int{1, 2, 3, 4, 5})
		result := s.Take(-1).ToSlice()

		if len(result) != 0 {
			t.Errorf("Expected empty slice for Take(-1), got %v", result)
		}
	})

	t.Run("Take with negative value on empty stream", func(t *testing.T) {
		s := Empty[int]()
		result := s.Take(-5).ToSlice()

		if len(result) != 0 {
			t.Errorf("Expected empty slice for Take(-5) on empty stream, got %v", result)
		}
	})

	t.Run("Take with zero value", func(t *testing.T) {
		s := From([]int{1, 2, 3})
		result := s.Take(0).ToSlice()

		if len(result) != 0 {
			t.Errorf("Expected empty slice for Take(0), got %v", result)
		}
	})

	t.Run("Take with large negative value", func(t *testing.T) {
		s := From([]int{1, 2, 3})
		result := s.Take(-1000).ToSlice()

		if len(result) != 0 {
			t.Errorf("Expected empty slice for Take(-1000), got %v", result)
		}
	})

	t.Run("Take negative preserves size information", func(t *testing.T) {
		s := From([]int{1, 2, 3})
		result := s.Take(-1)

		size, ok := result.Size()
		if !ok || size != 0 {
			t.Errorf("Expected size 0 for Take(-1), got size=%d, ok=%v", size, ok)
		}
	})
}

// TestSkip_NegativeValue tests that Skip with negative values treats them as 0
func TestSkip_NegativeValue(t *testing.T) {
	t.Run("Skip with negative value", func(t *testing.T) {
		s := From([]int{1, 2, 3, 4, 5})
		result := s.Skip(-1).ToSlice()

		expected := []int{1, 2, 3, 4, 5}
		if len(result) != len(expected) {
			t.Errorf("Expected %v for Skip(-1), got %v", expected, result)
		}
		for i, v := range result {
			if v != expected[i] {
				t.Errorf("Expected %d at index %d, got %d", expected[i], i, v)
			}
		}
	})

	t.Run("Skip with large negative value", func(t *testing.T) {
		s := From([]int{1, 2, 3})
		result := s.Skip(-1000).ToSlice()

		expected := []int{1, 2, 3}
		if len(result) != len(expected) {
			t.Errorf("Expected %v for Skip(-1000), got %v", expected, result)
		}
	})

	t.Run("Skip with zero value", func(t *testing.T) {
		s := From([]int{1, 2, 3})
		result := s.Skip(0).ToSlice()

		expected := []int{1, 2, 3}
		if len(result) != len(expected) {
			t.Errorf("Expected %v for Skip(0), got %v", expected, result)
		}
	})

	t.Run("Skip negative preserves size information", func(t *testing.T) {
		s := From([]int{1, 2, 3, 4, 5})
		result := s.Skip(-1)

		size, ok := result.Size()
		if !ok || size != 5 {
			t.Errorf("Expected size 5 for Skip(-1), got size=%d, ok=%v", size, ok)
		}
	})

	t.Run("Skip more than size returns empty", func(t *testing.T) {
		s := From([]int{1, 2, 3})
		result := s.Skip(10).ToSlice()

		if len(result) != 0 {
			t.Errorf("Expected empty slice for Skip(10) on stream of size 3, got %v", result)
		}
	})

	t.Run("Skip exactly size returns empty", func(t *testing.T) {
		s := From([]int{1, 2, 3})
		result := s.Skip(3).ToSlice()

		if len(result) != 0 {
			t.Errorf("Expected empty slice for Skip(3) on stream of size 3, got %v", result)
		}
	})
}

// TestRange_NegativeValue tests that Range with negative count returns empty stream
func TestRange_NegativeValue(t *testing.T) {
	t.Run("Range with negative count", func(t *testing.T) {
		result := Range(1, -1).ToSlice()

		if len(result) != 0 {
			t.Errorf("Expected empty slice for Range(1, -1), got %v", result)
		}
	})

	t.Run("Range with zero count", func(t *testing.T) {
		result := Range(1, 0).ToSlice()

		if len(result) != 0 {
			t.Errorf("Expected empty slice for Range(1, 0), got %v", result)
		}
	})

	t.Run("Range with large negative count", func(t *testing.T) {
		result := Range(1, -1000).ToSlice()

		if len(result) != 0 {
			t.Errorf("Expected empty slice for Range(1, -1000), got %v", result)
		}
	})

	t.Run("Range negative preserves size information", func(t *testing.T) {
		s := Range(1, -5)

		size, ok := s.Size()
		if !ok || size != 0 {
			t.Errorf("Expected size 0 for Range(1, -5), got size=%d, ok=%v", size, ok)
		}
	})

	t.Run("Range with positive count works normally", func(t *testing.T) {
		result := Range(1, 5).ToSlice()
		expected := []int{1, 2, 3, 4, 5}

		if len(result) != len(expected) {
			t.Errorf("Expected %v for Range(1, 5), got %v", expected, result)
		}
		for i, v := range result {
			if v != expected[i] {
				t.Errorf("Expected %d at index %d, got %d", expected[i], i, v)
			}
		}
	})
}

// TestChunk_NegativeValue tests that Chunk with negative or zero size returns nil
func TestChunk_NegativeValue(t *testing.T) {
	t.Run("Chunk with negative size", func(t *testing.T) {
		s := From([]int{1, 2, 3, 4, 5})
		result := s.Chunk(-1)

		if result != nil {
			t.Errorf("Expected nil for Chunk(-1), got %v", result)
		}
	})

	t.Run("Chunk with zero size", func(t *testing.T) {
		s := From([]int{1, 2, 3, 4, 5})
		result := s.Chunk(0)

		if result != nil {
			t.Errorf("Expected nil for Chunk(0), got %v", result)
		}
	})

	t.Run("Chunk with large negative size", func(t *testing.T) {
		s := From([]int{1, 2, 3, 4, 5})
		result := s.Chunk(-1000)

		if result != nil {
			t.Errorf("Expected nil for Chunk(-1000), got %v", result)
		}
	})

	t.Run("Chunk with positive size works normally", func(t *testing.T) {
		s := From([]int{1, 2, 3, 4, 5, 6, 7})
		result := s.Chunk(3)

		expected := [][]int{{1, 2, 3}, {4, 5, 6}, {7}}
		if len(result) != len(expected) {
			t.Errorf("Expected %v for Chunk(3), got %v", expected, result)
		}
	})
}

// TestTakeOrderedBy_NegativeValue tests that TakeOrderedBy with negative values returns empty stream
func TestTakeOrderedBy_NegativeValue(t *testing.T) {
	t.Run("TakeOrderedBy with negative value", func(t *testing.T) {
		s := From([]int{5, 2, 8, 1, 9, 3})
		result := TakeOrderedBy(s, -1, func(a, b int) bool { return a < b }).ToSlice()

		if len(result) != 0 {
			t.Errorf("Expected empty slice for TakeOrderedBy(-1), got %v", result)
		}
	})

	t.Run("TakeOrderedBy with zero value", func(t *testing.T) {
		s := From([]int{5, 2, 8, 1, 9, 3})
		result := TakeOrderedBy(s, 0, func(a, b int) bool { return a < b }).ToSlice()

		if len(result) != 0 {
			t.Errorf("Expected empty slice for TakeOrderedBy(0), got %v", result)
		}
	})

	t.Run("TakeOrderedBy with positive value works normally", func(t *testing.T) {
		s := From([]int{5, 2, 8, 1, 9, 3})
		result := TakeOrderedBy(s, 3, func(a, b int) bool { return a < b }).ToSlice()

		// TakeOrderedBy returns the 3 smallest elements in sorted order
		if len(result) != 3 {
			t.Errorf("Expected length 3 for TakeOrderedBy(3), got %d", len(result))
		}
		// Check that we got 3 elements and they are sorted
		if len(result) >= 1 && result[0] != 1 {
			t.Errorf("Expected smallest element to be 1, got %d", result[0])
		}
		// Verify all elements are from the source
		validElements := map[int]bool{1: true, 2: true, 3: true, 5: true, 8: true, 9: true}
		for _, v := range result {
			if !validElements[v] {
				t.Errorf("Got unexpected element %d", v)
			}
		}
		// Verify they are in ascending order
		for i := 1; i < len(result); i++ {
			if result[i-1] > result[i] {
				t.Errorf("Result not sorted: %v", result)
			}
		}
	})
}

// TestNegativeValues_Chain tests chaining operations with negative values
func TestNegativeValues_Chain(t *testing.T) {
	t.Run("Chain with Take negative", func(t *testing.T) {
		s := From([]int{1, 2, 3, 4, 5})
		result := s.Where(func(x int) bool { return x > 2 }).
			Take(-1).
			Select(func(x int) int { return x * 2 }).
			ToSlice()

		if len(result) != 0 {
			t.Errorf("Expected empty slice for chain with Take(-1), got %v", result)
		}
	})

	t.Run("Chain with Skip negative", func(t *testing.T) {
		s := From([]int{1, 2, 3, 4, 5})
		result := s.Skip(-1).
			Where(func(x int) bool { return x > 2 }).
			Take(2).
			ToSlice()

		expected := []int{3, 4}
		if len(result) != len(expected) {
			t.Errorf("Expected %v for chain with Skip(-1), got %v", expected, result)
		}
	})

	t.Run("Multiple negative operations", func(t *testing.T) {
		s := From([]int{1, 2, 3})
		result := s.Take(-5).Skip(-10).ToSlice()

		if len(result) != 0 {
			t.Errorf("Expected empty slice for Take(-5).Skip(-10), got %v", result)
		}
	})
}
