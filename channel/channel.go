// Package channel provides concurrent, streaming functional operations over Go channels.
package channel

import (
	"context"
	"errors"
	"iter"
	"sort"

	"golang.org/x/exp/constraints"
)

// Monoid represents any type that supports the `+` operator and whose zero
// value is the identity element of the `+` operator.
type Monoid interface {
	constraints.Integer | constraints.Float | constraints.Complex | ~string
}

// Monad is a backward-compatible type alias for Monoid.
type Monad = Monoid

// Map returns a channel streaming elements transformed by function f.
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

// Flatten merges a channel of channels sequentially into a single channel.
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

// FlatMap maps elements to channels using f and flattens the resulting channels.
func FlatMap[T, U any](ctx context.Context, channel chan T, f func(T) chan U) chan U {
	return Flatten(ctx, Map(ctx, channel, f))
}

// Filter streams only elements satisfying predicate p.
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

// FoldLeft accumulates elements from channel starting from u using f until channel is closed or ctx is cancelled.
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

// FoldRight buffers the channel into a slice and folds right-to-left.
func FoldRight[T, U any](ctx context.Context, channel chan T, f func(t T, u U) U, u U) U {
	s := ToSlice(ctx, channel)
	result := u
	for i := len(s) - 1; i >= 0; i-- {
		result = f(s[i], result)
	}
	return result
}

// Reduce reduces elements of channel using op starting with initial.
func Reduce[T any](ctx context.Context, channel chan T, op func(t1, t2 T) T, initial T) T {
	return FoldLeft(ctx, channel, op, initial)
}

// Sum calculates the sum of all numeric or string elements received from channel.
func Sum[M Monoid](ctx context.Context, elements chan M) M {
	var identity M
	return Reduce(ctx, elements, func(a, b M) M { return a + b }, identity)
}

// JoinErrs accumulates non-nil errors from errs into a single combined error using errors.Join.
func JoinErrs(ctx context.Context, errs chan error) error {
	var errList []error
	for {
		select {
		case <-ctx.Done():
			return errors.Join(errList...)
		case err, ok := <-errs:
			if !ok {
				return errors.Join(errList...)
			}
			if err != nil {
				errList = append(errList, err)
			}
		}
	}
}

// Join concatenates string-like channel elements with sep delimiter.
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

// Pair represents a 2-tuple of values streamed concurrently.
type Pair[T1, T2 any] struct {
	Fst T1
	Snd T2
}

// Zip combines elements from chan1 and chan2 into Pairs until either channel closes.
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

// UnZip splits a channel of Pairs into two separate channels.
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

// Sorted buffers all elements of channel, sorts them, and yields them in ascending order.
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

// Distinct yields only distinct elements received from channel.
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

// FromSlice converts a slice into a buffered, closed channel containing its elements.
func FromSlice[T any](slice []T) chan T {
	channel := make(chan T, len(slice))
	for _, t := range slice {
		channel <- t
	}
	close(channel)
	return channel
}

// ToSlice drains all elements of channel into a slice until closed or ctx is cancelled.
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

// Generate repeatedly calls supplier and streams values until ctx is cancelled.
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

// Iterate generates elements starting from seed, continuing while hasNext is true, and advancing via next.
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

// Range generates integers from startInclusive to endExclusive - 1.
func Range[T constraints.Integer](ctx context.Context, startInclusive, endExclusive T) chan T {
	return Iterate(ctx, startInclusive, func(t T) bool { return t < endExclusive }, func(t T) T { t++; return t })
}

// RangeClosed generates integers from startInclusive to endInclusive.
func RangeClosed[T constraints.Integer](ctx context.Context, startInclusive, endInclusive T) chan T {
	return Iterate(ctx, startInclusive, func(t T) bool { return t <= endInclusive }, func(t T) T { t++; return t })
}

// Limit yields at most max elements from channel.
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

// Skip ignores the first n elements from channel and yields the rest.
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

// AllMatch reports whether all elements in channel satisfy predicate p.
// It short-circuits on the first element for which p returns false.
func AllMatch[T any](ctx context.Context, channel chan T, p func(T) bool) bool {
	for {
		select {
		case <-ctx.Done():
			return false
		case t, ok := <-channel:
			if !ok {
				return true
			}
			if !p(t) {
				return false
			}
		}
	}
}

// AnyMatch reports whether any element in channel satisfies predicate p.
// It short-circuits on the first element for which p returns true.
func AnyMatch[T any](ctx context.Context, channel chan T, p func(T) bool) bool {
	for {
		select {
		case <-ctx.Done():
			return false
		case t, ok := <-channel:
			if !ok {
				return false
			}
			if p(t) {
				return true
			}
		}
	}
}

// TakeWhile yields elements as long as predicate p is satisfied, then closes.
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

// Count drains channel and counts the number of elements received.
func Count[T any](ctx context.Context, channel chan T) int64 {
	var count int64
	for {
		select {
		case <-ctx.Done():
			return count
		case _, ok := <-channel:
			if !ok {
				return count
			}
			count++
		}
	}
}

// Concat streams all elements of chan1 followed by all elements of chan2.
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

// Peek executes consumer on each element received before forwarding it down the returned channel.
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

// ForEach consumes each element from channel and passes it to consumer until closed or ctx is cancelled.
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

// Of returns a channel containing the provided variadic elements.
func Of[T any](ts ...T) chan T {
	return FromSlice(ts)
}

// Partition groups incoming items into slices of size elements and yields them.
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

// Clone duplicates elements received from channel to numClones separate output channels.
// Each output channel is buffered independently so that consumers do not deadlock if read
// at different rates. When ctx is cancelled or channel closes, all cloned channels are closed.
func Clone[T any](ctx context.Context, channel chan T, numClones int) []chan T {
	clones := make([]chan T, numClones)
	for i := 0; i < numClones; i++ {
		clones[i] = make(chan T)
	}
	if numClones <= 0 {
		return clones
	}

	inChans := make([]chan T, numClones)
	for i := 0; i < numClones; i++ {
		inChans[i] = make(chan T)
		go func(in chan T, out chan T) {
			defer close(out)
			var queue []T
			inOpen := true
			for {
				if len(queue) == 0 {
					if !inOpen {
						return
					}
					select {
					case <-ctx.Done():
						return
					case t, ok := <-in:
						if !ok {
							return
						}
						queue = append(queue, t)
					}
				} else {
					var nextIn chan T
					if inOpen {
						nextIn = in
					}
					select {
					case <-ctx.Done():
						return
					case t, ok := <-nextIn:
						if !ok {
							inOpen = false
						} else {
							queue = append(queue, t)
						}
					case out <- queue[0]:
						queue = queue[1:]
						if len(queue) == 0 {
							queue = nil
						}
					}
				}
			}
		}(inChans[i], clones[i])
	}

	go func() {
		defer func() {
			for i := 0; i < numClones; i++ {
				close(inChans[i])
			}
		}()
		for {
			select {
			case <-ctx.Done():
				return
			case t, ok := <-channel:
				if !ok {
					return
				}
				for i := 0; i < numClones; i++ {
					select {
					case <-ctx.Done():
						return
					case inChans[i] <- t:
					}
				}
			}
		}
	}()

	return clones
}

// Stream transforms a Go 1.23+ iter.Seq into a concurrent channel.
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
