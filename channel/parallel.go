package channel

import (
	"runtime"
	"sync"
)

func ParallelMap[T, U any](channel chan T, f func(T) U) chan U {
	mapped := make(chan U)
	go func() {
		concurrency := runtime.NumCPU()
		waitGroup := sync.WaitGroup{}
		for i := 0; i < concurrency; i++ {
			// spawn worker
			waitGroup.Add(1)
			go func() {
				defer waitGroup.Done()
				for t := range channel {
					mapped <- f(t)
				}
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
		concurrency := runtime.NumCPU()
		waitGroup := sync.WaitGroup{}
		for i := 0; i < concurrency; i++ {
			// spawn worker
			waitGroup.Add(1)
			go func() {
				defer waitGroup.Done()
				for c := range channel {
					for t := range c {
						flat <- t
					}
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
		concurrency := runtime.NumCPU()
		waitGroup := sync.WaitGroup{}
		for i := 0; i < concurrency; i++ {
			// spawn worker
			waitGroup.Add(1)
			go func() {
				defer waitGroup.Done()
				for t := range channel {
					if p(t) {
						filtered <- t
					}
				}
			}()
		}
		waitGroup.Wait()
		close(filtered)
	}()
	return filtered
}
