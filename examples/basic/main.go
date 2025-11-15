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

	// Пример 3: Select - преобразование в другой тип (int -> string)
	fmt.Println("Example 3: Select (int -> string)")
	numbers3 := []int{1, 2, 3, 4, 5}
	strings := glinq.Select(
		glinq.From(numbers3),
		func(x int) string { return fmt.Sprintf("Number: %d", x) },
	).ToSlice()
	fmt.Printf("Input: %v\n", numbers3)
	fmt.Printf("Strings: %v\n\n", strings)

	// Пример 4: Select с Where (combined)
	fmt.Println("Example 4: Where + Select (different type)")
	numbers4 := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	filtered := glinq.Select(
		glinq.From(numbers4).Where(func(x int) bool { return x%2 == 0 }),
		func(x int) string { return fmt.Sprintf("Even: %d", x) },
	).ToSlice()
	fmt.Printf("Input: %v\n", numbers4)
	fmt.Printf("Even numbers as strings: %v\n\n", filtered)

	// Пример 5: Select в структуру
	fmt.Println("Example 5: Select to struct")
	type User struct {
		ID   int
		Name string
	}

	ids := []int{1, 2, 3}
	users := glinq.Select(
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

	// Пример 9: Lazy evaluation с Select
	fmt.Println("Example 9: Lazy evaluation with Select")
	fmt.Println("Notice: only 3 elements are processed due to Take(3)")
	lazy := glinq.Select(
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

	// Пример 13: Комбинация Skip, Select, Take
	fmt.Println("Example 13: Skip + Select + Take")
	numbers8 := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	combined := glinq.Select(
		glinq.From(numbers8).
			Skip(3).
			Where(func(x int) bool { return x < 8 }).
			Take(3),
		func(x int) string { return fmt.Sprintf("[%d]", x*10) },
	).ToSlice()
	fmt.Printf("Input: %v\n", numbers8)
	fmt.Printf("Skip(3) -> Where(x < 8) -> Take(3) -> Select: %v\n\n", combined)

	// Пример 14: SelectWithIndex - метод с индексом
	fmt.Println("Example 14: SelectWithIndex (method with index)")
	numbers9 := []int{1, 2, 3, 4, 5}
	indexed := glinq.From(numbers9).
		SelectWithIndex(func(x int, idx int) int { return x * idx }).
		ToSlice()
	fmt.Printf("Input: %v\n", numbers9)
	fmt.Printf("SelectWithIndex(x * idx): %v\n\n", indexed)

	// Пример 15: SelectWithIndex - функция с индексом (разные типы)
	fmt.Println("Example 15: SelectWithIndex function (different types with index)")
	numbers10 := []int{10, 20, 30}
	indexedStrings := glinq.SelectWithIndex(
		glinq.From(numbers10),
		func(x int, idx int) string { return fmt.Sprintf("Value_%d_at_Index_%d", x, idx) },
	).ToSlice()
	fmt.Printf("Input: %v\n", numbers10)
	fmt.Printf("SelectWithIndex to string: %v\n\n", indexedStrings)

	// Пример 16: SelectWithIndex с Where
	fmt.Println("Example 16: Where + SelectWithIndex")
	numbers11 := []int{1, 2, 3, 4, 5, 6}
	filteredIndexed := glinq.SelectWithIndex(
		glinq.From(numbers11).
			Where(func(x int) bool { return x%2 == 0 }),
		func(x int, idx int) string {
			return fmt.Sprintf("Even[%d]=%d", idx, x)
		},
	).ToSlice()
	fmt.Printf("Input: %v\n", numbers11)
	fmt.Printf("Where(even) -> SelectWithIndex: %v\n\n", filteredIndexed)

	// Пример 17: Aggregate - сумма
	fmt.Println("Example 17: Aggregate - Sum")
	numbers12 := []int{1, 2, 3, 4, 5}
	sum := glinq.From(numbers12).Aggregate(0, func(acc, x int) int { return acc + x })
	fmt.Printf("Input: %v\n", numbers12)
	fmt.Printf("Aggregate sum: %d\n\n", sum)

	// Пример 18: Aggregate - произведение
	fmt.Println("Example 18: Aggregate - Product")
	numbers13 := []int{2, 3, 4}
	product := glinq.From(numbers13).Aggregate(1, func(acc, x int) int { return acc * x })
	fmt.Printf("Input: %v\n", numbers13)
	fmt.Printf("Aggregate product: %d\n\n", product)

	// Пример 19: Aggregate - конкатенация строк
	fmt.Println("Example 19: Aggregate - String concatenation")
	words := []string{"Hello", " ", "World", "!"}
	concatenated := glinq.From(words).Aggregate("", func(acc, x string) string { return acc + x })
	fmt.Printf("Input: %v\n", words)
	fmt.Printf("Aggregate concatenation: '%s'\n\n", concatenated)

	// Пример 20: Aggregate с фильтрацией
	fmt.Println("Example 20: Aggregate with Where filter")
	numbers14 := []int{1, 2, 3, 4, 5, 6, 7, 8}
	sumOfEvens := glinq.From(numbers14).
		Where(func(x int) bool { return x%2 == 0 }).
		Aggregate(0, func(acc, x int) int { return acc + x })
	fmt.Printf("Input: %v\n", numbers14)
	fmt.Printf("Sum of even numbers: %d\n\n", sumOfEvens)

	// Пример 21: Aggregate с кастомным типом
	fmt.Println("Example 21: Aggregate with custom type")
	type Point struct {
		X, Y int
	}
	points := []Point{{1, 2}, {3, 4}, {5, 6}}
	totalPoint := glinq.From(points).Aggregate(
		Point{0, 0},
		func(acc, p Point) Point {
			return Point{acc.X + p.X, acc.Y + p.Y}
		},
	)
	fmt.Printf("Input points: %+v\n", points)
	fmt.Printf("Aggregate sum: Point{X:%d, Y:%d}\n\n", totalPoint.X, totalPoint.Y)

	fmt.Println("\n=== End of Examples ===")
}
