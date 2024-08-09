package channels

import "golang.org/x/exp/constraints"

// Monad represents any type that can use the `+` operator and whose zero
// value is the identity element the `+` operator
type Monad interface {
	constraints.Integer | constraints.Float | constraints.Complex | ~string
}

func Map[T any, U any](channel chan T, f func(T) U) chan U {
	mapped := make(chan U)
	go func() {
		for t := range channel {
			mapped <- f(t)
		}
		close(mapped)
	}()
	return mapped
}

func FoldLeft[T any, U any](channel chan T, f func(u U, t T) U, u U) U {
	result := u
	for t := range channel {
		result = f(result, t)
	}
	return result
}

func FoldRight[T any, U any](channel chan T, f func(t T, u U) U, u U) U {
	result := u
	for t := range channel {
		result = f(t, FoldRight[T, U](channel, f, u))
	}
	return result
}

func Reduce[T any](channel chan T, op func(t1, t2 T) T, initial T) T {
	return FoldLeft(channel, op, initial)
}

func Filter[T any](channel chan T, p func(T) bool) chan T {
	filtered := make(chan T)
	go func() {
		for t := range channel {
			if p(t) {
				filtered <- t
			}
		}
		close(filtered)
	}()
	return filtered
}

func Sum[M Monad](numbers chan M) M {
	var identity M
	return Reduce(numbers, func(a, b M) M { return a + b }, identity)
}

func Iterate[T constraints.Integer](start, end T) chan T {
	c := make(chan T)
	go func() {
		for i := start; i < end; i++ {
			c <- i
		}
		close(c)
	}()
	return c
}

func Range[T constraints.Integer](end T) chan T {
	c := make(chan T)
	go func() {
		var start T
		for i := start; i < end; i++ {
			c <- i
		}
		close(c)
	}()
	return c
}
