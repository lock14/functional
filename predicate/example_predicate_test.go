package predicate_test

import (
	"fmt"

	"github.com/lock14/functional/predicate"
)

func ExampleIsNil() {
	var ptr *int
	val := 10

	fmt.Println(predicate.IsNil(ptr))
	fmt.Println(predicate.IsNil(&val))

	// Output:
	// true
	// false
}

func ExampleNot() {
	isEven := func(n int) bool { return n%2 == 0 }
	isOdd := predicate.Not(isEven)

	fmt.Println(isOdd(3))
	fmt.Println(isOdd(4))

	// Output:
	// true
	// false
}
