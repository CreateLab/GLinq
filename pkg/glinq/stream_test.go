package glinq

import (
	"testing"
)

func TestFrom(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	stream := From(slice)
	result := stream.ToSlice()

	if len(result) != len(slice) {
		t.Errorf("expected length %d, got %d", len(slice), len(result))
	}

	for i, v := range result {
		if v != slice[i] {
			t.Errorf("expected %d, got %d at index %d", slice[i], v, i)
		}
	}
}

func TestEmpty(t *testing.T) {
	stream := Empty[int]()
	result := stream.ToSlice()

	if len(result) != 0 {
		t.Errorf("expected empty stream, got %v", result)
	}
}

func TestRange(t *testing.T) {
	stream := Range(1, 5)
	result := stream.ToSlice()

	expected := []int{1, 2, 3, 4, 5}
	if len(result) != len(expected) {
		t.Errorf("expected length %d, got %d", len(expected), len(result))
	}

	for i, v := range result {
		if v != expected[i] {
			t.Errorf("expected %d, got %d at index %d", expected[i], v, i)
		}
	}
}

func TestRangeZeroCount(t *testing.T) {
	stream := Range(1, 0)
	result := stream.ToSlice()

	if len(result) != 0 {
		t.Errorf("expected empty stream, got %v", result)
	}
}

// TestFromReference проверяет что From держит ссылку на оригинальный слайс
func TestFromReference(t *testing.T) {
	data := []int{1, 2, 3}
	stream := From(data)

	// Изменяем исходный слайс
	data[0] = 999

	result := stream.ToSlice()
	if result[0] != 999 {
		t.Errorf("Expected From to reference original slice, got %d, expected 999", result[0])
	}
	if len(result) != 3 {
		t.Errorf("Expected length 3, got %d", len(result))
	}
}

// TestFromSafeCopy проверяет что FromSafe копирует слайс
func TestFromSafeCopy(t *testing.T) {
	data := []int{1, 2, 3}
	stream := FromSafe(data)

	// Изменяем исходный слайс
	data[0] = 999

	result := stream.ToSlice()
	if result[0] != 1 {
		t.Errorf("Expected FromSafe to copy slice, got %d, expected 1", result[0])
	}
	if len(result) != 3 {
		t.Errorf("Expected length 3, got %d", len(result))
	}
	// Проверяем что остальные элементы тоже не изменились
	if result[1] != 2 || result[2] != 3 {
		t.Errorf("Expected [1, 2, 3], got %v", result)
	}
}
