package iterator

import (
	"cmp"
	"errors"
	"github.com/lock14/functional/slice"
	"golang.org/x/exp/constraints"
	"iter"
	"slices"
)

// Monad represents any type that can use the `+` operator and whose zero
// value is the identity element the `+` operator
type Monad interface {
	constraints.Integer | constraints.Float | constraints.Complex | ~string
}

func Map[T, U any](itr iter.Seq[T], f func(T) U) iter.Seq[U] {
	return func(yield func(U) bool) {
		for t := range itr {
			if !yield(f(t)) {
				break
			}
		}
	}
}

func Flatten[T any](itrs iter.Seq[iter.Seq[T]]) iter.Seq[T] {
	return func(yield func(T) bool) {
	Loop:
		for itr := range itrs {
			for t := range itr {
				if !yield(t) {
					break Loop
				}
			}
		}
	}
}

func FlatMap[T, U any](iter iter.Seq[T], f func(T) iter.Seq[U]) iter.Seq[U] {
	return Flatten(Map(iter, f))
}

func Filter[T any](itr iter.Seq[T], p func(T) bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		for t := range itr {
			if p(t) {
				if !yield(t) {
					break
				}
			}
		}
	}
}

func FoldLeft[T, U any](itr iter.Seq[T], f func(U, T) U, u U) U {
	result := u
	for t := range itr {
		result = f(result, t)
	}
	return result
}

func FoldRight[T, U any](itr iter.Seq[T], f func(T, U) U, u U) U {
	result := u
	for t := range itr {
		result = f(t, FoldRight[T, U](itr, f, u))
	}
	return result
}

func Reduce[T any](itr iter.Seq[T], f func(T, T) T, t T) T {
	return FoldLeft(itr, f, t)
}

func Sum[M Monad](itr iter.Seq[M]) M {
	var identity M
	return Reduce(itr, func(a, b M) M { return a + b }, identity)
}

func JoinErrs(itr iter.Seq[error]) error {
	return Reduce(itr, func(e1, e2 error) error { return errors.Join(e1, e2) }, nil)
}

func Join[T ~string](itr iter.Seq[T], sep T) T {
	first := true
	var result T
	for t := range itr {
		if first {
			first = false
		} else {
			t += sep
		}
		result += t
	}
	return result
}

func Zip[T, U any](itr1 iter.Seq[T], itr2 iter.Seq[U]) iter.Seq2[T, U] {
	return func(yield func(T, U) bool) {
		next1, stop1 := iter.Pull(itr1)
		defer stop1()
		next2, stop2 := iter.Pull(itr2)
		defer stop2()

		t, ok1 := next1()
		u, ok2 := next2()
		for ok1 && ok2 && yield(t, u) {
			t, ok1 = next1()
			u, ok2 = next2()
		}
	}
}

func UnZip[T, U any](itr iter.Seq2[T, U]) (iter.Seq[T], iter.Seq[U]) {
	// TODO
	return nil, nil
}

func Sorted[T cmp.Ordered](itr iter.Seq[T]) iter.Seq[T] {
	return slices.Values(slices.Sorted(itr))
}

func Distinct[T comparable](itr iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		set := make(map[T]struct{})
		for t := range itr {
			if _, ok := set[t]; !ok {
				set[t] = struct{}{}
				if !yield(t) {
					break
				}
			}
		}
	}
}

func Iterate[T any](seed T, hasNext func(T) bool, next func(T) T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for cur := seed; hasNext(cur); cur = next(cur) {
			if !yield(cur) {
				break
			}
		}
	}
}

func Range[T constraints.Integer](startInclusive, endExclusive T) iter.Seq[T] {
	return Iterate(startInclusive, func(t T) bool { return t < endExclusive }, func(t T) T { t++; return t })
}

func RangeClosed[T constraints.Integer](startInclusive, endInclusive T) iter.Seq[T] {
	return Iterate(startInclusive, func(t T) bool { return t <= endInclusive }, func(t T) T { t++; return t })
}

func Limit[T any](itr iter.Seq[T], max int64) iter.Seq[T] {
	return func(yield func(T) bool) {
		var count int64
		for t := range itr {
			if count == max || !yield(t) {
				break
			}
			count++
		}
	}
}

func Skip[T any](itr iter.Seq[T], n int64) iter.Seq[T] {
	return func(yield func(T) bool) {
		var count int64
		for t := range itr {
			if count >= n {
				if !yield(t) {
					break
				}
			}
			count++
		}
	}
}

func AllMatch[T any](itr iter.Seq[T], p func(T) bool) bool {
	return Reduce(Map(itr, p), func(t1, t2 bool) bool { return t1 && t2 }, true)
}

func AnyMatch[T any](itr iter.Seq[T], p func(T) bool) bool {
	return Reduce(Map(itr, p), func(t1, t2 bool) bool { return t1 || t2 }, false)
}

func Count[T any](itr iter.Seq[T]) int64 {
	return Sum(Map(itr, func(t T) int64 { return 1 }))
}

func Concat[T any](itr1, itr2 iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for t := range itr1 {
			if !yield(t) {
				break
			}
		}
		for t := range itr2 {
			if !yield(t) {
				break
			}
		}
	}
}

func Peek[T any](itr iter.Seq[T], consumer func(T)) iter.Seq[T] {
	return func(yield func(T) bool) {
		for t := range itr {
			consumer(t)
			if !yield(t) {
				break
			}
		}
	}
}

func Of[T any](ts ...T) iter.Seq[T] {
	return slices.Values(ts)
}

func Partition[T any](itr iter.Seq[T], size int) iter.Seq[iter.Seq[T]] {
	return slices.Values[[]iter.Seq[T]](slice.Map(slice.Partition(slices.Collect(itr), size), slices.Values))
}
