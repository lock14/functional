package main

import (
	"fmt"
	"github.com/lock14/functional/channels"
)

func main() {
	channel := channels.FromSlice([]int{1, 1, 2, 2, 3, 3, 4, 4})
	distinct := channels.Distinct(channel)
	slice := channels.ToSlice(distinct)
	fmt.Println(slice)
}
