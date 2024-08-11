package slice

import (
	"errors"
	"golang.org/x/exp/constraints"
)

// Monad represents any type that can use the `+` operator and whose zero
// value is the identity element the `+` operator
type Monad interface {
	constraints.Integer | constraints.Float | constraints.Complex | ~string
}

func Map[T any, U any](slice []T, f func(T) U) []U {
	mapped := make([]U, 0, len(slice))
	for _, t := range slice {
		mapped = append(mapped, f(t))
	}
	return mapped
}

func Flatten[T any](slices [][]T) []T {
	var flattened []T
	for _, slice := range slices {
		for _, t := range slice {
			flattened = append(flattened, t)
		}
	}
	return flattened
}

func FlatMap[T, U any](slice []T, f func(T) []U) []U {
	return Flatten(Map(slice, f))
}

func FoldLeft[T any, U any](slice []T, f func(u U, t T) U, u U) U {
	result := u
	for _, t := range slice {
		result = f(result, t)
	}
	return result
}

func FoldRight[T any, U any](slice []T, f func(t T, u U) U, u U) U {
	result := u
	for i := len(slice) - 1; i >= 0; i-- {
		result = f(slice[i], result)
	}
	return result
}

func Reduce[T any](slice []T, op func(t1, t2 T) T, initial T) T {
	return FoldLeft(slice, op, initial)
}

func Filter[T any](slice []T, p func(T) bool) []T {
	var filtered []T
	for _, t := range slice {
		if p(t) {
			filtered = append(filtered, t)
		}
	}
	return filtered
}

func Sum[M Monad](numbers []M) M {
	var identity M
	return Reduce(numbers, func(a, b M) M { return a + b }, identity)
}

func JoinErrs(errs []error) error {
	return Reduce(errs, func(e1, e2 error) error { return errors.Join(e1, e2) }, nil)
}

func Join[T ~string](strings []T, sep T) T {
	if len(strings) == 0 {
		var zero T
		return zero
	}
	first := strings[0]
	strings = strings[1:]
	return first + Reduce(strings, func(a, b T) T { return a + sep + b }, "")
}

type Pair[T1, T2 any] struct {
	fst T1
	snd T2
}

func Zip[T, U any](slice1 []T, slice2 []U) []Pair[T, U] {
	len1 := len(slice1)
	len2 := len(slice2)
	minLen := len1
	if len2 < minLen {
		minLen = len2
	}
	zipped := make([]Pair[T, U], 0, minLen)
	for i := 0; i < minLen; i++ {
		zipped = append(zipped, Pair[T, U]{slice1[i], slice2[i]})
	}
	return zipped
}

func UnZip[T, U any](slice []Pair[T, U]) ([]T, []U) {
	ts := make([]T, 0, len(slice))
	us := make([]U, 0, len(slice))
	for _, p := range slice {
		ts = append(ts, p.fst)
		us = append(us, p.snd)
	}
	return ts, us
}

func Concat[T any](slice1, slice2 []T) []T {
	c := make([]T, 0, len(slice1)+len(slice2))
	for _, t := range slice1 {
		c = append(c, t)
	}
	for _, t := range slice2 {
		c = append(c, t)
	}
	return c
}

func Partition[T any](slice []T, size int) [][]T {
	partitioned := make([][]T, 0, len(slice)/size+1)
	count := 0
	partition := make([]T, 0, size)
	for _, t := range slice {
		if count == size {
			partitioned = append(partitioned, partition)
			partition = make([]T, 0, size)
			count = 0
		}
		partition = append(partition, t)
		count++
	}
	if count > 0 {
		partitioned = append(partitioned, partition)
	}
	return partitioned
}
