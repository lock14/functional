package main

import (
	"fmt"
	"github.com/lock14/functional/channels"
)

func main() {
	fmt.Println(
		channels.ToSlice(
			channels.Distinct(
				channels.FromSlice([]int{1, 1, 2, 2, 3, 3, 4, 4}),
			),
		),
	)
}
