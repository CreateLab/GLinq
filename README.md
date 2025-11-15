# glinq

LINQ-подобный API для Go с поддержкой отложенного выполнения (lazy evaluation).

glinq предоставляет функциональный подход к работе со слайсами и картами в Go, вдохновленный Microsoft LINQ.
Все операции выполняются лениво и не начинаются до вызова терминальной операции.

## Особенности

- **Lazy Evaluation**: Все промежуточные операции выполняются только при материализации результата
- **Type Safe**: Полная поддержка generics (Go 1.18+)
- **Composable**: Операции легко комбинируются в цепочки
- **Zero Dependencies**: Не требует внешних зависимостей
- **Map Support**: Встроенная поддержка работы с картами

## Установка

```bash
go get github.com/yourusername/glinq
```

## Быстрый старт

```go
package main

import (
    "fmt"
    "github.com/yourusername/glinq/pkg/glinq"
)

func main() {
    numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
    result := glinq.From(numbers).
        Where(func(x int) bool { return x > 5 }).
        Select(func(x int) int { return x * 2 }).
        ToSlice()
    
    fmt.Println(result) // [12, 14, 16, 18, 20]
}
```

## Примеры использования

### Фильтрация (Where)

```go
numbers := []int{1, 2, 3, 4, 5}
evens := glinq.From(numbers).
    Where(func(x int) bool { return x%2 == 0 }).
    ToSlice()
// [2, 4]
```

### Трансформация (Select)

```go
numbers := []int{1, 2, 3}
squared := glinq.From(numbers).
    Select(func(x int) int { return x * x }).
    ToSlice()
// [1, 4, 9]
```

### Ограничение элементов (Take и Skip)

```go
numbers := []int{1, 2, 3, 4, 5}

// Take первые 3 элемента
first3 := glinq.From(numbers).Take(3).ToSlice()
// [1, 2, 3]

// Skip первые 2 элемента
rest := glinq.From(numbers).Skip(2).ToSlice()
// [3, 4, 5]
```

### Работа с картами

```go
data := map[string]int{
    "apple":  5,
    "banana": 3,
    "orange": 8,
}

filtered := glinq.FromMap(data).
    Where(func(kv glinq.KeyValue[string, int]) bool {
        return kv.Value > 4
    }).
    ToMap()
// map[apple:5 orange:8]
```

### Проверка условий (Any и All)

```go
numbers := []int{1, 2, 3, 4, 5}

hasEven := glinq.From(numbers).Any(func(x int) bool { 
    return x%2 == 0 
})
// true

allPositive := glinq.From(numbers).All(func(x int) bool { 
    return x > 0 
})
// true
```

### Подсчет элементов (Count)

```go
numbers := []int{1, 2, 3, 4, 5}
count := glinq.From(numbers).
    Where(func(x int) bool { return x > 2 }).
    Count()
// 3
```

### Выполнение действия для каждого элемента (ForEach)

```go
numbers := []int{1, 2, 3}
glinq.From(numbers).ForEach(func(x int) {
    fmt.Println(x)
})
// 1
// 2
// 3
```

### Получение первого элемента (First)

```go
numbers := []int{1, 2, 3}
first, ok := glinq.From(numbers).First()
// first = 1, ok = true
```

### Демонстрация lazy evaluation

```go
// Благодаря lazy evaluation, фильтр применяется только к необходимым элементам
result := glinq.Range(1, 1000000).
    Where(func(x int) bool { return x%2 == 0 }).
    Take(3).
    ToSlice()
// [2, 4, 6]
// Обработано только ~6 элементов, а не миллион!
```

## Поддерживаемые операции

### Создатели (Creators)

- `From[T any](slice []T) *Stream[T]` - создать Stream из слайса
- `Empty[T any]() *Stream[T]` - создать пустой Stream
- `Range(start, count int) *Stream[int]` - создать Stream целых чисел
- `FromMap[K, V](m map[K]V) *Stream[KeyValue[K, V]]` - создать Stream из карты

### Операторы (Operators)

- `Where(predicate func(T) bool) *Stream[T]` - фильтрация по условию
- `Select[R any](mapper func(T) R) *Stream[R]` - преобразование элементов
- `Take(n int) *Stream[T]` - взять первые n элементов
- `Skip(n int) *Stream[T]` - пропустить первые n элементов

### Терминальные операции (Terminal Operations)

- `ToSlice() []T` - преобразовать Stream в слайс
- `First() (T, bool)` - получить первый элемент
- `Count() int` - подсчитать количество элементов
- `Any(predicate func(T) bool) bool` - проверить наличие элемента
- `All(predicate func(T) bool) bool` - проверить все элементы
- `ForEach(action func(T))` - выполнить действие для каждого элемента

### Вспомогательные функции

- `Keys[K, V](stream *Stream[KeyValue[K, V]]) *Stream[K]` - извлечь ключи
- `Values[K, V](stream *Stream[KeyValue[K, V]]) *Stream[V]` - извлечь значения
- `ToMap[K, V](stream *Stream[KeyValue[K, V]]) map[K]V` - преобразовать в карту

## Требования

- Go 1.18+ (для поддержки generics)

## Тестирование

```bash
go test ./...
```

## Запуск примеров

```bash
go run examples/basic/main.go
```

## Лицензия

MIT

## Автор

your-username
