package iterator

import (
	"slices"
	"testing"
)

func BenchmarkMap(b *testing.B) {
	b.ReportAllocs()
	f := func(x int) int { return x * 2 }

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		seq := Map(Range(0, 1000), f)
		for range seq {
		}
	}
}

func BenchmarkFilter(b *testing.B) {
	b.ReportAllocs()
	p := func(x int) bool { return x%2 == 0 }

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		seq := Filter(Range(0, 1000), p)
		for range seq {
		}
	}
}

func BenchmarkDistinct(b *testing.B) {
	b.ReportAllocs()
	data := make([]int, 1000)
	for i := range data {
		data[i] = i % 100
	}
	seq := slices.Values(data)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d := Distinct(seq)
		for range d {
		}
	}
}
