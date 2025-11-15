package glinq

// Numeric represents numeric types that support addition.
type Numeric interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

// Ordered represents types that can be compared using <, <=, >, >= operators.
type Ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64 |
		~string
}

// Sum calculates the sum of all elements in the Stream.
// Returns zero value if Stream is empty.
//
// Example:
//
//	numbers := []int{1, 2, 3, 4, 5}
//	sum := Sum(From(numbers))
//	// 15
func Sum[T Numeric](s Stream[T]) T {
	var sum T
	for {
		value, ok := s.Next()
		if !ok {
			break
		}
		sum += value
	}
	return sum
}

// Min returns the minimum element in the Stream.
// Returns zero value and false if Stream is empty.
//
// Example:
//
//	numbers := []int{5, 2, 8, 1, 9}
//	min, ok := Min(From(numbers))
//	// min = 1, ok = true
func Min[T Ordered](s Stream[T]) (T, bool) {
	var minVal T
	var found bool

	for {
		value, ok := s.Next()
		if !ok {
			break
		}
		if !found || value < minVal {
			minVal = value
			found = true
		}
	}

	return minVal, found
}

// Max returns the maximum element in the Stream.
// Returns zero value and false if Stream is empty.
//
// Example:
//
//	numbers := []int{5, 2, 8, 1, 9}
//	max, ok := Max(From(numbers))
//	// max = 9, ok = true
func Max[T Ordered](s Stream[T]) (T, bool) {
	var maxVal T
	var found bool

	for {
		value, ok := s.Next()
		if !ok {
			break
		}
		if !found || value > maxVal {
			maxVal = value
			found = true
		}
	}

	return maxVal, found
}
