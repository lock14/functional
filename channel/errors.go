package channel

func MapWithErr[T, U any](channel chan T, f func(T) (U, error)) (chan U, chan error) {
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

func FlatMapWithErr[T, U any](channel chan T, f func(T) (chan U, error)) (chan U, chan error) {
	channels, errs := MapWithErr(channel, f)
	return Flatten(channels), errs
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
