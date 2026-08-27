// Package iterator provides lazy functional operations over Go 1.23+ iterators (iter.Seq and iter.Seq2).
package iterator

import (
	"cmp"
	"errors"
	"iter"
	"slices"

	"github.com/lock14/functional/slice"
	"golang.org/x/exp/constraints"
)

// Monoid represents any type that supports the `+` operator and whose zero
// value is the identity element of the `+` operator.
type Monoid interface {
	constraints.Integer | constraints.Float | constraints.Complex | ~string
}

// Monad is a backward-compatible type alias for Monoid.
type Monad = Monoid

// Map returns a new lazy iterator transforming each element of itr with f.
func Map[T, U any](itr iter.Seq[T], f func(T) U) iter.Seq[U] {
	return func(yield func(U) bool) {
		for t := range itr {
			if !yield(f(t)) {
				return
			}
		}
	}
}

// Flatten flattens a nested sequence of iterators into a single lazy iterator.
func Flatten[T any](itrs iter.Seq[iter.Seq[T]]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for itr := range itrs {
			for t := range itr {
				if !yield(t) {
					return
				}
			}
		}
	}
}

// FlatMap maps each element of iter to a sub-iterator using f and flattens the result.
func FlatMap[T, U any](iter iter.Seq[T], f func(T) iter.Seq[U]) iter.Seq[U] {
	return Flatten(Map(iter, f))
}

// Filter returns a lazy iterator yielding only elements of itr that satisfy predicate p.
func Filter[T any](itr iter.Seq[T], p func(T) bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		for t := range itr {
			if p(t) {
				if !yield(t) {
					return
				}
			}
		}
	}
}

// FoldLeft eagerly accumulates values of itr from left to right starting with u using f.
func FoldLeft[T, U any](itr iter.Seq[T], f func(U, T) U, u U) U {
	result := u
	for t := range itr {
		result = f(result, t)
	}
	return result
}

// FoldRight eagerly accumulates values of itr from right to left starting with u using f.
func FoldRight[T, U any](itr iter.Seq[T], f func(T, U) U, u U) U {
	s := slices.Collect(itr)
	return slice.FoldRight(s, f, u)
}

// Reduce reduces itr using associative operation f and initial value t.
func Reduce[T any](itr iter.Seq[T], f func(T, T) T, t T) T {
	return FoldLeft(itr, f, t)
}

// Sum computes the sum of all elements in itr using their zero value as identity.
func Sum[M Monoid](itr iter.Seq[M]) M {
	var identity M
	return Reduce(itr, func(a, b M) M { return a + b }, identity)
}

// JoinErrs combines multiple errors yielded by itr into a single joined error.
func JoinErrs(itr iter.Seq[error]) error {
	var errs []error
	for err := range itr {
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

// Join concatenates elements of a string-like iterator separated by sep.
func Join[T ~string](itr iter.Seq[T], sep T) T {
	first := true
	var result T
	for t := range itr {
		if first {
			first = false
		} else {
			result += sep
		}
		result += t
	}
	return result
}

// Zip pairs elements from itr1 and itr2 into an iter.Seq2 until either iterator is exhausted.
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

// UnZip splits an iter.Seq2 into two separate iterators.
func UnZip[T, U any](itr iter.Seq2[T, U]) (iter.Seq[T], iter.Seq[U]) {
	ts, us := slice.Collect(itr)
	return slices.Values(ts), slices.Values(us)
}

// Sorted collects elements of itr, sorts them, and returns a new iterator.
func Sorted[T cmp.Ordered](itr iter.Seq[T]) iter.Seq[T] {
	return slices.Values(slices.Sorted(itr))
}

// Distinct returns a lazy iterator yielding only unique elements from itr in order of first appearance.
func Distinct[T comparable](itr iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		set := make(map[T]struct{})
		for t := range itr {
			if _, ok := set[t]; !ok {
				set[t] = struct{}{}
				if !yield(t) {
					return
				}
			}
		}
	}
}

// Generate returns an infinite lazy iterator that repeatedly calls supplier.
func Generate[T any](supplier func() T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for yield(supplier()) {
		}
	}
}

// Iterate returns a lazy iterator generating elements starting from seed, continuing
// while hasNext returns true, and advancing via next.
func Iterate[T any](seed T, hasNext func(T) bool, next func(T) T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for cur := seed; hasNext(cur); cur = next(cur) {
			if !yield(cur) {
				return
			}
		}
	}
}

// Range returns an iterator generating integer values from startInclusive to endExclusive - 1.
func Range[T constraints.Integer](startInclusive, endExclusive T) iter.Seq[T] {
	return Iterate(startInclusive, func(t T) bool { return t < endExclusive }, func(t T) T { t++; return t })
}

// RangeClosed returns an iterator generating integer values from startInclusive to endInclusive.
func RangeClosed[T constraints.Integer](startInclusive, endInclusive T) iter.Seq[T] {
	return Iterate(startInclusive, func(t T) bool { return t <= endInclusive }, func(t T) T { t++; return t })
}

// Limit returns a lazy iterator yielding at most max elements from itr.
func Limit[T any](itr iter.Seq[T], max int64) iter.Seq[T] {
	return func(yield func(T) bool) {
		var count int64
		for t := range itr {
			if count >= max || !yield(t) {
				return
			}
			count++
		}
	}
}

// Skip returns a lazy iterator skipping the first n elements of itr.
func Skip[T any](itr iter.Seq[T], n int64) iter.Seq[T] {
	return func(yield func(T) bool) {
		var count int64
		for t := range itr {
			if count >= n {
				if !yield(t) {
					return
				}
			}
			count++
		}
	}
}

// AllMatch reports whether all elements of itr satisfy predicate p.
// It short-circuits on the first element for which p returns false.
func AllMatch[T any](itr iter.Seq[T], p func(T) bool) bool {
	for t := range itr {
		if !p(t) {
			return false
		}
	}
	return true
}

// AnyMatch reports whether any element of itr satisfies predicate p.
// It short-circuits on the first element for which p returns true.
func AnyMatch[T any](itr iter.Seq[T], p func(T) bool) bool {
	for t := range itr {
		if p(t) {
			return true
		}
	}
	return false
}

// Count eagerly counts the number of elements in itr.
func Count[T any](itr iter.Seq[T]) int64 {
	return Sum(Map(itr, func(t T) int64 { return 1 }))
}

// Concat concatenates multiple iterators in sequence into a single lazy iterator.
func Concat[T any](itrs ...iter.Seq[T]) iter.Seq[T] {
	return Flatten(slices.Values(itrs))
}

// Peek returns a lazy iterator that executes consumer on each element before yielding it.
func Peek[T any](itr iter.Seq[T], consumer func(T)) iter.Seq[T] {
	return func(yield func(T) bool) {
		for t := range itr {
			consumer(t)
			if !yield(t) {
				return
			}
		}
	}
}

// Of returns a lazy iterator over the provided variadic elements.
func Of[T any](ts ...T) iter.Seq[T] {
	return slices.Values(ts)
}

// Partition returns a lazy sequence of chunk iterators, each yielding up to size elements.
func Partition[T any](itr iter.Seq[T], size int) iter.Seq[iter.Seq[T]] {
	return func(yield func(iter.Seq[T]) bool) {
		if size <= 0 {
			return
		}
		var chunk []T
		for t := range itr {
			chunk = append(chunk, t)
			if len(chunk) == size {
				if !yield(slices.Values(chunk)) {
					return
				}
				chunk = nil
			}
		}
		if len(chunk) > 0 {
			yield(slices.Values(chunk))
		}
	}
}
