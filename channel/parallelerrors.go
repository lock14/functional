package channel

import "sync"

func ParallelMapWithErr[T, U any](channel chan T, f func(T) (U, error)) (chan U, chan error) {
	mapped := make(chan U)
	errs := make(chan error)
	go func() {
		waitGroup := sync.WaitGroup{}
		for t := range channel {
			waitGroup.Add(1)
			go func() {
				defer waitGroup.Done()
				u, err := f(t)
				if err != nil {
					errs <- err
				} else {
					mapped <- u
				}
			}()
		}
		waitGroup.Wait()
		close(mapped)
	}()
	return mapped, errs
}

func ParallelFlatMapWithErr[T, U any](channel chan T, f func(T) (chan U, error)) (chan U, chan error) {
	channels, errs := ParallelMapWithErr(channel, f)
	return ParallelFlatten(channels), errs
}

func ParallelFilterWithErr[T any](channel chan T, p func(T) (bool, error)) (chan T, chan error) {
	filtered := make(chan T)
	errs := make(chan error)
	go func() {
		waitGroup := sync.WaitGroup{}
		for t := range channel {
			waitGroup.Add(1)
			go func() {
				defer waitGroup.Done()
				ok, err := p(t)
				if err != nil {
					errs <- err
				} else if ok {
					filtered <- t
				}
			}()
		}
		waitGroup.Wait()
		close(filtered)
	}()
	return filtered, errs
}
