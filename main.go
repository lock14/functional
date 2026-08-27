package main

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"

	"github.com/lock14/functional/channel"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	generator := channel.Generate(ctx, rand.Int)
	for s := range channel.ParallelMap(ctx, 4, channel.Limit(ctx, generator, 3), strconv.Itoa) {
		fmt.Println(s)
	}
}

func chanTest() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	items := channel.FromSlice([]int{1, 1, 2, 2, 3, 3, 4, 4})
	distinct := channel.Distinct(ctx, items)
	slice := channel.ToSlice(ctx, distinct)
	fmt.Println(slice)
	generator := channel.Generate(ctx, func() int { return 1 })
	fmt.Println(channel.ToSlice(ctx, channel.Limit(ctx, generator, 10)))

	// to test cancellation
	ctx2, cancel2 := context.WithCancel(context.Background())
	generator2 := channel.Generate(ctx2, func() int { return 1 })
	cancel2()
	val, ok := <-generator2
	fmt.Printf("val: %v, ok: %v\n", val, ok)

	fmt.Println(channel.ToSlice(ctx, channel.Of(1, 2, 3)))
	fmt.Println(channel.ToSlice(ctx, channel.Zip(ctx, channel.Of(1, 2, 3), channel.Of("bob", "mary"))))
	fmt.Println(channel.Join(ctx, channel.Of("[", "]"), channel.Join(ctx, channel.Of("bob", "mary", "jain"), ", ")))
	fmt.Println(channel.Join(ctx, channel.Of("bob"), ", "))

	// partition now returns chan []T, so to view we just range or use a custom ToSlice for slices
	partChan := channel.Partition(ctx, channel.Range(ctx, 0, 10), 3)
	for p := range partChan {
		fmt.Println(p)
	}
}
