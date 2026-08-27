package iterator_test

import (
	"fmt"

	"github.com/lock14/functional/iterator"
)

func ExampleRange() {
	for n := range iterator.Range(1, 4) {
		fmt.Println(n)
	}

	// Output:
	// 1
	// 2
	// 3
}

func ExampleFilter() {
	numbers := iterator.Range(1, 7)
	evens := iterator.Filter(numbers, func(n int) bool { return n%2 == 0 })
	for n := range evens {
		fmt.Println(n)
	}

	// Output:
	// 2
	// 4
	// 6
}

func ExampleLimit() {
	numbers := iterator.Range(1, 100)
	firstThree := iterator.Limit(numbers, 3)
	for n := range firstThree {
		fmt.Println(n)
	}

	// Output:
	// 1
	// 2
	// 3
}

func ExampleDistinct() {
	items := iterator.Of(1, 2, 2, 3, 1, 4)
	distinct := iterator.Distinct(items)
	for n := range distinct {
		fmt.Println(n)
	}

	// Output:
	// 1
	// 2
	// 3
	// 4
}
