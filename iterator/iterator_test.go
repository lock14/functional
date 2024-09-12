package iterator

import (
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/lock14/functional/slice"
	"iter"
	"maps"
	"slices"
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

			input := slices.Values(tc.input)
			mappedItr := Map(input, tc.mappingFunc)
			got := slices.Collect(mappedItr)
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("unexpected result (-got, +want): %s", diff)
			}
		})
	}
}

func TestFlatMap(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name        string
		input       []int
		mappingFunc func(int) iter.Seq[string]
		want        []string
	}{
		{
			name:  "map_empty",
			input: []int{},
			mappingFunc: func(i int) iter.Seq[string] {
				t.Error("mapping function was called when it should not have been")
				return nil
			},
			want: nil,
		},
		{
			name:  "map_one",
			input: []int{1},
			mappingFunc: func(n int) iter.Seq[string] {
				return func(yield func(string) bool) {
					for i := 0; i <= n; i++ {
						if !yield(strconv.Itoa(i)) {
							break
						}
					}
				}
			},
			want: []string{"0", "1"},
		},
		{
			name:  "map_many",
			input: []int{1, 2, 3},
			mappingFunc: func(n int) iter.Seq[string] {
				return func(yield func(string) bool) {
					for i := 0; i <= n; i++ {
						if !yield(strconv.Itoa(i)) {
							break
						}
					}
				}
			},
			want: []string{"0", "1", "0", "1", "2", "0", "1", "2", "3"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			input := slices.Values(tc.input)
			mappedItr := FlatMap(input, tc.mappingFunc)
			got := slices.Collect(mappedItr)
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("unexpected result (-got, +want): %s", diff)
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

			input := slices.Values(tc.input)
			filteredItr := Filter(input, tc.filterFunc)
			got := slices.Collect(filteredItr)
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("unexpected result (-got, +want): %s", diff)
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

			input := slices.Values(tc.input)
			got := FoldLeft(input, tc.foldingFunc, tc.initialValue)
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("unexpected result (-got, +want): %s", diff)
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

			input := slices.Values(tc.input)
			got := FoldRight(input, tc.foldingFunc, tc.initialValue)
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("unexpected result (-got, +want): %s", diff)
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

			input := slices.Values(tc.input)
			got := Reduce(input, tc.reducerFunc, tc.initialValue)
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("unexpected result (-got, +want): %s", diff)
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

			input := slices.Values(tc.input)
			got := Sum(input)
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("unexpected result (-got, +want): %s", diff)
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

			input := slices.Values(tc.input)
			got := JoinErrs(input)
			if diff := DiffErr(got, tc.want); diff != "" {
				t.Errorf("unexpected result (-got, +want): %s", diff)
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
			name:  "join_one",
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

			input := slices.Values(tc.input)
			got := Join(input, tc.sep)
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("unexpected result (-got, +want): %s", diff)
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
		wantLeft   []int
		wantRight  []string
	}{
		{
			name:       "both_empty",
			leftInput:  []int{},
			rightInput: []string{},
			wantLeft:   nil,
			wantRight:  nil,
		},
		{
			name:       "left_empty",
			leftInput:  []int{},
			rightInput: []string{"bob", "mary", "jane"},
			wantLeft:   nil,
			wantRight:  nil,
		},
		{
			name:       "right_empty",
			leftInput:  []int{1, 2, 3},
			rightInput: []string{},
			wantLeft:   nil,
			wantRight:  nil,
		},
		{
			name:       "left_shorter",
			leftInput:  []int{1},
			rightInput: []string{"bob", "mary", "jane"},
			wantLeft:   []int{1},
			wantRight:  []string{"bob"},
		},
		{
			name:       "right_shorter",
			leftInput:  []int{1, 2, 3},
			rightInput: []string{"bob"},
			wantLeft:   []int{1},
			wantRight:  []string{"bob"},
		},
		{
			name:       "same_length",
			leftInput:  []int{1, 2, 3},
			rightInput: []string{"bob", "mary", "jane"},
			wantLeft:   []int{1, 2, 3},
			wantRight:  []string{"bob", "mary", "jane"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			leftInput := slices.Values(tc.leftInput)
			rightInput := slices.Values(tc.rightInput)
			zipped := Zip(leftInput, rightInput)
			gotLeft, gotRight := slice.Collect(zipped)
			if diff := cmp.Diff(gotLeft, tc.wantLeft); diff != "" {
				t.Errorf("unexpected result (-got, +want): %s", diff)
			}
			if diff := cmp.Diff(gotRight, tc.wantRight); diff != "" {
				t.Errorf("unexpected result (-got, +want): %s", diff)
			}
		})
	}
}

func TestUnZip(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name      string
		input     map[int]string
		wantLeft  []int
		wantRight []string
	}{
		{
			name:      "empty",
			input:     map[int]string{},
			wantLeft:  nil,
			wantRight: nil,
		},
		{
			name: "one",
			input: map[int]string{
				1: "bob",
			},
			wantLeft:  []int{1},
			wantRight: []string{"bob"},
		},
		{
			name: "many",
			input: map[int]string{
				1: "bob",
				2: "mary",
				3: "jane",
			},
			wantLeft:  []int{1, 2, 3},
			wantRight: []string{"bob", "mary", "jane"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			input := maps.All(tc.input)
			unzippedLeft, unzippedRight := UnZip(input)
			gotLeft, gotRight := slices.Collect(unzippedLeft), slices.Collect(unzippedRight)
			if diff := cmp.Diff(gotLeft, tc.wantLeft); diff != "" {
				t.Errorf("unexpected result (-got, +want): %s", diff)
			}
			if diff := cmp.Diff(gotRight, tc.wantRight); diff != "" {
				t.Errorf("unexpected result (-got, +want): %s", diff)
			}
		})
	}
}

func TestSorted(t *testing.T) {
	t.Parallel()

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
			input := slices.Values(tc.input)
			sorted := Sorted(input)
			got := slices.Collect(sorted)
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("unexpected result (-got, +want): %s", diff)
			}
		})
	}
}

func TestDistinct(t *testing.T) {
	t.Parallel()

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
			input := slices.Values(tc.input)
			distinct := Distinct(input)
			got := slices.Collect(distinct)
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("unexpected result (-got, +want): %s", diff)
			}
		})
	}
}

type StatefulSupplier struct {
	state int
}

func (s *StatefulSupplier) Supply() int {
	value := s.state
	s.state++
	return value
}

func (s *StatefulSupplier) NumCalls() int {
	return s.state
}

func TestGenerate(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name      string
		supplier  *StatefulSupplier
		numReads  int
		want      []int
		wantCalls int
	}{
		{
			name:      "read_one",
			supplier:  &StatefulSupplier{},
			numReads:  1,
			want:      []int{0},
			wantCalls: 1,
		},
		{
			name:      "read_many",
			supplier:  &StatefulSupplier{},
			numReads:  10,
			want:      []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			wantCalls: 10,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			generator := Generate(tc.supplier.Supply)
			var got []int
			generator(func(i int) bool {
				got = append(got, i)
				return len(got) < tc.numReads
			})
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("unexpected result (-got, +want): %s", diff)
			}
			if diff := cmp.Diff(tc.supplier.NumCalls(), tc.wantCalls); diff != "" {
				t.Errorf("unexpected number of calls (-got, +want): %s", diff)
			}
		})
	}
}

func TestIterate(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		seed    int
		hasNext func(int) bool
		next    func(int) int
		want    []int
	}{
		{
			name:    "count",
			seed:    1,
			hasNext: func(i int) bool { return i <= 10 },
			next:    func(i int) int { return i + 1 },
			want:    []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := slices.Collect(Iterate(tc.seed, tc.hasNext, tc.next))
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("unexpected result (-got, +want): %s", diff)
			}
		})
	}
}

func TestRange(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name  string
		start int
		end   int
		want  []int
	}{
		{
			name:  "0_to_9",
			start: 0,
			end:   10,
			want:  []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := slices.Collect(Range(tc.start, tc.end))
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("unexpected result (-got, +want): %s", diff)
			}
		})
	}
}

func TestRangeClosed(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name  string
		start int
		end   int
		want  []int
	}{
		{
			name:  "1_to_10",
			start: 1,
			end:   10,
			want:  []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := slices.Collect(RangeClosed(tc.start, tc.end))
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("unexpected result (-got, +want): %s", diff)
			}
		})
	}
}

func TestLimit(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name  string
		input []int
		limit int64
		want  []int
	}{
		{
			name:  "limit_less_than_size",
			input: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			limit: 3,
			want:  []int{1, 2, 3},
		},
		{
			name:  "limit_equals_size",
			input: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			limit: 10,
			want:  []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		},
		{
			name:  "limit_greater_than_size",
			input: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			limit: 30,
			want:  []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := slices.Collect(Limit(slices.Values(tc.input), tc.limit))
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("unexpected result (-got, +want): %s", diff)
			}
		})
	}
}

func TestSkip(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name  string
		input []int
		skip  int64
		want  []int
	}{
		{
			name:  "skip_less_than_size",
			input: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			skip:  3,
			want:  []int{4, 5, 6, 7, 8, 9, 10},
		},
		{
			name:  "skip_equals_size",
			input: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			skip:  10,
			want:  nil,
		},
		{
			name:  "skip_greater_than_size",
			input: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			skip:  30,
			want:  nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := slices.Collect(Skip(slices.Values(tc.input), tc.skip))
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("unexpected result (-got, +want): %s", diff)
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

			input := slices.Values(tc.input)
			peekedIter := Peek(input, tc.statefulConsumer.Consume)
			got := slices.Collect(peekedIter)
			// make sure peek didn't mutate the data
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("unexpected result (-got, +want): %s", diff)
			}
			// make sure we consumed what we expected
			if diff := cmp.Diff(tc.statefulConsumer.Consumed(), tc.wantConsumed); diff != "" {
				t.Errorf("unexpected result for consumed (-got, +want): %s", diff)
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
