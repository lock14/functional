package channel

import (
	"context"
	"sync"
)

// ParallelMapWithErr maps elements concurrently using workers goroutines with a fallible function f.
func ParallelMapWithErr[T, U any](ctx context.Context, workers int, channel chan T, f func(T) (U, error)) (chan U, chan error) {
	mapped := make(chan U)
	errs := make(chan error)
	go func() {
		defer close(mapped)
		defer close(errs)
		if workers <= 0 {
			return
		}
		var wg sync.WaitGroup
		for i := 0; i < workers; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
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
		}
		wg.Wait()
	}()
	return mapped, errs
}

// ParallelFlatMapWithErr maps elements to sub-channels concurrently with a fallible function f, flattens, and splits errors.
func ParallelFlatMapWithErr[T, U any](ctx context.Context, workers int, channel chan T, f func(T) (chan U, error)) (chan U, chan error) {
	channels, errs := ParallelMapWithErr(ctx, workers, channel, f)
	return ParallelFlatten(ctx, workers, channels), errs
}

// ParallelFilterWithErr filters elements concurrently using workers goroutines with a fallible predicate p.
func ParallelFilterWithErr[T any](ctx context.Context, workers int, channel chan T, p func(T) (bool, error)) (chan T, chan error) {
	filtered := make(chan T)
	errs := make(chan error)
	go func() {
		defer close(filtered)
		defer close(errs)
		if workers <= 0 {
			return
		}
		var wg sync.WaitGroup
		for i := 0; i < workers; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
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
		}
		wg.Wait()
	}()
	return filtered, errs
}
