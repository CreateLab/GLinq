package glinq

import (
	"testing"
)

func TestToSlice(t *testing.T) {
	slice := []int{1, 2, 3}
	stream := From(slice)
	result := stream.ToSlice()

	if len(result) != 3 {
		t.Errorf("expected length 3, got %d", len(result))
	}
}

func TestFirstWithElements(t *testing.T) {
	slice := []int{1, 2, 3}
	stream := From(slice)
	value, ok := stream.First()

	if !ok {
		t.Errorf("expected ok=true, got false")
	}
	if value != 1 {
		t.Errorf("expected 1, got %d", value)
	}
}

func TestFirstEmpty(t *testing.T) {
	stream := Empty[int]()
	value, ok := stream.First()

	if ok {
		t.Errorf("expected ok=false, got true")
	}
	if value != 0 {
		t.Errorf("expected 0, got %d", value)
	}
}

func TestCount(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	stream := From(slice)
	count := stream.Count()

	if count != 5 {
		t.Errorf("expected count 5, got %d", count)
	}
}

func TestCountEmpty(t *testing.T) {
	stream := Empty[int]()
	count := stream.Count()

	if count != 0 {
		t.Errorf("expected count 0, got %d", count)
	}
}

func TestCountWithFilter(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	stream := From(slice).Where(func(x int) bool { return x > 2 })
	count := stream.Count()

	if count != 3 {
		t.Errorf("expected count 3, got %d", count)
	}
}

func TestAnyTrue(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	stream := From(slice)
	result := stream.Any(func(x int) bool { return x == 3 })

	if !result {
		t.Errorf("expected true, got false")
	}
}

func TestAnyFalse(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	stream := From(slice)
	result := stream.Any(func(x int) bool { return x == 10 })

	if result {
		t.Errorf("expected false, got true")
	}
}

func TestAnyEmpty(t *testing.T) {
	stream := Empty[int]()
	result := stream.Any(func(x int) bool { return x == 1 })

	if result {
		t.Errorf("expected false, got true")
	}
}

func TestAllTrue(t *testing.T) {
	slice := []int{2, 4, 6, 8}
	stream := From(slice)
	result := stream.All(func(x int) bool { return x%2 == 0 })

	if !result {
		t.Errorf("expected true, got false")
	}
}

func TestAllFalse(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	stream := From(slice)
	result := stream.All(func(x int) bool { return x > 2 })

	if result {
		t.Errorf("expected false, got true")
	}
}

func TestAllEmpty(t *testing.T) {
	stream := Empty[int]()
	result := stream.All(func(x int) bool { return x == 1 })

	if !result {
		t.Errorf("expected true for empty stream, got false")
	}
}

func TestForEach(t *testing.T) {
	slice := []int{1, 2, 3}
	stream := From(slice)
	var result []int

	stream.ForEach(func(x int) {
		result = append(result, x*2)
	})

	expected := []int{2, 4, 6}
	if len(result) != len(expected) {
		t.Errorf("expected length %d, got %d", len(expected), len(result))
	}

	for i, v := range result {
		if v != expected[i] {
			t.Errorf("expected %d, got %d at index %d", expected[i], v, i)
		}
	}
}
