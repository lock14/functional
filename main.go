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
	generator, closeGenerator := channel.Generate(func() int { return 1 })
	fmt.Println(channel.ToSlice(channel.Limit(generator, 10)))
	closeGenerator()
	val, ok := <-generator
	fmt.Printf("val: %v, ok: %v\n", val, ok)
	fmt.Println(channel.ToSlice(channel.Of(1, 2, 3)))
	fmt.Println(channel.ToSlice(channel.Zip(channel.Of(1, 2, 3), channel.Of("bob", "mary"))))
	fmt.Println(channel.Join(channel.Of("[", "]"), channel.Join(channel.Of("bob", "mary", "jain"), ", ")))
	fmt.Println(channel.Join(channel.Of("bob"), ", "))
}
