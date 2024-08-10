package channels

import (
	"errors"
	"golang.org/x/exp/constraints"
	"sort"
	"sync"
)

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

func MapWithErr[T any, U any](channel chan T, f func(T) (U, error)) (chan U, chan error) {
	mapped := make(chan U)
	errs := make(chan error)
	go func() {
		for t := range channel {
			u, err := f(t)
			if err != nil {
				errs <- err
			} else {
				mapped <- u
			}
		}
		close(mapped)
		close(errs)
	}()
	return mapped, errs
}

func ParallelMap[T any, U any](channel chan T, f func(T) U) chan U {
	mapped := make(chan U)
	go func() {
		waitGroup := sync.WaitGroup{}
		for t := range channel {
			waitGroup.Add(1)
			go func() {
				defer waitGroup.Done()
				mapped <- f(t)
			}()
		}
		waitGroup.Wait()
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

func FilterWithErr[T any](channel chan T, p func(T) (bool, error)) (chan T, chan error) {
	filtered := make(chan T)
	errs := make(chan error)
	go func() {
		for t := range channel {
			ok, err := p(t)
			if err != nil {
				errs <- err
			} else if ok {
				filtered <- t
			}
		}
		close(filtered)
		close(errs)
	}()
	return filtered, errs
}

func ParallelFilter[T any](channel chan T, p func(T) bool) chan T {
	filtered := make(chan T)
	go func() {
		waitGroup := sync.WaitGroup{}
		for t := range channel {
			waitGroup.Add(1)
			go func() {
				defer waitGroup.Done()
				if p(t) {
					filtered <- t
				}
			}()
		}
		waitGroup.Wait()
		close(filtered)
	}()
	return filtered
}

func Sum[M Monad](numbers chan M) M {
	var identity M
	return Reduce(numbers, func(a, b M) M { return a + b }, identity)
}

func JoinErrs(errs chan error) error {
	return Reduce(errs, func(e1, e2 error) error {
		return errors.Join(e1, e2)
	}, nil)
}

func Sorted[T constraints.Ordered](channel chan T) chan T {
	ordered := make(chan T)
	go func() {
		var buf []T
		for t := range channel {
			buf = append(buf, t)
		}
		sort.Slice(buf, func(i, j int) bool {
			return buf[i] < buf[j]
		})
		for _, t := range buf {
			ordered <- t
		}
		close(ordered)
	}()
	return ordered
}

func Distinct[T comparable](channel chan T) chan T {
	distinct := make(chan T)
	go func() {
		set := make(map[T]struct{})
		for t := range channel {
			if _, ok := set[t]; !ok {
				set[t] = struct{}{}
				distinct <- t
			}
		}
		close(distinct)
	}()
	return distinct
}

func FromSlice[T any](slice []T) chan T {
	channel := make(chan T)
	go func() {
		for _, t := range slice {
			channel <- t
		}
		close(channel)
	}()
	return channel
}

func ToSlice[T any](channel chan T) []T {
	var slice []T
	for t := range channel {
		slice = append(slice, t)
	}
	return slice
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
