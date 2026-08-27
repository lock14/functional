package channel

import (
	"context"
	"testing"
)

func BenchmarkMap(b *testing.B) {
	b.ReportAllocs()
	ctx := context.Background()
	f := func(x int) int { return x * 2 }

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		in := Range(ctx, 0, 1000)
		mapped := Map(ctx, in, f)
		for range mapped {
		}
	}
}

func BenchmarkParallelMap(b *testing.B) {
	b.ReportAllocs()
	ctx := context.Background()
	f := func(x int) int { return x * 2 }

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		in := Range(ctx, 0, 1000)
		mapped := ParallelMap(ctx, 4, in, f)
		for range mapped {
		}
	}
}

func BenchmarkClone(b *testing.B) {
	b.ReportAllocs()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		in := Range(ctx, 0, 100)
		clones := Clone(ctx, in, 2)
		go func() {
			for range clones[0] {
			}
		}()
		for range clones[1] {
		}
	}
}
