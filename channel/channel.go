package channel

import (
	"errors"
	"golang.org/x/exp/constraints"
	"sort"
	"sync/atomic"
)

// Monad represents any type that can use the `+` operator and whose zero
// value is the identity element the `+` operator
type Monad interface {
	constraints.Integer | constraints.Float | constraints.Complex | ~string
}

func Map[T, U any](channel chan T, f func(T) U) chan U {
	mapped := make(chan U)
	go func() {
		for t := range channel {
			mapped <- f(t)
		}
		close(mapped)
	}()
	return mapped
}

func Flatten[T any](channels chan chan T) chan T {
	flat := make(chan T)
	go func() {
		for channel := range channels {
			for t := range channel {
				flat <- t
			}
		}
		close(flat)
	}()
	return flat
}

func FlatMap[T, U any](channel chan T, f func(T) chan U) chan U {
	return Flatten(Map(channel, f))
}

func FoldLeft[T, U any](channel chan T, f func(u U, t T) U, u U) U {
	result := u
	for t := range channel {
		result = f(result, t)
	}
	return result
}

func FoldRight[T, U any](channel chan T, f func(t T, u U) U, u U) U {
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

func Sum[M Monad](elements chan M) M {
	var identity M
	return Reduce(elements, func(a, b M) M { return a + b }, identity)
}

func JoinErrs(errs chan error) error {
	return Reduce(errs, func(e1, e2 error) error { return errors.Join(e1, e2) }, nil)
}

func Join[T ~string](strings chan T, sep T) T {
	first, ok := <-strings
	if !ok {
		return first
	}
	return first + Reduce(strings, func(a, b T) T { return a + sep + b }, "")
}

type Pair[T1, T2 any] struct {
	fst T1
	snd T2
}

func Zip[T, U any](chan1 chan T, chan2 chan U) chan Pair[T, U] {
	zipped := make(chan Pair[T, U])
	go func() {
		t, ok1 := <-chan1
		u, ok2 := <-chan2
		for ok1 && ok2 {
			zipped <- Pair[T, U]{fst: t, snd: u}
			t, ok1 = <-chan1
			u, ok2 = <-chan2
		}
		close(zipped)
	}()
	return zipped
}

func UnZip[T, U any](channel chan Pair[T, U]) (chan T, chan U) {
	ts := make(chan T)
	us := make(chan U)
	go func() {
		for p := range channel {
			ts <- p.fst
			us <- p.snd
		}
		close(ts)
		close(us)
	}()
	return ts, us
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

func Generate[T any](supplier func() T) (chan T, func()) {
	c := make(chan T)
	keepGoing := atomic.Bool{}
	keepGoing.Store(true)
	closeFunc := func() {
		keepGoing.Store(false)
		// read from the channel to unblock the goroutine so it can read the bool
		// and close the channel.
		_ = <-c
	}
	go func() {
		for keepGoing.Load() {
			c <- supplier()
		}
		close(c)
	}()
	return c, closeFunc
}

func Iterate[T any](seed T, hasNext func(T) bool, next func(T) T) chan T {
	c := make(chan T)
	go func() {
		for cur := seed; hasNext(cur); cur = next(cur) {
			c <- cur
		}
		close(c)
	}()
	return c
}

func Range[T constraints.Integer](startInclusive, endExclusive T) chan T {
	return Iterate(startInclusive, func(t T) bool { return t < endExclusive }, func(t T) T { t++; return t })
}

func RangeClosed[T constraints.Integer](startInclusive, endInclusive T) chan T {
	return Iterate(startInclusive, func(t T) bool { return t <= endInclusive }, func(t T) T { t++; return t })
}

func Limit[T any](channel chan T, max int64) chan T {
	c := make(chan T)
	go func() {
		count := int64(0)
		for t := range channel {
			if count < max {
				c <- t
				count++
			} else {
				break
			}
		}
		close(c)
	}()
	return c
}

func Skip[T any](channel chan T, n int64) chan T {
	c := make(chan T)
	go func() {
		count := int64(0)
		for t := range channel {
			if count >= n {
				c <- t
			}
			count++
		}
		close(c)
	}()
	return c
}

func AllMatch[T any](channel chan T, p func(T) bool) bool {
	return Reduce(Map(channel, p), func(t1, t2 bool) bool { return t1 && t2 }, true)
}

func AnyMatch[T any](channel chan T, p func(T) bool) bool {
	return Reduce(Map(channel, p), func(t1, t2 bool) bool { return t1 || t2 }, false)
}

func TakeWhile[T any](chanel chan T, p func(T) bool) chan T {
	c := make(chan T)
	go func() {
		for t := range chanel {
			if p(t) {
				c <- t
			} else {
				break
			}
		}
		close(c)
	}()
	return c
}

func Count[T any](channel chan T) int64 {
	return Sum(Map(channel, func(t T) int64 { return 1 }))
}

func Concat[T any](chan1, chan2 chan T) chan T {
	c := make(chan T)
	go func() {
		for t := range chan1 {
			c <- t
		}
		for t := range chan2 {
			c <- t
		}
		close(c)
	}()
	return c
}

func Peek[T any](channel chan T, consumer func(T)) chan T {
	c := make(chan T)
	go func() {
		for t := range channel {
			consumer(t)
			c <- t
		}
		close(c)
	}()
	return c
}

func ForEach[T any](channel chan T, consumer func(T)) {
	for t := range channel {
		consumer(t)
	}
}

func Of[T any](ts ...T) chan T {
	c := make(chan T)
	go func() {
		for _, t := range ts {
			c <- t
		}
		close(c)
	}()
	return c
}

func Partition[T any](channel chan T, size int) chan []T {
	partitioned := make(chan []T)
	go func() {
		count := 0
		slice := make([]T, 0, size)
		for t := range channel {
			if count == size {
				partitioned <- slice
				slice = make([]T, 0, size)
				count = 0
			}
			slice = append(slice, t)
			count++
		}
		if count > 0 {
			partitioned <- slice
		}
		close(partitioned)
	}()
	return partitioned
}
