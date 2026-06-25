package channel

import (
	"context"
	"errors"
	"iter"
	"sort"
	"sync"

	"golang.org/x/exp/constraints"
)

// Monad represents any type that can use the `+` operator and whose zero
// value is the identity element the `+` operator
type Monad interface {
	constraints.Integer | constraints.Float | constraints.Complex | ~string
}

func Map[T, U any](ctx context.Context, channel chan T, f func(T) U) chan U {
	mapped := make(chan U)
	go func() {
		defer close(mapped)
		for {
			select {
			case <-ctx.Done():
				return
			case t, ok := <-channel:
				if !ok {
					return
				}
				select {
				case <-ctx.Done():
					return
				case mapped <- f(t):
				}
			}
		}
	}()
	return mapped
}

func Flatten[T any](ctx context.Context, channels chan chan T) chan T {
	flat := make(chan T)
	go func() {
		defer close(flat)
		for {
			select {
			case <-ctx.Done():
				return
			case channel, ok := <-channels:
				if !ok {
					return
				}
				for {
					select {
					case <-ctx.Done():
						return
					case t, ok2 := <-channel:
						if !ok2 {
							goto nextChannel
						}
						select {
						case <-ctx.Done():
							return
						case flat <- t:
						}
					}
				}
			nextChannel:
			}
		}
	}()
	return flat
}

func FlatMap[T, U any](ctx context.Context, channel chan T, f func(T) chan U) chan U {
	return Flatten(ctx, Map(ctx, channel, f))
}

func Filter[T any](ctx context.Context, channel chan T, p func(T) bool) chan T {
	filtered := make(chan T)
	go func() {
		defer close(filtered)
		for {
			select {
			case <-ctx.Done():
				return
			case t, ok := <-channel:
				if !ok {
					return
				}
				if p(t) {
					select {
					case <-ctx.Done():
						return
					case filtered <- t:
					}
				}
			}
		}
	}()
	return filtered
}

func FoldLeft[T, U any](ctx context.Context, channel chan T, f func(u U, t T) U, u U) U {
	result := u
	for {
		select {
		case <-ctx.Done():
			return result
		case t, ok := <-channel:
			if !ok {
				return result
			}
			result = f(result, t)
		}
	}
}

func FoldRight[T, U any](ctx context.Context, channel chan T, f func(t T, u U) U, u U) U {
	s := ToSlice(ctx, channel)
	result := u
	for i := len(s) - 1; i >= 0; i-- {
		result = f(s[i], result)
	}
	return result
}

func Reduce[T any](ctx context.Context, channel chan T, op func(t1, t2 T) T, initial T) T {
	return FoldLeft(ctx, channel, op, initial)
}

func Sum[M Monad](ctx context.Context, elements chan M) M {
	var identity M
	return Reduce(ctx, elements, func(a, b M) M { return a + b }, identity)
}

func JoinErrs(ctx context.Context, errs chan error) error {
	return Reduce(ctx, errs, func(e1, e2 error) error { return errors.Join(e1, e2) }, nil)
}

func Join[T ~string](ctx context.Context, strings chan T, sep T) T {
	select {
	case <-ctx.Done():
		var zero T
		return zero
	case first, ok := <-strings:
		if !ok {
			var zero T
			return zero
		}
		return first + Reduce(ctx, strings, func(a, b T) T { return a + sep + b }, "")
	}
}

type Pair[T1, T2 any] struct {
	Fst T1
	Snd T2
}

func Zip[T, U any](ctx context.Context, chan1 chan T, chan2 chan U) chan Pair[T, U] {
	zipped := make(chan Pair[T, U])
	go func() {
		defer close(zipped)
		for {
			select {
			case <-ctx.Done():
				return
			case t, ok1 := <-chan1:
				if !ok1 {
					return
				}
				select {
				case <-ctx.Done():
					return
				case u, ok2 := <-chan2:
					if !ok2 {
						return
					}
					select {
					case <-ctx.Done():
						return
					case zipped <- Pair[T, U]{Fst: t, Snd: u}:
					}
				}
			}
		}
	}()
	return zipped
}

func UnZip[T, U any](ctx context.Context, channel chan Pair[T, U]) (chan T, chan U) {
	ts := make(chan T)
	us := make(chan U)
	go func() {
		clones := Clone(ctx, channel, 2)
		go func() {
			defer close(ts)
			for {
				select {
				case <-ctx.Done():
					return
				case p, ok := <-clones[0]:
					if !ok {
						return
					}
					select {
					case <-ctx.Done():
						return
					case ts <- p.Fst:
					}
				}
			}
		}()
		go func() {
			defer close(us)
			for {
				select {
				case <-ctx.Done():
					return
				case p, ok := <-clones[1]:
					if !ok {
						return
					}
					select {
					case <-ctx.Done():
						return
					case us <- p.Snd:
					}
				}
			}
		}()
	}()
	return ts, us
}

func Sorted[T constraints.Ordered](ctx context.Context, channel chan T) chan T {
	ordered := make(chan T)
	go func() {
		defer close(ordered)
		var buf []T
		for {
			select {
			case <-ctx.Done():
				return
			case t, ok := <-channel:
				if !ok {
					goto sortPhase
				}
				buf = append(buf, t)
			}
		}
	sortPhase:
		sort.Slice(buf, func(i, j int) bool {
			return buf[i] < buf[j]
		})
		for _, t := range buf {
			select {
			case <-ctx.Done():
				return
			case ordered <- t:
			}
		}
	}()
	return ordered
}

func Distinct[T comparable](ctx context.Context, channel chan T) chan T {
	distinct := make(chan T)
	go func() {
		defer close(distinct)
		set := make(map[T]struct{})
		for {
			select {
			case <-ctx.Done():
				return
			case t, ok := <-channel:
				if !ok {
					return
				}
				if _, ok := set[t]; !ok {
					set[t] = struct{}{}
					select {
					case <-ctx.Done():
						return
					case distinct <- t:
					}
				}
			}
		}
	}()
	return distinct
}

func FromSlice[T any](slice []T) chan T {
	channel := make(chan T, len(slice))
	for _, t := range slice {
		channel <- t
	}
	close(channel)
	return channel
}

func ToSlice[T any](ctx context.Context, channel chan T) []T {
	var slice []T
	for {
		select {
		case <-ctx.Done():
			return slice
		case t, ok := <-channel:
			if !ok {
				return slice
			}
			slice = append(slice, t)
		}
	}
}

func Generate[T any](ctx context.Context, supplier func() T) chan T {
	c := make(chan T)
	go func() {
		defer close(c)
		for {
			select {
			case <-ctx.Done():
				return
			case c <- supplier():
			}
		}
	}()
	return c
}

func Iterate[T any](ctx context.Context, seed T, hasNext func(T) bool, next func(T) T) chan T {
	c := make(chan T)
	go func() {
		defer close(c)
		for cur := seed; hasNext(cur); cur = next(cur) {
			select {
			case <-ctx.Done():
				return
			case c <- cur:
			}
		}
	}()
	return c
}

func Range[T constraints.Integer](ctx context.Context, startInclusive, endExclusive T) chan T {
	return Iterate(ctx, startInclusive, func(t T) bool { return t < endExclusive }, func(t T) T { t++; return t })
}

func RangeClosed[T constraints.Integer](ctx context.Context, startInclusive, endInclusive T) chan T {
	return Iterate(ctx, startInclusive, func(t T) bool { return t <= endInclusive }, func(t T) T { t++; return t })
}

func Limit[T any](ctx context.Context, channel chan T, max int64) chan T {
	c := make(chan T)
	go func() {
		defer close(c)
		count := int64(0)
		for {
			if count >= max {
				return
			}
			select {
			case <-ctx.Done():
				return
			case t, ok := <-channel:
				if !ok {
					return
				}
				select {
				case <-ctx.Done():
					return
				case c <- t:
					count++
				}
			}
		}
	}()
	return c
}

func Skip[T any](ctx context.Context, channel chan T, n int64) chan T {
	c := make(chan T)
	go func() {
		defer close(c)
		count := int64(0)
		for {
			select {
			case <-ctx.Done():
				return
			case t, ok := <-channel:
				if !ok {
					return
				}
				if count >= n {
					select {
					case <-ctx.Done():
						return
					case c <- t:
					}
				}
				count++
			}
		}
	}()
	return c
}

func AllMatch[T any](ctx context.Context, channel chan T, p func(T) bool) bool {
	return Reduce(ctx, Map(ctx, channel, p), func(t1, t2 bool) bool { return t1 && t2 }, true)
}

func AnyMatch[T any](ctx context.Context, channel chan T, p func(T) bool) bool {
	return Reduce(ctx, Map(ctx, channel, p), func(t1, t2 bool) bool { return t1 || t2 }, false)
}

func TakeWhile[T any](ctx context.Context, channel chan T, p func(T) bool) chan T {
	c := make(chan T)
	go func() {
		defer close(c)
		for {
			select {
			case <-ctx.Done():
				return
			case t, ok := <-channel:
				if !ok {
					return
				}
				if p(t) {
					select {
					case <-ctx.Done():
						return
					case c <- t:
					}
				} else {
					return
				}
			}
		}
	}()
	return c
}

func Count[T any](ctx context.Context, channel chan T) int64 {
	return Sum(ctx, Map(ctx, channel, func(t T) int64 { return 1 }))
}

func Concat[T any](ctx context.Context, chan1, chan2 chan T) chan T {
	c := make(chan T)
	go func() {
		defer close(c)
		for {
			select {
			case <-ctx.Done():
				return
			case t, ok := <-chan1:
				if !ok {
					goto secondChannel
				}
				select {
				case <-ctx.Done():
					return
				case c <- t:
				}
			}
		}
	secondChannel:
		for {
			select {
			case <-ctx.Done():
				return
			case t, ok := <-chan2:
				if !ok {
					return
				}
				select {
				case <-ctx.Done():
					return
				case c <- t:
				}
			}
		}
	}()
	return c
}

func Peek[T any](ctx context.Context, channel chan T, consumer func(T)) chan T {
	c := make(chan T)
	go func() {
		defer close(c)
		for {
			select {
			case <-ctx.Done():
				return
			case t, ok := <-channel:
				if !ok {
					return
				}
				consumer(t)
				select {
				case <-ctx.Done():
					return
				case c <- t:
				}
			}
		}
	}()
	return c
}

func ForEach[T any](ctx context.Context, channel chan T, consumer func(T)) {
	for {
		select {
		case <-ctx.Done():
			return
		case t, ok := <-channel:
			if !ok {
				return
			}
			consumer(t)
		}
	}
}

func Of[T any](ts ...T) chan T {
	return FromSlice(ts)
}

func Partition[T any](ctx context.Context, channel chan T, size int) chan []T {
	partitioned := make(chan []T)
	go func() {
		defer close(partitioned)
		if size <= 0 {
			return
		}
		var chunk []T
		for {
			select {
			case <-ctx.Done():
				return
			case t, ok := <-channel:
				if !ok {
					if len(chunk) > 0 {
						select {
						case <-ctx.Done():
						case partitioned <- chunk:
						}
					}
					return
				}
				chunk = append(chunk, t)
				if len(chunk) == size {
					select {
					case <-ctx.Done():
						return
					case partitioned <- chunk:
						chunk = nil
					}
				}
			}
		}
	}()
	return partitioned
}

func Clone[T any](ctx context.Context, channel chan T, numClones int) []chan T {
	clones := make([]chan T, numClones)
	for i := 0; i < numClones; i++ {
		clones[i] = make(chan T)
	}
	go func() {
		waitGroups := make([]*sync.WaitGroup, len(clones))
		for i := 0; i < numClones; i++ {
			waitGroups[i] = &sync.WaitGroup{}
		}
		orders := make([]chan uint64, len(clones))
		for i := 0; i < numClones; i++ {
			orders[i] = make(chan uint64, 1)
			orders[i] <- 0
		}
		count := uint64(0)
		for {
			select {
			case <-ctx.Done():
				goto cleanup
			case t, ok := <-channel:
				if !ok {
					goto cleanup
				}
				for i := 0; i < numClones; i++ {
					waitGroups[i].Add(1)
					go func(order uint64, i int) {
						defer waitGroups[i].Done()
						for {
							select {
							case <-ctx.Done():
								return
							case o := <-orders[i]:
								if o == order {
									select {
									case <-ctx.Done():
									case clones[i] <- t:
										orders[i] <- order + 1
									}
									return
								}
								orders[i] <- o
							}
						}
					}(count, i)
				}
				count++
			}
		}
	cleanup:
		for i := 0; i < numClones; i++ {
			go func(i int) {
				waitGroups[i].Wait()
				close(clones[i])
			}(i)
		}
	}()
	return clones
}

func Stream[T any](ctx context.Context, seq iter.Seq[T]) chan T {
	c := make(chan T)
	go func() {
		defer close(c)
		for t := range seq {
			select {
			case <-ctx.Done():
				return
			case c <- t:
			}
		}
	}()
	return c
}
