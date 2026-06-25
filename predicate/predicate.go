package predicate

import "reflect"

func IsNil[T any](t T) bool {
	v := reflect.ValueOf(t)
	if !v.IsValid() {
		return true
	}
	switch v.Type().Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Pointer,
		reflect.UnsafePointer, reflect.Interface, reflect.Slice:
		return v.IsNil()
	default:
		return false
	}
}

func NotNil[T any](t T) bool {
	return !IsNil(t)
}

func Not[T any](p func(T) bool) func(T) bool {
	return func(t T) bool {
		return !p(t)
	}
}

func True[T any](t T) bool {
	return true
}

func False[T any](t T) bool {
	return false
}
