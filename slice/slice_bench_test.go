package slice

import (
	"testing"
)

func BenchmarkMap(b *testing.B) {
	b.ReportAllocs()
	data := make([]int, 1000)
	for i := range data {
		data[i] = i
	}
	f := func(x int) int { return x * 2 }

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Map(data, f)
	}
}

func BenchmarkFilter(b *testing.B) {
	b.ReportAllocs()
	data := make([]int, 1000)
	for i := range data {
		data[i] = i
	}
	p := func(x int) bool { return x%2 == 0 }

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Filter(data, p)
	}
}

func BenchmarkFoldLeft(b *testing.B) {
	b.ReportAllocs()
	data := make([]int, 1000)
	for i := range data {
		data[i] = i
	}
	f := func(acc, x int) int { return acc + x }

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = FoldLeft(data, f, 0)
	}
}

func BenchmarkPartition(b *testing.B) {
	b.ReportAllocs()
	data := make([]int, 1000)
	for i := range data {
		data[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Partition(data, 50)
	}
}
