package predicate

import (
	"testing"
)

func BenchmarkIsNil(b *testing.B) {
	b.ReportAllocs()
	val := 42
	ptr := &val
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = IsNil(ptr)
	}
}

func BenchmarkNot(b *testing.B) {
	b.ReportAllocs()
	isEven := func(n int) bool { return n%2 == 0 }
	isOdd := Not(isEven)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = isOdd(i)
	}
}
