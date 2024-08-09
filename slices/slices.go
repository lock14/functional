package slices

import "golang.org/x/exp/constraints"

// Monad represents any type
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
