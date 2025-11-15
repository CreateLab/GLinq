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
