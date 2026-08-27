package channel

import (
	"context"
	"testing"
)

func TestParallelMap(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c := ParallelMap(ctx, 4, Range(ctx, 1, 4), func(i int) int { return i * 2 })
	s := Sorted(ctx, c)
	res := ToSlice(ctx, s)
	if len(res) != 3 || res[0] != 2 || res[1] != 4 || res[2] != 6 {
		t.Errorf("ParallelMap failed: %v", res)
	}

	c0 := ParallelMap(ctx, 0, Range(ctx, 1, 4), func(i int) int { return i * 2 })
	<-c0
}

func TestParallelFlatten(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Flatten channels of channels
	c := ParallelFlatten(ctx, 4, Of(Of(1, 2), Of(3, 4)))
	s := Sorted(ctx, c)
	res := ToSlice(ctx, s)
	if len(res) != 4 || res[0] != 1 || res[3] != 4 {
		t.Errorf("ParallelFlatten failed: %v", res)
	}

	c0 := ParallelFlatten(ctx, 0, Of(Of(1, 2)))
	<-c0
}

func TestParallelFlatMap(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c := ParallelFlatMap(ctx, 4, Range(ctx, 1, 3), func(i int) chan int { return Of(i, i*10) })
	s := Sorted(ctx, c)
	res := ToSlice(ctx, s)
	if len(res) != 4 || res[0] != 1 || res[1] != 2 || res[2] != 10 || res[3] != 20 {
		t.Errorf("ParallelFlatMap failed: %v", res)
	}
}

func TestParallelFilter(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c := ParallelFilter(ctx, 4, Range(ctx, 1, 5), func(i int) bool { return i%2 == 0 })
	s := Sorted(ctx, c)
	res := ToSlice(ctx, s)
	if len(res) != 2 || res[0] != 2 || res[1] != 4 {
		t.Errorf("ParallelFilter failed: %v", res)
	}

	c0 := ParallelFilter(ctx, 0, Range(ctx, 1, 5), func(i int) bool { return i%2 == 0 })
	<-c0
}
