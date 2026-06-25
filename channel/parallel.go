package channel

import (
	"context"
	"sync"
)

func ParallelMap[T, U any](ctx context.Context, workers int, channel chan T, f func(T) U) chan U {
	mapped := make(chan U)
	go func() {
		defer close(mapped)
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
						u := f(t)
						select {
						case <-ctx.Done():
							return
						case mapped <- u:
						}
					}
				}
			}()
		}
		wg.Wait()
	}()
	return mapped
}

func ParallelFlatten[T any](ctx context.Context, workers int, channel chan chan T) chan T {
	flat := make(chan T)
	go func() {
		defer close(flat)
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
					case c, ok := <-channel:
						if !ok {
							return
						}
						for {
							select {
							case <-ctx.Done():
								return
							case t, ok2 := <-c:
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
		}
		wg.Wait()
	}()
	return flat
}

func ParallelFlatMap[T, U any](ctx context.Context, workers int, channel chan T, f func(T) chan U) chan U {
	return ParallelFlatten(ctx, workers, ParallelMap(ctx, workers, channel, f))
}

func ParallelFilter[T any](ctx context.Context, workers int, channel chan T, p func(T) bool) chan T {
	filtered := make(chan T)
	go func() {
		defer close(filtered)
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
		}
		wg.Wait()
	}()
	return filtered
}
