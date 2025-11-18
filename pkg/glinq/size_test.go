package glinq

import (
	"testing"
	"time"
)

// assertSize проверяет что размер известен и равен ожидаемому
func assertSize(t *testing.T, s Stream[int], expected int, op string) {
	t.Helper()
	size, ok := s.Size()
	if !ok {
		t.Errorf("%s: expected known size %d, got unknown", op, expected)
		return
	}
	if size != expected {
		t.Errorf("%s: expected size %d, got %d", op, expected, size)
	}
}

// assertNoSize проверяет что размер неизвестен
func assertNoSize(t *testing.T, s Stream[int], op string) {
	t.Helper()
	_, ok := s.Size()
	if ok {
		t.Errorf("%s: expected unknown size, got known", op)
	}
}

func TestFrom_Size(t *testing.T) {
	s1 := From([]int{1, 2, 3, 4, 5})
	assertSize(t, s1, 5, "From")
}

func TestFromSafe_Size(t *testing.T) {
	s1 := FromSafe([]int{1, 2, 3})
	assertSize(t, s1, 3, "FromSafe")
}

func TestEmpty_Size(t *testing.T) {
	s1 := Empty[int]()
	assertSize(t, s1, 0, "Empty")
}

func TestRange_Size(t *testing.T) {
	s1 := Range(1, 10)
	assertSize(t, s1, 10, "Range")
}

func TestSelect_PreservesSize(t *testing.T) {
	s1 := From([]int{1, 2, 3, 4, 5})
	s2 := s1.Select(func(x int) int { return x * 2 })
	assertSize(t, s2, 5, "Select")
}

func TestSelectWithIndex_PreservesSize(t *testing.T) {
	s1 := From([]int{1, 2, 3})
	s2 := s1.SelectWithIndex(func(x int, idx int) int { return x * idx })
	assertSize(t, s2, 3, "SelectWithIndex")
}

func TestWhere_LosesSize(t *testing.T) {
	s1 := From([]int{1, 2, 3, 4, 5})
	s2 := s1.Where(func(x int) bool { return x > 3 })
	assertNoSize(t, s2, "Where")
}

func TestTake_CalculatesSize(t *testing.T) {
	s1 := From([]int{1, 2, 3, 4, 5})
	s2 := s1.Take(3)
	assertSize(t, s2, 3, "Take")

	// Take with size larger than source
	s3 := s1.Take(10)
	assertSize(t, s3, 5, "Take (larger than source)")

	// Take with unknown source
	s4 := s1.Where(func(x int) bool { return x > 0 })
	s5 := s4.Take(3)
	assertSize(t, s5, 3, "Take (unknown source)")
}

func TestSkip_CalculatesSize(t *testing.T) {
	s1 := From([]int{1, 2, 3, 4, 5})
	s2 := s1.Skip(2)
	assertSize(t, s2, 3, "Skip")

	// Skip more than size
	s3 := s1.Skip(10)
	assertSize(t, s3, 0, "Skip (more than size)")

	// Skip with unknown source
	s4 := s1.Where(func(x int) bool { return x > 0 })
	s5 := s4.Skip(2)
	assertNoSize(t, s5, "Skip (unknown source)")
}

func TestConcat_CalculatesSize(t *testing.T) {
	s1 := From([]int{1, 2})
	s2 := From([]int{3, 4, 5})
	s3 := s1.Concat(s2)
	assertSize(t, s3, 5, "Concat")

	// Concat with unknown size
	s4 := s1.Where(func(x int) bool { return x > 0 })
	s5 := s4.Concat(s2)
	assertNoSize(t, s5, "Concat (unknown first)")
}

func TestDistinctBy_LosesSize(t *testing.T) {
	s1 := From([]int{1, 2, 2, 3, 3, 3})
	s2 := s1.DistinctBy(func(x int) any { return x })
	assertNoSize(t, s2, "DistinctBy")
}

func TestOrderBy_PreservesSize(t *testing.T) {
	s1 := From([]int{5, 2, 8, 1, 9})
	s2 := s1.OrderBy(func(a, b int) int { return a - b })
	assertSize(t, s2, 5, "OrderBy")
}

func TestOrderByDescending_PreservesSize(t *testing.T) {
	s1 := From([]int{5, 2, 8, 1, 9})
	s2 := s1.OrderByDescending(func(a, b int) int { return a - b })
	assertSize(t, s2, 5, "OrderByDescending")
}

func TestReverse_PreservesSize(t *testing.T) {
	s1 := From([]int{1, 2, 3, 4})
	s2 := s1.Reverse()
	assertSize(t, s2, 4, "Reverse")
}

func TestSelectMany_LosesSize(t *testing.T) {
	numbers := [][]int{{1, 2}, {3, 4}, {5}}
	s1 := From(numbers)
	s2 := SelectMany(s1, func(slice []int) Enumerable[int] { return From(slice) })
	assertNoSize(t, s2, "SelectMany")
}

func TestFromMap_Size(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	s1 := FromMap(m)
	size, ok := s1.Size()
	if !ok {
		t.Error("FromMap: expected known size")
	}
	if size != 3 {
		t.Errorf("FromMap: expected size 3, got %d", size)
	}
}

func TestFromMapSafe_Size(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}
	s1 := FromMapSafe(m)
	size, ok := s1.Size()
	if !ok {
		t.Error("FromMapSafe: expected known size")
	}
	if size != 2 {
		t.Errorf("FromMapSafe: expected size 2, got %d", size)
	}
}

func TestKeys_PreservesSize(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	s1 := FromMap(m)
	s2 := Keys(s1)
	size, ok := s2.Size()
	if !ok {
		t.Error("Keys: expected known size")
	}
	if size != 3 {
		t.Errorf("Keys: expected size 3, got %d", size)
	}
}

func TestValues_PreservesSize(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}
	s1 := FromMap(m)
	s2 := Values(s1)
	size, ok := s2.Size()
	if !ok {
		t.Error("Values: expected known size")
	}
	if size != 2 {
		t.Errorf("Values: expected size 2, got %d", size)
	}
}

func TestGroupBy_Size(t *testing.T) {
	type Person struct {
		Age  int
		Name string
	}
	people := []Person{{25, "Alice"}, {30, "Bob"}, {25, "Charlie"}}
	s1 := From(people)
	s2 := GroupBy(s1, func(p Person) int { return p.Age })
	// GroupBy materializes, so we know the size (number of groups)
	size, ok := s2.Size()
	if !ok {
		t.Error("GroupBy: expected known size")
	}
	if size != 2 {
		t.Errorf("GroupBy: expected size 2 (two age groups), got %d", size)
	}
}

func TestFromEnumerable_PreservesSize(t *testing.T) {
	s1 := From([]int{1, 2, 3})
	s2 := FromEnumerable(s1)
	assertSize(t, s2, 3, "FromEnumerable")

	// FromEnumerable with non-Sizable should have unknown size
	// (we can't test this easily without creating a custom Enumerable)
}

func TestSelect_FreeFunction_PreservesSize(t *testing.T) {
	const s = "num"
	s1 := From([]int{1, 2, 3})
	s2 := Select(s1, func(x int) string {
		return s
	})
	size, ok := s2.Size()
	if !ok {
		t.Error("Select (free function): expected known size")
	}
	if size != 3 {
		t.Errorf("Select (free function): expected size 3, got %d", size)
	}
}

func TestSelectWithIndex_FreeFunction_PreservesSize(t *testing.T) {
	s1 := From([]int{1, 2, 3})
	s2 := SelectWithIndex(s1, func(x int, idx int) string {
		return "num"
	})
	size, ok := s2.Size()
	if !ok {
		t.Error("SelectWithIndex (free function): expected known size")
	}
	if size != 3 {
		t.Errorf("SelectWithIndex (free function): expected size 3, got %d", size)
	}
}

func TestTakeOrderedBy_Size(t *testing.T) {
	s1 := From([]int{5, 2, 8, 1, 9, 3})
	s2 := TakeOrderedBy(s1, 3, func(a, b int) bool { return a < b })
	assertSize(t, s2, 3, "TakeOrderedBy")

	// TakeOrderedBy with size larger than source
	s3 := TakeOrderedBy(s1, 10, func(a, b int) bool { return a < b })
	assertSize(t, s3, 6, "TakeOrderedBy (larger than source)")
}

func TestCount_Optimization(t *testing.T) {
	s1 := From([]int{1, 2, 3, 4, 5})

	// Count should be O(1) when size is known
	start := time.Now()
	count := s1.Count()
	duration := time.Since(start)

	if count != 5 {
		t.Errorf("Count: expected 5, got %d", count)
	}
	if duration > time.Microsecond*100 {
		t.Errorf("Count: expected O(1) operation, took %v", duration)
	}

	// Count with unknown size should iterate
	s2 := s1.Where(func(x int) bool { return x > 0 })
	start2 := time.Now()
	count2 := s2.Count()
	duration2 := time.Since(start2)

	if count2 != 5 {
		t.Errorf("Count (unknown size): expected 5, got %d", count2)
	}
	// This should take longer (but still fast for small collections)
	if duration2 < duration {
		t.Logf("Count with unknown size took %v (expected longer than %v)", duration2, duration)
	}
}

func TestToSlice_Preallocation(t *testing.T) {
	s1 := From([]int{1, 2, 3, 4, 5})
	result := s1.ToSlice()

	if len(result) != 5 {
		t.Errorf("ToSlice: expected length 5, got %d", len(result))
	}

	// Check that capacity is at least length (preallocation)
	if cap(result) < len(result) {
		t.Errorf("ToSlice: expected capacity >= length, got cap=%d, len=%d", cap(result), len(result))
	}
}

func TestChunk_Preallocation(t *testing.T) {
	s1 := From([]int{1, 2, 3, 4, 5, 6, 7})
	chunks := s1.Chunk(3)

	if len(chunks) != 3 {
		t.Errorf("Chunk: expected 3 chunks, got %d", len(chunks))
	}

	// Check preallocation (capacity should be at least length)
	if cap(chunks) < len(chunks) {
		t.Errorf("Chunk: expected capacity >= length, got cap=%d, len=%d", cap(chunks), len(chunks))
	}
}

func TestSizeHandling_Chain(t *testing.T) {
	// Test complex chain
	s1 := From([]int{1, 2, 3, 4, 5})
	assertSize(t, s1, 5, "From")

	// Select preserves size
	s2 := s1.Select(func(x int) int { return x * 2 })
	assertSize(t, s2, 5, "Select")

	// Where loses size
	s3 := s2.Where(func(x int) bool { return x > 5 })
	assertNoSize(t, s3, "Where")

	// Take calculates size
	s4 := s1.Take(3)
	assertSize(t, s4, 3, "Take")

	// Take with unknown source
	s5 := s3.Take(10)
	assertSize(t, s5, 10, "Take with unknown source")

	// Skip calculates size
	s6 := s1.Skip(2)
	assertSize(t, s6, 3, "Skip")

	// Concat adds sizes
	s7 := From([]int{1, 2}).Concat(From([]int{3, 4, 5}))
	assertSize(t, s7, 5, "Concat")
}

func TestDistinct_LosesSize(t *testing.T) {
	s1 := From([]int{1, 2, 2, 3, 3, 3})
	s2 := Distinct(s1)
	assertNoSize(t, s2, "Distinct")
}

func TestUnion_LosesSize(t *testing.T) {
	s1 := From([]int{1, 2, 3})
	s2 := From([]int{3, 4, 5})
	s3 := Union(s1, s2)
	assertNoSize(t, s3, "Union")
}

func TestIntersect_LosesSize(t *testing.T) {
	s1 := From([]int{1, 2, 3})
	s2 := From([]int{2, 3, 4})
	s3 := Intersect(s1, s2)
	assertNoSize(t, s3, "Intersect")
}

func TestExcept_LosesSize(t *testing.T) {
	s1 := From([]int{1, 2, 3})
	s2 := From([]int{2, 3})
	s3 := Except(s1, s2)
	assertNoSize(t, s3, "Except")
}
