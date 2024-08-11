package channel

import "sync"

func ParallelMap[T, U any](channel chan T, f func(T) U) chan U {
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

func ParallelFlatten[T any](channel chan chan T) chan T {
	flat := make(chan T)
	go func() {
		waitGroup := sync.WaitGroup{}
		for c := range channel {
			waitGroup.Add(1)
			go func() {
				defer waitGroup.Done()
				for t := range c {
					flat <- t
				}
			}()
		}
		waitGroup.Wait()
		close(flat)
	}()
	return flat
}

func ParallelFlatMap[T, U any](channel chan T, f func(T) chan U) chan U {
	return ParallelFlatten(ParallelMap(channel, f))
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
