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
			name:  "filter_many",
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

func TestFoldLeft(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name         string
		input        []int
		initialValue string
		foldingFunc  func(string, int) string
		want         string
	}{
		{
			name:         "fold_left_empty",
			input:        []int{},
			initialValue: "init",
			foldingFunc: func(s string, i int) string {
				t.Error("folding function was called when it should not have been")
				return "bad-value"
			},
			want: "init",
		},
		{
			name:         "fold_left_one",
			input:        []int{1},
			initialValue: "init",
			foldingFunc: func(s string, i int) string {
				return s + strconv.Itoa(i)
			},
			want: "init1",
		},
		{
			name:         "fold_left_many",
			input:        []int{1, 2, 3},
			initialValue: "",
			foldingFunc: func(s string, i int) string {
				return s + strconv.Itoa(i)
			},
			want: "123",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			input := FromSlice(tc.input)
			got := FoldLeft(input, tc.foldingFunc, tc.initialValue)
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("unexpected result (-got, +want): %s", diff)
			}
			// check that channel is closed now
			_, ok := <-input
			if ok {
				t.Error("expected input to be closed ")
			}
		})
	}
}

func TestFoldRight(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name         string
		input        []int
		initialValue string
		foldingFunc  func(i int, s string) string
		want         string
	}{
		{
			name:         "fold_right_empty",
			input:        []int{},
			initialValue: "init",
			foldingFunc: func(i int, s string) string {
				t.Error("folding function was called when it should not have been")
				return "bad-value"
			},
			want: "init",
		},
		{
			name:         "fold_right_one",
			input:        []int{1},
			initialValue: "init",
			foldingFunc: func(i int, s string) string {
				return strconv.Itoa(i) + s
			},
			want: "1init",
		},
		{
			name:         "fold_right_many",
			input:        []int{1, 2, 3},
			initialValue: "",
			foldingFunc: func(i int, s string) string {
				return strconv.Itoa(i) + s
			},
			want: "123",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			input := FromSlice(tc.input)
			got := FoldRight(input, tc.foldingFunc, tc.initialValue)
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("unexpected result (-got, +want): %s", diff)
			}
			// check that channel is closed now
			_, ok := <-input
			if ok {
				t.Error("expected input to be closed ")
			}
		})
	}
}

func TestReduce(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name         string
		input        []string
		initialValue string
		reducerFunc  func(string, string) string
		want         string
	}{
		{
			name:         "reduce_empty",
			input:        []string{},
			initialValue: "init",
			reducerFunc: func(s1, s2 string) string {
				t.Error("folding function was called when it should not have been")
				return "bad-value"
			},
			want: "init",
		},
		{
			name:         "reduce_one",
			input:        []string{"1"},
			initialValue: "init",
			reducerFunc: func(s1, s2 string) string {
				return s1 + s2
			},
			want: "init1",
		},
		{
			name:         "reduce_many",
			input:        []string{"1", "2", "3"},
			initialValue: "",
			reducerFunc: func(s1, s2 string) string {
				return s1 + s2
			},
			want: "123",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			input := FromSlice(tc.input)
			got := Reduce(input, tc.reducerFunc, tc.initialValue)
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("unexpected result (-got, +want): %s", diff)
			}
			// check that channel is closed now
			_, ok := <-input
			if ok {
				t.Error("expected input to be closed ")
			}
		})
	}
}

type StatefulConsumer[T any] struct {
	consumed []T
}

func (c *StatefulConsumer[T]) Consume(s T) {
	c.consumed = append(c.consumed, s)
}

func (c *StatefulConsumer[T]) Consumed() []T {
	return c.consumed
}

func TestPeek(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name             string
		input            []int
		statefulConsumer *StatefulConsumer[int]
		want             []int
		wantConsumed     []int
	}{
		{
			name:             "peek_empty",
			input:            []int{},
			statefulConsumer: &StatefulConsumer[int]{},
			want:             nil,
			wantConsumed:     nil,
		},
		{
			name:             "peek_one",
			input:            []int{2},
			statefulConsumer: &StatefulConsumer[int]{},
			want:             []int{2},
			wantConsumed:     []int{2},
		},
		{
			name:             "peek_many",
			input:            []int{1, 2, 3, 4, 5},
			statefulConsumer: &StatefulConsumer[int]{},
			want:             []int{1, 2, 3, 4, 5},
			wantConsumed:     []int{1, 2, 3, 4, 5},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			input := FromSlice(tc.input)
			peekedChan := Peek(input, tc.statefulConsumer.Consume)
			got := ToSlice(peekedChan)
			// make sure peek didn't mutate the data
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("unexpected result (-got, +want): %s", diff)
			}
			// make sure we consumed what we expected
			if diff := cmp.Diff(tc.statefulConsumer.Consumed(), tc.wantConsumed); diff != "" {
				t.Errorf("unexpected result for consumed (-got, +want): %s", diff)
			}
			// check that both channels are closed now
			_, ok := <-input
			if ok {
				t.Error("expected input to be closed ")
			}
			_, ok = <-peekedChan
			if ok {
				t.Error("expected peekedChan to be closed ")
			}
		})
	}
}
