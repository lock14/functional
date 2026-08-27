// Package slice provides eager functional operations over Go slices.
package slice

import (
	"cmp"
	"errors"
	"iter"

	"golang.org/x/exp/constraints"
)

// Monoid represents any type that supports the `+` operator and whose zero
// value is the identity element of the `+` operator.
type Monoid interface {
	constraints.Integer | constraints.Float | constraints.Complex | ~string
}

// Monad is a backward-compatible type alias for Monoid.
type Monad = Monoid

// Map transforms each element of slice using function f, returning a new slice.
func Map[T any, U any](slice []T, f func(T) U) []U {
	mapped := make([]U, 0, len(slice))
	for _, t := range slice {
		mapped = append(mapped, f(t))
	}
	return mapped
}

// Flatten flattens a slice of slices into a single contiguous slice.
func Flatten[T any](slices [][]T) []T {
	var totalLen int
	for _, s := range slices {
		totalLen += len(s)
	}
	flattened := make([]T, 0, totalLen)
	for _, s := range slices {
		flattened = append(flattened, s...)
	}
	return flattened
}

// FlatMap applies f to each element of slice and flattens the resulting slices.
func FlatMap[T, U any](slice []T, f func(T) []U) []U {
	return Flatten(Map(slice, f))
}

// Filter returns a new slice containing only elements that satisfy predicate p.
func Filter[T any](slice []T, p func(T) bool) []T {
	var filtered []T
	for _, t := range slice {
		if p(t) {
			filtered = append(filtered, t)
		}
	}
	return filtered
}

// FoldLeft accumulates values starting from u from left to right using f.
func FoldLeft[T any, U any](slice []T, f func(u U, t T) U, u U) U {
	result := u
	for _, t := range slice {
		result = f(result, t)
	}
	return result
}

// FoldRight accumulates values starting from u from right to left using f.
func FoldRight[T any, U any](slice []T, f func(t T, u U) U, u U) U {
	result := u
	for i := len(slice) - 1; i >= 0; i-- {
		result = f(slice[i], result)
	}
	return result
}

// Reduce reduces slice using associative operation op and initial value initial.
func Reduce[T any](slice []T, op func(t1, t2 T) T, initial T) T {
	return FoldLeft(slice, op, initial)
}

// Sum computes the sum of all elements in numbers using their zero value as identity.
func Sum[M Monoid](numbers []M) M {
	var identity M
	return Reduce(numbers, func(a, b M) M { return a + b }, identity)
}

// JoinErrs combines multiple errors into a single error using errors.Join.
func JoinErrs(errs []error) error {
	return errors.Join(errs...)
}

// Join concatenates elements of a string-like slice separated by sep.
func Join[T ~string](strings []T, sep T) T {
	if len(strings) == 0 {
		var zero T
		return zero
	}
	first := strings[0]
	strings = strings[1:]
	return first + Reduce(strings, func(a, b T) T { return a + sep + b }, "")
}

// Pair represents a 2-tuple of values.
type Pair[T1, T2 any] struct {
	Fst T1
	Snd T2
}

// Zip combines two slices into a slice of Pairs up to the length of the shorter slice.
func Zip[T, U any](slice1 []T, slice2 []U) []Pair[T, U] {
	minLen := min(len(slice1), len(slice2))
	zipped := make([]Pair[T, U], 0, minLen)
	for i := 0; i < minLen; i++ {
		zipped = append(zipped, Pair[T, U]{Fst: slice1[i], Snd: slice2[i]})
	}
	return zipped
}

// UnZip splits a slice of Pairs into two separate slices.
func UnZip[T, U any](slice []Pair[T, U]) ([]T, []U) {
	ts := make([]T, 0, len(slice))
	us := make([]U, 0, len(slice))
	for _, p := range slice {
		ts = append(ts, p.Fst)
		us = append(us, p.Snd)
	}
	return ts, us
}

// Concat concatenates two slices into a new single slice.
func Concat[T any](slice1, slice2 []T) []T {
	c := make([]T, 0, len(slice1)+len(slice2))
	c = append(c, slice1...)
	c = append(c, slice2...)
	return c
}

// Partition divides slice into sub-slices of at most size elements each.
func Partition[T any](slice []T, size int) [][]T {
	if size <= 0 {
		return [][]T{}
	}
	partitioned := make([][]T, 0, (len(slice)+size-1)/size)
	for i := 0; i < len(slice); i += size {
		end := min(i+size, len(slice))
		partitioned = append(partitioned, slice[i:end])
	}
	return partitioned
}

// Collect consumes an iter.Seq2 and collects elements into two separate slices.
func Collect[T, U any](seq2 iter.Seq2[T, U]) ([]T, []U) {
	var ts []T
	var us []U
	for t, u := range seq2 {
		ts = append(ts, t)
		us = append(us, u)
	}
	return ts, us
}

// Ordered represents ordered types for constraints.
type Ordered = cmp.Ordered
