package main

import (
	"fmt"

	"github.com/CreateLab/glinq/pkg/glinq"
)

func main() {
	fmt.Println("=== glinq Examples ===")
	fmt.Println()

	// Пример 1: Select - преобразование в тот же тип
	fmt.Println("Example 1: Select (same type transformation)")
	numbers := []int{1, 2, 3, 4, 5}
	doubled := glinq.From(numbers).
		Select(func(x int) int { return x * 2 }).
		ToSlice()
	fmt.Printf("Input: %v\n", numbers)
	fmt.Printf("Doubled (Select): %v\n\n", doubled)

	// Пример 2: Select с Where (chain)
	fmt.Println("Example 2: Where + Select chain")
	numbers2 := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	result := glinq.From(numbers2).
		Where(func(x int) bool { return x > 5 }).
		Select(func(x int) int { return x * 2 }).
		ToSlice()
	fmt.Printf("Input: %v\n", numbers2)
	fmt.Printf("Filtered (> 5) and mapped (x * 2): %v\n\n", result)

	// Пример 3: Map - преобразование в другой тип (int -> string)
	fmt.Println("Example 3: Map (int -> string)")
	numbers3 := []int{1, 2, 3, 4, 5}
	strings := glinq.Map(
		glinq.From(numbers3),
		func(x int) string { return fmt.Sprintf("Number: %d", x) },
	).ToSlice()
	fmt.Printf("Input: %v\n", numbers3)
	fmt.Printf("Strings: %v\n\n", strings)

	// Пример 4: Map с Where (combined)
	fmt.Println("Example 4: Where + Map (different type)")
	numbers4 := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	filtered := glinq.Map(
		glinq.From(numbers4).Where(func(x int) bool { return x%2 == 0 }),
		func(x int) string { return fmt.Sprintf("Even: %d", x) },
	).ToSlice()
	fmt.Printf("Input: %v\n", numbers4)
	fmt.Printf("Even numbers as strings: %v\n\n", filtered)

	// Пример 5: Map в структуру
	fmt.Println("Example 5: Map to struct")
	type User struct {
		ID   int
		Name string
	}

	ids := []int{1, 2, 3}
	users := glinq.Map(
		glinq.From(ids),
		func(id int) User {
			return User{ID: id, Name: fmt.Sprintf("User%d", id)}
		},
	).ToSlice()
	fmt.Printf("Input IDs: %v\n", ids)
	fmt.Printf("Users: %+v\n\n", users)

	// Пример 6: Работа с map
	fmt.Println("Example 6: Filter a map by values")
	data := map[string]int{
		"apple":  5,
		"banana": 3,
		"orange": 8,
	}
	filteredMap := glinq.ToMap(
		glinq.FromMap(data).
			Where(func(kv glinq.KeyValue[string, int]) bool {
				return kv.Value > 4
			}),
	)
	fmt.Printf("Input map: %v\n", data)
	fmt.Printf("Filtered map (values > 4): %v\n\n", filteredMap)

	// Пример 7: Работа с Keys и Values
	fmt.Println("Example 7: Extract keys and values from map")
	data2 := map[string]int{
		"x": 10,
		"y": 20,
		"z": 30,
	}
	keys := glinq.Keys(glinq.FromMap(data2)).ToSlice()
	values := glinq.Values(glinq.FromMap(data2)).ToSlice()
	fmt.Printf("Original map: %v\n", data2)
	fmt.Printf("Keys: %v\n", keys)
	fmt.Printf("Values: %v\n\n", values)

	// Пример 8: Демонстрация lazy evaluation
	fmt.Println("Example 8: Lazy evaluation demonstration")
	fmt.Println("Creating stream with Range(1, 10) -> Where(x > 5) -> Take(3)")
	fmt.Println("Only the necessary elements are processed:")
	fmt.Println()
	lazyResult := glinq.Range(1, 10).
		Where(func(x int) bool {
			fmt.Printf("  Checking %d\n", x)
			return x > 5
		}).
		Take(3).
		ToSlice()
	fmt.Printf("Result: %v\n\n", lazyResult)

	// Пример 9: Lazy evaluation с Map
	fmt.Println("Example 9: Lazy evaluation with Map")
	fmt.Println("Notice: only 3 elements are processed due to Take(3)")
	lazy := glinq.Map(
		glinq.Range(1, 20).
			Where(func(x int) bool {
				fmt.Printf("  Checking %d\n", x)
				return x%2 == 0
			}).
			Take(3),
		func(x int) string { return fmt.Sprintf("Result: %d", x) },
	).ToSlice()
	fmt.Printf("Lazy result: %v\n\n", lazy)

	// Пример 10: Count и Any
	fmt.Println("Example 10: Count and Any operations")
	numbers5 := []int{1, 2, 3, 4, 5}
	count := glinq.From(numbers5).
		Where(func(x int) bool { return x%2 == 0 }).
		Count()
	fmt.Printf("Count of even numbers in %v: %d\n", numbers5, count)
	hasEven := glinq.From(numbers5).Any(func(x int) bool { return x%2 == 0 })
	fmt.Printf("Has even numbers: %v\n\n", hasEven)

	// Пример 11: All и ForEach
	fmt.Println("Example 11: All and ForEach operations")
	numbers6 := []int{2, 4, 6, 8}
	allEven := glinq.From(numbers6).All(func(x int) bool { return x%2 == 0 })
	fmt.Printf("All numbers in %v are even: %v\n", numbers6, allEven)
	fmt.Print("Doubling each number: ")
	glinq.From(numbers6).ForEach(func(x int) {
		fmt.Printf("%d ", x*2)
	})
	fmt.Println()

	// Пример 12: Сложная цепочка операций
	fmt.Println("Example 12: Complex chained operations")
	numbers7 := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	complex := glinq.From(numbers7).
		Skip(2).
		Where(func(x int) bool { return x%2 == 0 }).
		Select(func(x int) int { return x * x }).
		Take(3).
		ToSlice()
	fmt.Printf("Input: %v\n", numbers7)
	fmt.Printf("Skip(2) -> Where(even) -> Select(x*x) -> Take(3): %v\n\n", complex)

	// Пример 13: Комбинация Skip, Map, Take
	fmt.Println("Example 13: Skip + Map + Take")
	numbers8 := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	combined := glinq.Map(
		glinq.From(numbers8).
			Skip(3).
			Where(func(x int) bool { return x < 8 }).
			Take(3),
		func(x int) string { return fmt.Sprintf("[%d]", x*10) },
	).ToSlice()
	fmt.Printf("Input: %v\n", numbers8)
	fmt.Printf("Skip(3) -> Where(x < 8) -> Take(3) -> Map: %v\n", combined)

	fmt.Println("\n=== End of Examples ===")
}
