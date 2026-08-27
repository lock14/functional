// Package predicate provides generic, composable predicate functions and nil-safety helpers.
package predicate

import "reflect"

// IsNil reports whether t is nil or an invalid reflect value.
// For non-nillable types (e.g. primitive ints or structs), IsNil returns false.
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

// NotNil reports whether t is not nil.
func NotNil[T any](t T) bool {
	return !IsNil(t)
}

// Not returns the logical negation of predicate p.
func Not[T any](p func(T) bool) func(T) bool {
	return func(t T) bool {
		return !p(t)
	}
}

// True is a predicate that always returns true regardless of input.
func True[T any](t T) bool {
	return true
}

// False is a predicate that always returns false regardless of input.
func False[T any](t T) bool {
	return false
}
