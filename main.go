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
