package glinq

import (
	"reflect"
	"testing"
)

func TestChunk(t *testing.T) {
	t.Run("Chunk with exact division", func(t *testing.T) {
		slice := []int{1, 2, 3, 4, 5, 6}
		result := From(slice).Chunk(3)

		expected := [][]int{{1, 2, 3}, {4, 5, 6}}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("Chunk with remainder", func(t *testing.T) {
		slice := []int{1, 2, 3, 4, 5, 6, 7}
		result := From(slice).Chunk(3)

		expected := [][]int{{1, 2, 3}, {4, 5, 6}, {7}}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("Chunk with single element chunks", func(t *testing.T) {
		slice := []int{1, 2, 3}
		result := From(slice).Chunk(1)

		expected := [][]int{{1}, {2}, {3}}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("Chunk larger than slice", func(t *testing.T) {
		slice := []int{1, 2, 3}
		result := From(slice).Chunk(10)

		expected := [][]int{{1, 2, 3}}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("Chunk empty slice", func(t *testing.T) {
		slice := []int{}
		result := From(slice).Chunk(3)

		if len(result) != 0 {
			t.Errorf("Expected empty result, got %v", result)
		}
	})

	t.Run("Chunk with zero size", func(t *testing.T) {
		slice := []int{1, 2, 3}
		result := From(slice).Chunk(0)

		if result != nil {
			t.Errorf("Expected nil, got %v", result)
		}
	})

	t.Run("Chunk with negative size", func(t *testing.T) {
		slice := []int{1, 2, 3}
		result := From(slice).Chunk(-1)

		if result != nil {
			t.Errorf("Expected nil, got %v", result)
		}
	})

	t.Run("Chunk with filtered stream", func(t *testing.T) {
		slice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		result := From(slice).
			Where(func(x int) bool { return x%2 == 0 }).
			Chunk(2)

		expected := [][]int{{2, 4}, {6, 8}, {10}}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("Chunk with strings", func(t *testing.T) {
		slice := []string{"a", "b", "c", "d", "e"}
		result := From(slice).Chunk(2)

		expected := [][]string{{"a", "b"}, {"c", "d"}, {"e"}}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})
}
