package predicates

import "reflect"

func IsNil[T any](t T) bool {
	switch reflect.ValueOf(t).Type().Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Pointer,
		reflect.UnsafePointer, reflect.Interface, reflect.Slice:
		return reflect.ValueOf(t).IsNil()
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
