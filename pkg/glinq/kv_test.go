package glinq

import (
	"fmt"
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

// TestFromMapOnDemand проверяет что FromMap читает значения по требованию
func TestFromMapOnDemand(t *testing.T) {
	m := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
	}
	stream := FromMap(m)

	// Изменяем значение в мапе после создания стрима
	m["a"] = 999

	// Читаем из стрима - должно получить новое значение
	result := stream.ToSlice()

	// Проверяем что значение было прочитано по требованию (новое значение)
	found := false
	for _, kv := range result {
		if kv.Key == "a" {
			if kv.Value != 999 {
				t.Errorf("Expected FromMap to read values on-demand, got %d, expected 999", kv.Value)
			}
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected to find key 'a' in result")
	}
}

// TestFromMapSafeSnapshot проверяет что FromMapSafe делает полный снимок
func TestFromMapSafeSnapshot(t *testing.T) {
	m := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
	}
	stream := FromMapSafe(m)

	// Изменяем значение в мапе после создания стрима
	m["a"] = 999

	// Читаем из стрима - должно получить старое значение (снимок)
	result := stream.ToSlice()

	// Проверяем что значение было скопировано (старое значение)
	found := false
	for _, kv := range result {
		if kv.Key == "a" {
			if kv.Value != 1 {
				t.Errorf("Expected FromMapSafe to take snapshot, got %d, expected 1", kv.Value)
			}
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected to find key 'a' in result")
	}

	// Проверяем что остальные значения тоже не изменились
	for _, kv := range result {
		if kv.Key == "b" && kv.Value != 2 {
			t.Errorf("Expected value 2 for key 'b', got %d", kv.Value)
		}
		if kv.Key == "c" && kv.Value != 3 {
			t.Errorf("Expected value 3 for key 'c', got %d", kv.Value)
		}
	}
}

// TestFromMapKeysOnly проверяет что FromMap копирует только ключи
func TestFromMapKeysOnly(t *testing.T) {
	// Используем большую структуру для значений чтобы проверить что она не копируется
	type LargeValue struct {
		Data [1000]int
	}

	m := make(map[int]LargeValue)
	for i := 0; i < 10; i++ {
		m[i] = LargeValue{Data: [1000]int{i}}
	}

	stream := FromMap(m)
	result := stream.ToSlice()

	if len(result) != 10 {
		t.Errorf("Expected 10 elements, got %d", len(result))
	}

	// Проверяем что значения корректны
	for _, kv := range result {
		if kv.Value.Data[0] != kv.Key {
			t.Errorf("Expected value Data[0] == key, got %d != %d", kv.Value.Data[0], kv.Key)
		}
	}
}

// BenchmarkFromO1 проверяет что From() имеет O(1) сложность создания
func BenchmarkFromO1(b *testing.B) {
	sizes := []int{10, 100, 1000, 10000, 100000}
	for _, size := range sizes {
		data := make([]int, size)
		for i := range data {
			data[i] = i
		}

		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = From(data)
			}
		})
	}
}

// BenchmarkFromSafeON проверяет что FromSafe() имеет O(n) сложность создания
func BenchmarkFromSafeON(b *testing.B) {
	sizes := []int{10, 100, 1000, 10000, 100000}
	for _, size := range sizes {
		data := make([]int, size)
		for i := range data {
			data[i] = i
		}

		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = FromSafe(data)
			}
		})
	}
}

func TestGroupBy(t *testing.T) {
	t.Run("Group by age", func(t *testing.T) {
		type Person struct {
			Age  int
			Name string
		}

		people := []Person{
			{25, "Alice"},
			{30, "Bob"},
			{25, "Charlie"},
			{30, "David"},
		}

		result := GroupBy(From(people), func(p Person) int { return p.Age }).ToSlice()

		if len(result) != 2 {
			t.Errorf("Expected 2 groups, got %d", len(result))
		}

		// Build a map for easier checking
		groupMap := make(map[int][]Person)
		for _, kv := range result {
			groupMap[kv.Key] = kv.Value
		}

		// Check group 25
		group25, ok := groupMap[25]
		if !ok {
			t.Errorf("Expected group with key 25")
		}
		if len(group25) != 2 {
			t.Errorf("Expected 2 people in group 25, got %d", len(group25))
		}

		// Check group 30
		group30, ok := groupMap[30]
		if !ok {
			t.Errorf("Expected group with key 30")
		}
		if len(group30) != 2 {
			t.Errorf("Expected 2 people in group 30, got %d", len(group30))
		}
	})

	t.Run("Group by first letter", func(t *testing.T) {
		words := []string{"apple", "banana", "apricot", "blueberry", "avocado"}

		result := GroupBy(From(words), func(s string) rune { return rune(s[0]) }).ToSlice()

		if len(result) != 2 {
			t.Errorf("Expected 2 groups, got %d", len(result))
		}

		groupMap := make(map[rune][]string)
		for _, kv := range result {
			groupMap[kv.Key] = kv.Value
		}

		groupA, ok := groupMap['a']
		if !ok {
			t.Errorf("Expected group with key 'a'")
		}
		if len(groupA) != 3 {
			t.Errorf("Expected 3 words in group 'a', got %d", len(groupA))
		}

		groupB, ok := groupMap['b']
		if !ok {
			t.Errorf("Expected group with key 'b'")
		}
		if len(groupB) != 2 {
			t.Errorf("Expected 2 words in group 'b', got %d", len(groupB))
		}
	})

	t.Run("Group empty stream", func(t *testing.T) {
		result := GroupBy(Empty[int](), func(x int) int { return x }).ToSlice()

		if len(result) != 0 {
			t.Errorf("Expected empty result, got %v", result)
		}
	})

	t.Run("Group with single element", func(t *testing.T) {
		input := []int{42}
		result := GroupBy(From(input), func(x int) int { return x }).ToSlice()

		if len(result) != 1 {
			t.Errorf("Expected 1 group, got %d", len(result))
		}

		if result[0].Key != 42 {
			t.Errorf("Expected key 42, got %d", result[0].Key)
		}

		if len(result[0].Value) != 1 || result[0].Value[0] != 42 {
			t.Errorf("Expected value [42], got %v", result[0].Value)
		}
	})

	t.Run("GroupBy with chaining", func(t *testing.T) {
		type Item struct {
			Category string
			Value    int
		}

		items := []Item{
			{"A", 1},
			{"A", 2},
			{"B", 3},
			{"A", 4},
		}

		// Group by category and then filter groups with more than 2 items
		result := GroupBy(
			From(items),
			func(i Item) string { return i.Category },
		).Where(func(kv KeyValue[string, []Item]) bool {
			return len(kv.Value) > 2
		}).ToSlice()

		if len(result) != 1 {
			t.Errorf("Expected 1 group, got %d", len(result))
		}

		if result[0].Key != "A" {
			t.Errorf("Expected key 'A', got %s", result[0].Key)
		}

		if len(result[0].Value) != 3 {
			t.Errorf("Expected 3 items in group 'A', got %d", len(result[0].Value))
		}
	})
}
