package glinq

import (
	"reflect"
	"sort"
	"testing"
)

func TestConcat(t *testing.T) {
	stream1 := From([]int{1, 2, 3})
	stream2 := From([]int{4, 5, 6})

	result := stream1.Concat(stream2).ToSlice()
	expected := []int{1, 2, 3, 4, 5, 6}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestUnion(t *testing.T) {
	stream1 := From([]int{1, 2, 3})
	stream2 := From([]int{3, 4, 5})

	result := Union(stream1, stream2).ToSlice()

	// Union preserves order but removes duplicates
	expected := []int{1, 2, 3, 4, 5}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestIntersect(t *testing.T) {
	stream1 := From([]int{1, 2, 3, 4})
	stream2 := From([]int{3, 4, 5, 6})

	result := Intersect(stream1, stream2).ToSlice()
	sort.Ints(result) // Sort for comparison

	expected := []int{3, 4}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestExcept(t *testing.T) {
	stream1 := From([]int{1, 2, 3, 4})
	stream2 := From([]int{3, 4, 5, 6})

	result := Except(stream1, stream2).ToSlice()
	sort.Ints(result) // Sort for comparison

	expected := []int{1, 2}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestZip(t *testing.T) {
	t.Run("Zip with equal length sequences", func(t *testing.T) {
		numbers := From([]int{1, 2, 3})
		letters := From([]string{"a", "b", "c"})

		result := Zip(numbers, letters, func(n int, s string) string {
			return string(rune('0'+n)) + ":" + s
		}).ToSlice()

		expected := []string{"1:a", "2:b", "3:c"}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("Zip with first sequence shorter", func(t *testing.T) {
		numbers := From([]int{1, 2})
		letters := From([]string{"a", "b", "c", "d"})

		result := Zip(numbers, letters, func(n int, s string) string {
			return string(rune('0'+n)) + ":" + s
		}).ToSlice()

		expected := []string{"1:a", "2:b"}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("Zip with second sequence shorter", func(t *testing.T) {
		numbers := From([]int{1, 2, 3, 4})
		letters := From([]string{"a", "b"})

		result := Zip(numbers, letters, func(n int, s string) string {
			return string(rune('0'+n)) + ":" + s
		}).ToSlice()

		expected := []string{"1:a", "2:b"}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("Zip with empty first sequence", func(t *testing.T) {
		numbers := Empty[int]()
		letters := From([]string{"a", "b", "c"})

		result := Zip(numbers, letters, func(n int, s string) string {
			return string(rune('0'+n)) + ":" + s
		}).ToSlice()

		if len(result) != 0 {
			t.Errorf("Expected empty slice, got %v", result)
		}
	})

	t.Run("Zip with empty second sequence", func(t *testing.T) {
		numbers := From([]int{1, 2, 3})
		letters := Empty[string]()

		result := Zip(numbers, letters, func(n int, s string) string {
			return string(rune('0'+n)) + ":" + s
		}).ToSlice()

		if len(result) != 0 {
			t.Errorf("Expected empty slice, got %v", result)
		}
	})

	t.Run("Zip with different types", func(t *testing.T) {
		numbers := From([]int{10, 20, 30})
		multipliers := From([]float64{1.5, 2.0, 3.5})

		result := Zip(numbers, multipliers, func(n int, m float64) float64 {
			return float64(n) * m
		}).ToSlice()

		expected := []float64{15.0, 40.0, 105.0}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("Zip with struct result", func(t *testing.T) {
		type Pair struct {
			Num int
			Str string
		}

		numbers := From([]int{1, 2, 3})
		letters := From([]string{"a", "b", "c"})

		result := Zip(numbers, letters, func(n int, s string) Pair {
			return Pair{Num: n, Str: s}
		}).ToSlice()

		expected := []Pair{{Num: 1, Str: "a"}, {Num: 2, Str: "b"}, {Num: 3, Str: "c"}}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})
}
