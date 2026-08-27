package predicate

import (
	"errors"
	"testing"
)

func TestIsNil(t *testing.T) {
	t.Parallel()

	var nilPtr *int
	validInt := 42
	validPtr := &validInt
	var nilErr error
	validErr := errors.New("sample error")
	var nilSlice []int
	validSlice := []int{1, 2, 3}
	var nilMap map[string]int
	validMap := map[string]int{"a": 1}
	var nilChan chan int
	validChan := make(chan int)
	var nilFunc func()
	validFunc := func() {}

	cases := []struct {
		name  string
		input any
		want  bool
	}{
		{name: "nil_pointer", input: nilPtr, want: true},
		{name: "non_nil_pointer", input: validPtr, want: false},
		{name: "primitive_int", input: 5, want: false},
		{name: "primitive_string", input: "hello", want: false},
		{name: "primitive_struct", input: struct{}{}, want: false},
		{name: "nil_interface", input: nilErr, want: true},
		{name: "non_nil_interface", input: validErr, want: false},
		{name: "nil_slice", input: nilSlice, want: true},
		{name: "non_nil_slice", input: validSlice, want: false},
		{name: "nil_map", input: nilMap, want: true},
		{name: "non_nil_map", input: validMap, want: false},
		{name: "nil_chan", input: nilChan, want: true},
		{name: "non_nil_chan", input: validChan, want: false},
		{name: "nil_func", input: nilFunc, want: true},
		{name: "non_nil_func", input: validFunc, want: false},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := IsNil(tc.input)
			if got != tc.want {
				t.Errorf("IsNil(%v) = %v, want %v", tc.input, got, tc.want)
			}
		})
	}
}

func TestNotNil(t *testing.T) {
	t.Parallel()

	var nilPtr *int
	validInt := 5
	validPtr := &validInt

	cases := []struct {
		name  string
		input any
		want  bool
	}{
		{name: "nil_pointer", input: nilPtr, want: false},
		{name: "non_nil_pointer", input: validPtr, want: true},
		{name: "primitive_int", input: 10, want: true},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := NotNil(tc.input)
			if got != tc.want {
				t.Errorf("NotNil(%v) = %v, want %v", tc.input, got, tc.want)
			}
		})
	}
}

func TestNot(t *testing.T) {
	t.Parallel()

	isEven := func(n int) bool { return n%2 == 0 }
	isOdd := Not(isEven)

	cases := []struct {
		name  string
		input int
		want  bool
	}{
		{name: "even_number_negated", input: 2, want: false},
		{name: "odd_number_negated", input: 3, want: true},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := isOdd(tc.input)
			if got != tc.want {
				t.Errorf("Not(isEven)(%d) = %v, want %v", tc.input, got, tc.want)
			}
		})
	}
}

func TestTrue(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name  string
		input int
		want  bool
	}{
		{name: "zero", input: 0, want: true},
		{name: "positive", input: 42, want: true},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := True(tc.input)
			if got != tc.want {
				t.Errorf("True(%d) = %v, want %v", tc.input, got, tc.want)
			}
		})
	}
}

func TestFalse(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name  string
		input int
		want  bool
	}{
		{name: "zero", input: 0, want: false},
		{name: "positive", input: 42, want: false},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := False(tc.input)
			if got != tc.want {
				t.Errorf("False(%d) = %v, want %v", tc.input, got, tc.want)
			}
		})
	}
}
