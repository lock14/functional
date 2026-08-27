package channel

import (
	"context"
)

// MapWithErr maps elements using a fallible function f, splitting values and errors into separate channels.
func MapWithErr[T, U any](ctx context.Context, channel chan T, f func(T) (U, error)) (chan U, chan error) {
	mapped := make(chan U)
	errs := make(chan error)
	go func() {
		defer close(mapped)
		defer close(errs)
		for {
			select {
			case <-ctx.Done():
				return
			case t, ok := <-channel:
				if !ok {
					return
				}
				u, err := f(t)
				if err != nil {
					select {
					case <-ctx.Done():
						return
					case errs <- err:
					}
				} else {
					select {
					case <-ctx.Done():
						return
					case mapped <- u:
					}
				}
			}
		}
	}()
	return mapped, errs
}

// FlatMapWithErr maps elements using a fallible function f returning sub-channels and flattens them, splitting errors.
func FlatMapWithErr[T, U any](ctx context.Context, channel chan T, f func(T) (chan U, error)) (chan U, chan error) {
	channels, errs := MapWithErr(ctx, channel, f)
	return Flatten(ctx, channels), errs
}

// FilterWithErr filters elements using a fallible predicate p, splitting errors into a separate channel.
func FilterWithErr[T any](ctx context.Context, channel chan T, p func(T) (bool, error)) (chan T, chan error) {
	filtered := make(chan T)
	errs := make(chan error)
	go func() {
		defer close(filtered)
		defer close(errs)
		for {
			select {
			case <-ctx.Done():
				return
			case t, ok := <-channel:
				if !ok {
					return
				}
				ok2, err := p(t)
				if err != nil {
					select {
					case <-ctx.Done():
						return
					case errs <- err:
					}
				} else if ok2 {
					select {
					case <-ctx.Done():
						return
					case filtered <- t:
					}
				}
			}
		}
	}()
	return filtered, errs
}
