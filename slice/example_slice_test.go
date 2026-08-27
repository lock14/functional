package slice_test

import (
	"fmt"
	"strconv"

	"github.com/lock14/functional/slice"
)

func ExampleMap() {
	numbers := []int{1, 2, 3}
	strings := slice.Map(numbers, strconv.Itoa)
	fmt.Println(strings)

	// Output:
	// [1 2 3]
}

func ExampleFilter() {
	numbers := []int{1, 2, 3, 4, 5, 6}
	evens := slice.Filter(numbers, func(n int) bool { return n%2 == 0 })
	fmt.Println(evens)

	// Output:
	// [2 4 6]
}

func ExampleFoldLeft() {
	numbers := []int{1, 2, 3, 4}
	sum := slice.FoldLeft(numbers, func(acc, n int) int { return acc + n }, 0)
	fmt.Println(sum)

	// Output:
	// 10
}

func ExampleZip() {
	numbers := []int{1, 2, 3}
	names := []string{"Alice", "Bob", "Charlie"}
	pairs := slice.Zip(numbers, names)
	fmt.Printf("%+v\n", pairs)

	// Output:
	// [{Fst:1 Snd:Alice} {Fst:2 Snd:Bob} {Fst:3 Snd:Charlie}]
}

func ExamplePartition() {
	numbers := []int{1, 2, 3, 4, 5}
	chunks := slice.Partition(numbers, 2)
	fmt.Println(chunks)

	// Output:
	// [[1 2] [3 4] [5]]
}
