package channel

import (
	"github.com/google/go-cmp/cmp"
	"strconv"
	"testing"
)

func TestMap(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name        string
		input       []int
		mappingFunc func(int) string
		want        []string
	}{
		{
			name:  "map_empty",
			input: []int{},
			mappingFunc: func(i int) string {
				t.Error("mapping function was called when it should not have been")
				return ""
			},
			want: nil,
		},
		{
			name:        "map_one",
			input:       []int{1},
			mappingFunc: strconv.Itoa,
			want:        []string{"1"},
		},
		{
			name:        "map_many",
			input:       []int{1, 2, 3},
			mappingFunc: strconv.Itoa,
			want:        []string{"1", "2", "3"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			input := FromSlice(tc.input)
			mappedChan := Map(input, tc.mappingFunc)
			got := ToSlice(mappedChan)
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("unexpected result (-got, +want): %s", diff)
			}
			// check that both channels are closed now
			_, ok := <-input
			if ok {
				t.Error("expected input to be closed ")
			}
			_, ok = <-mappedChan
			if ok {
				t.Error("expected mappedChan to be closed ")
			}
		})
	}
}

func TestFlatMap(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name        string
		input       []int
		mappingFunc func(int) chan string
		want        []string
	}{
		{
			name:  "map_empty",
			input: []int{},
			mappingFunc: func(i int) chan string {
				t.Error("mapping function was called when it should not have been")
				return nil
			},
			want: nil,
		},
		{
			name:  "map_one",
			input: []int{1},
			mappingFunc: func(n int) chan string {
				c := make(chan string)
				go func() {
					for i := 0; i <= n; i++ {
						c <- strconv.Itoa(i)
					}
					close(c)
				}()
				return c
			},
			want: []string{"0", "1"},
		},
		{
			name:  "map_many",
			input: []int{1, 2, 3},
			mappingFunc: func(n int) chan string {
				c := make(chan string)
				go func() {
					for i := 0; i <= n; i++ {
						c <- strconv.Itoa(i)
					}
					close(c)
				}()
				return c
			},
			want: []string{"0", "1", "0", "1", "2", "0", "1", "2", "3"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			input := FromSlice(tc.input)
			mappedChan := FlatMap(input, tc.mappingFunc)
			got := ToSlice(mappedChan)
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("unexpected result (-got, +want): %s", diff)
			}
			// check that both channels are closed now
			_, ok := <-input
			if ok {
				t.Error("expected input to be closed ")
			}
			_, ok = <-mappedChan
			if ok {
				t.Error("expected mappedChan to be closed ")
			}
		})
	}
}

func TestFilter(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name       string
		input      []int
		filterFunc func(int) bool
		want       []int
	}{
		{
			name:  "filter_empty",
			input: []int{},
			filterFunc: func(i int) bool {
				t.Error("filter function was called when it should not have been")
				return true
			},
			want: nil,
		},
		{
			name:  "filter_one_true",
			input: []int{2},
			filterFunc: func(i int) bool {
				return i%2 == 0
			},
			want: []int{2},
		},
		{
			name:  "filter_one_false",
			input: []int{1},
			filterFunc: func(i int) bool {
				return i%2 == 0
			},
			want: nil,
		},
		{
			name:  "map_many",
			input: []int{1, 2, 3, 4, 5},
			filterFunc: func(i int) bool {
				return i%2 == 0
			},
			want: []int{2, 4},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			input := FromSlice(tc.input)
			filteredChan := Filter(input, tc.filterFunc)
			got := ToSlice(filteredChan)
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("unexpected result (-got, +want): %s", diff)
			}
			// check that both channels are closed now
			_, ok := <-input
			if ok {
				t.Error("expected input to be closed ")
			}
			_, ok = <-filteredChan
			if ok {
				t.Error("expected filteredChan to be closed ")
			}
		})
	}
}
