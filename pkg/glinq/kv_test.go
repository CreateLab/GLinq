package glinq

import (
	"testing"
)

func TestFromMap(t *testing.T) {
	m := map[string]int{
		"apple":  5,
		"banana": 3,
		"orange": 8,
	}
	stream := FromMap(m)
	result := stream.ToSlice()

	if len(result) != 3 {
		t.Errorf("expected length 3, got %d", len(result))
	}
}

func TestFromMapEmpty(t *testing.T) {
	m := make(map[string]int)
	stream := FromMap(m)
	result := stream.ToSlice()

	if len(result) != 0 {
		t.Errorf("expected empty result, got %v", result)
	}
}

func TestKeys(t *testing.T) {
	m := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
	}
	stream := FromMap(m)
	keysStream := Keys(stream)
	result := keysStream.ToSlice()

	if len(result) != 3 {
		t.Errorf("expected length 3, got %d", len(result))
	}

	// Check if all keys are present
	keyMap := make(map[string]bool)
	for _, k := range result {
		keyMap[k] = true
	}

	for k := range m {
		if !keyMap[k] {
			t.Errorf("expected key %s, not found in result", k)
		}
	}
}

func TestValues(t *testing.T) {
	m := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
	}
	stream := FromMap(m)
	valuesStream := Values(stream)
	result := valuesStream.ToSlice()

	if len(result) != 3 {
		t.Errorf("expected length 3, got %d", len(result))
	}

	// Check if all values are present
	valueMap := make(map[int]bool)
	for _, v := range result {
		valueMap[v] = true
	}

	for _, v := range m {
		if !valueMap[v] {
			t.Errorf("expected value %d, not found in result", v)
		}
	}
}

func TestToMap(t *testing.T) {
	m := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
	}
	stream := FromMap(m)
	result := ToMap(stream)

	if len(result) != len(m) {
		t.Errorf("expected length %d, got %d", len(m), len(result))
	}

	for k, v := range m {
		if result[k] != v {
			t.Errorf("expected %s: %d, got %s: %d", k, v, k, result[k])
		}
	}
}

func TestMapFilter(t *testing.T) {
	m := map[string]int{
		"apple":  5,
		"banana": 3,
		"orange": 8,
	}
	stream := FromMap(m).Where(func(kv KeyValue[string, int]) bool {
		return kv.Value > 4
	})
	result := ToMap(stream)

	if len(result) != 2 {
		t.Errorf("expected length 2, got %d", len(result))
	}

	if result["apple"] != 5 {
		t.Errorf("expected apple: 5, got %d", result["apple"])
	}

	if result["orange"] != 8 {
		t.Errorf("expected orange: 8, got %d", result["orange"])
	}

	if _, ok := result["banana"]; ok {
		t.Errorf("expected banana to be filtered out")
	}
}
