package main

import (
	"fmt"
	"github.com/lock14/functional/channel"
)

func main() {
	items := channel.FromSlice([]int{1, 1, 2, 2, 3, 3, 4, 4})
	distinct := channel.Distinct(items)
	slice := channel.ToSlice(distinct)
	fmt.Println(slice)
	fmt.Println(channel.ToSlice(channel.Limit(channel.Generate(func() int { return 1 }), 10)))
}
