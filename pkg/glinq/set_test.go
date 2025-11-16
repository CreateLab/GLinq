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
