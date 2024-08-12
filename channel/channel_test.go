package channel

import (
	"fmt"
	"github.com/google/go-cmp/cmp"
	"strconv"
	"strings"
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

func TestSum(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name  string
		input []int
		want  int
	}{
		{
			name:  "sum_empty",
			input: []int{},
			want:  0,
		},
		{
			name:  "sum_one",
			input: []int{1},
			want:  1,
		},
		{
			name:  "sum_many",
			input: []int{1, 2, 3},
			want:  6,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			input := FromSlice(tc.input)
			got := Sum(input)
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

func TestJoinErrs(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name  string
		input []error
		want  error
	}{
		{
			name:  "join_empty",
			input: []error{},
			want:  nil,
		},
		{
			name:  "join_one",
			input: []error{fmt.Errorf("err1")},
			want:  fmt.Errorf("err1"),
		},
		{
			name:  "join_many",
			input: []error{fmt.Errorf("err1"), fmt.Errorf("err2"), fmt.Errorf("err3")},
			want:  fmt.Errorf("err1\nerr2\nerr3"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			input := FromSlice(tc.input)
			got := JoinErrs(input)
			if diff := DiffErr(got, tc.want); diff != "" {
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

func TestJoin(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name  string
		input []string
		sep   string
		want  string
	}{
		{
			name:  "join_empty",
			input: []string{},
			sep:   ", ",
			want:  "",
		},
		{
			name:  "sum_one",
			input: []string{"a"},
			sep:   ", ",
			want:  "a",
		},
		{
			name:  "join_many",
			input: []string{"a", "b", "c"},
			sep:   ", ",
			want:  "a, b, c",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			input := FromSlice(tc.input)
			got := Join(input, tc.sep)
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

func TestZip(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name       string
		leftInput  []int
		rightInput []string
		want       []Pair[int, string]
	}{
		{
			name:       "both_empty",
			leftInput:  []int{},
			rightInput: []string{},
			want:       nil,
		},
		{
			name:       "left_empty",
			leftInput:  []int{},
			rightInput: []string{"bob", "mary", "jane"},
			want:       nil,
		},
		{
			name:       "right_empty",
			leftInput:  []int{1, 2, 3},
			rightInput: []string{},
			want:       nil,
		},
		{
			name:       "left_shorter",
			leftInput:  []int{1},
			rightInput: []string{"bob", "mary", "jane"},
			want: []Pair[int, string]{
				{Fst: 1, Snd: "bob"},
			},
		},
		{
			name:       "right_shorter",
			leftInput:  []int{1, 2, 3},
			rightInput: []string{"bob"},
			want: []Pair[int, string]{
				{Fst: 1, Snd: "bob"},
			},
		},
		{
			name:       "same_length",
			leftInput:  []int{1, 2, 3},
			rightInput: []string{"bob", "mary", "jane"},
			want: []Pair[int, string]{
				{Fst: 1, Snd: "bob"},
				{Fst: 2, Snd: "mary"},
				{Fst: 3, Snd: "jane"},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			leftInput := FromSlice(tc.leftInput)
			rightInput := FromSlice(tc.rightInput)
			zipped := Zip(leftInput, rightInput)
			got := ToSlice(zipped)
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("unexpected result (-got, +want): %s", diff)
			}
			// verify zip channel is closed
			_, ok := <-zipped
			if ok {
				t.Error("expected zipped to be closed ")
			}
		})
	}
}

func TestUnZip(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name      string
		input     []Pair[int, string]
		wantLeft  []int
		wantRight []string
	}{
		{
			name:      "empty",
			input:     []Pair[int, string]{},
			wantLeft:  nil,
			wantRight: nil,
		},
		{
			name: "one",
			input: []Pair[int, string]{
				{Fst: 1, Snd: "bob"},
			},
			wantLeft:  []int{1},
			wantRight: []string{"bob"},
		},
		{
			name: "many",
			input: []Pair[int, string]{
				{Fst: 1, Snd: "bob"},
				{Fst: 2, Snd: "mary"},
				{Fst: 3, Snd: "jane"},
			},
			wantLeft:  []int{1, 2, 3},
			wantRight: []string{"bob", "mary", "jane"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			input := FromSlice(tc.input)
			unzippedLeft, unzippedRight := UnZip(input)
			gotLeft, gotRight := ToSlice(unzippedLeft), ToSlice(unzippedRight)
			if diff := cmp.Diff(gotLeft, tc.wantLeft); diff != "" {
				t.Errorf("unexpected result (-got, +want): %s", diff)
			}
			if diff := cmp.Diff(gotRight, tc.wantRight); diff != "" {
				t.Errorf("unexpected result (-got, +want): %s", diff)
			}
			_, ok := <-input
			if ok {
				t.Error("expected input to be closed ")
			}
			_, ok = <-unzippedLeft
			if ok {
				t.Error("expected unzippedLeft to be closed ")
			}
			_, ok = <-unzippedRight
			if ok {
				t.Error("expected unzippedRight to be closed ")
			}
		})
	}
}

func TestSorted(t *testing.T) {
	cases := []struct {
		name  string
		input []int
		want  []int
	}{
		{
			name:  "empty",
			input: []int{},
			want:  nil,
		},
		{
			name:  "one",
			input: []int{1},
			want:  []int{1},
		},
		{
			name:  "many_sorted",
			input: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			want:  []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		},
		{
			name:  "many_reverse_sorted",
			input: []int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
			want:  []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		},
		{
			name:  "many_unordered",
			input: []int{2, 1, 3, 5, 10, 9, 6, 8, 4, 7},
			want:  []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			input := FromSlice(tc.input)
			sortedChan := Sorted(input)
			got := ToSlice(sortedChan)
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("unexpected result (-got, +want): %s", diff)
			}
			_, ok := <-input
			if ok {
				t.Error("expected input to be closed ")
			}
			_, ok = <-sortedChan
			if ok {
				t.Error("expected sortedChan to be closed ")
			}

		})
	}
}

func TestDistinct(t *testing.T) {
	cases := []struct {
		name  string
		input []int
		want  []int
	}{
		{
			name:  "empty",
			input: []int{},
			want:  nil,
		},
		{
			name:  "one",
			input: []int{1},
			want:  []int{1},
		},
		{
			name:  "many_distinct",
			input: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			want:  []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		},
		{
			name:  "many_duplicates",
			input: []int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 9, 2, 3, 6, 1, 7, 4, 5, 8, 10, 3, 4, 7, 10, 9, 6, 8, 1, 5, 2},
			want:  []int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			input := FromSlice(tc.input)
			distinctChan := Distinct(input)
			got := ToSlice(distinctChan)
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("unexpected result (-got, +want): %s", diff)
			}
			_, ok := <-input
			if ok {
				t.Error("expected input to be closed ")
			}
			_, ok = <-distinctChan
			if ok {
				t.Error("expected distinctChan to be closed ")
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

func DiffErr(got error, want error) string {
	if got == nil && want == nil {
		return ""
	}
	if got == nil {
		return fmt.Sprintf("got error <nil> but want an error containing %q", want)
	}
	if want == nil {
		return fmt.Sprintf("got error %q but want an error <nil>", got)
	}
	if gotMsg, wantMsg := got.Error(), want.Error(); !strings.Contains(gotMsg, wantMsg) {
		out := fmt.Sprintf("got error %q but want an error containing %q", gotMsg, want)

		// For long strings that will be hard to visually diff, include a diff.
		// Explanation of the &&'s and ||'s: if we're diffing a long error
		// message against a short one, a detailed diff isn't needed. The
		// difference will be obvious to the eye, and any extra message will
		// just be clutter. So only show the extra diff if the messages are both
		// long, or both multi-line.
		const msgLen = 20 // chosen arbitrarily
		bothAreLong := len(wantMsg) >= msgLen && len(gotMsg) >= msgLen
		bothAreMultiline := strings.Contains(wantMsg, "\n") && strings.Contains(gotMsg, "\n")
		if bothAreLong || bothAreMultiline {
			out += fmt.Sprintf("; diff was (-got,+want):\n%s", cmp.Diff(gotMsg, want))
		}
		return out
	}
	return ""
}
