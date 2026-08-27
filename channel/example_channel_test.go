package channel_test

import (
	"context"
	"fmt"
	"strconv"

	"github.com/lock14/functional/channel"
)

func ExampleMap() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	in := channel.Of(1, 2, 3)
	mapped := channel.Map(ctx, in, strconv.Itoa)

	for s := range mapped {
		fmt.Println(s)
	}

	// Output:
	// 1
	// 2
	// 3
}

func ExampleFilter() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	in := channel.Of(1, 2, 3, 4, 5)
	filtered := channel.Filter(ctx, in, func(n int) bool { return n%2 != 0 })

	for s := range filtered {
		fmt.Println(s)
	}

	// Output:
	// 1
	// 3
	// 5
}

func ExampleLimit() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	in := channel.Of(10, 20, 30, 40, 50)
	limited := channel.Limit(ctx, in, 3)

	for s := range limited {
		fmt.Println(s)
	}

	// Output:
	// 10
	// 20
	// 30
}
