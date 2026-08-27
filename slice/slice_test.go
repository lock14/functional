package slice

import (
	"errors"
	"iter"
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMap(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name  string
		input []int
		f     func(int) string
		want  []string
	}{
		{
			name:  "empty",
			input: []int{},
			f:     strconv.Itoa,
			want:  []string{},
		},
		{
			name:  "single",
			input: []int{1},
			f:     strconv.Itoa,
			want:  []string{"1"},
		},
		{
			name:  "multiple",
			input: []int{1, 2, 3},
			f:     strconv.Itoa,
			want:  []string{"1", "2", "3"},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := Map(tc.input, tc.f)
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("Map mismatch (-got +want):\n%s", diff)
			}
		})
	}
}

func TestFlatten(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name  string
		input [][]int
		want  []int
	}{
		{
			name:  "empty",
			input: [][]int{},
			want:  []int{},
		},
		{
			name:  "nested",
			input: [][]int{{1, 2}, {3, 4}},
			want:  []int{1, 2, 3, 4},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := Flatten(tc.input)
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("Flatten mismatch (-got +want):\n%s", diff)
			}
		})
	}
}

func TestFlatMap(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name  string
		input []int
		f     func(int) []int
		want  []int
	}{
		{
			name:  "empty",
			input: []int{},
			f:     func(i int) []int { return []int{i, i * 10} },
			want:  []int{},
		},
		{
			name:  "multiple",
			input: []int{1, 2},
			f:     func(i int) []int { return []int{i, i * 10} },
			want:  []int{1, 10, 2, 20},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := FlatMap(tc.input, tc.f)
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("FlatMap mismatch (-got +want):\n%s", diff)
			}
		})
	}
}

func TestFilter(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name  string
		input []int
		p     func(int) bool
		want  []int
	}{
		{
			name:  "empty",
			input: []int{},
			p:     func(i int) bool { return i%2 == 0 },
			want:  nil,
		},
		{
			name:  "filter_evens",
			input: []int{1, 2, 3, 4},
			p:     func(i int) bool { return i%2 == 0 },
			want:  []int{2, 4},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := Filter(tc.input, tc.p)
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("Filter mismatch (-got +want):\n%s", diff)
			}
		})
	}
}

func TestFoldLeft(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		input   []int
		initial string
		f       func(string, int) string
		want    string
	}{
		{
			name:    "empty",
			input:   []int{},
			initial: "init",
			f:       func(u string, i int) string { return u + strconv.Itoa(i) },
			want:    "init",
		},
		{
			name:    "multiple",
			input:   []int{1, 2, 3},
			initial: "",
			f:       func(u string, i int) string { return u + strconv.Itoa(i) },
			want:    "123",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := FoldLeft(tc.input, tc.f, tc.initial)
			if got != tc.want {
				t.Errorf("FoldLeft = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestFoldRight(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		input   []int
		initial string
		f       func(int, string) string
		want    string
	}{
		{
			name:    "empty",
			input:   []int{},
			initial: "init",
			f:       func(i int, u string) string { return u + strconv.Itoa(i) },
			want:    "init",
		},
		{
			name:    "multiple",
			input:   []int{1, 2, 3},
			initial: "",
			f:       func(i int, u string) string { return u + strconv.Itoa(i) },
			want:    "321",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := FoldRight(tc.input, tc.f, tc.initial)
			if got != tc.want {
				t.Errorf("FoldRight = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestReduce(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		input   []int
		op      func(int, int) int
		initial int
		want    int
	}{
		{
			name:    "sum_reduction",
			input:   []int{1, 2, 3},
			op:      func(a, b int) int { return a + b },
			initial: 0,
			want:    6,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := Reduce(tc.input, tc.op, tc.initial)
			if got != tc.want {
				t.Errorf("Reduce = %d, want %d", got, tc.want)
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
		{name: "empty", input: []int{}, want: 0},
		{name: "numbers", input: []int{1, 2, 3}, want: 6},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := Sum(tc.input)
			if got != tc.want {
				t.Errorf("Sum = %d, want %d", got, tc.want)
			}
		})
	}
}

func TestJoinErrs(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name  string
		input []error
		want  string
	}{
		{name: "empty", input: []error{}, want: ""},
		{name: "errors", input: []error{errors.New("a"), errors.New("b")}, want: "a\nb"},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := JoinErrs(tc.input)
			if tc.want == "" && got != nil {
				t.Errorf("JoinErrs expected nil, got %v", got)
			} else if tc.want != "" && (got == nil || got.Error() != tc.want) {
				t.Errorf("JoinErrs = %v, want %q", got, tc.want)
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
		{name: "empty", input: []string{}, sep: ",", want: ""},
		{name: "single", input: []string{"a"}, sep: ",", want: "a"},
		{name: "multiple", input: []string{"a", "b", "c"}, sep: ",", want: "a,b,c"},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := Join(tc.input, tc.sep)
			if got != tc.want {
				t.Errorf("Join = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestZip(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		input1 []int
		input2 []string
		want   []Pair[int, string]
	}{
		{
			name:   "matching_length",
			input1: []int{1, 2},
			input2: []string{"a", "b"},
			want:   []Pair[int, string]{{Fst: 1, Snd: "a"}, {Fst: 2, Snd: "b"}},
		},
		{
			name:   "unequal_length",
			input1: []int{1, 2, 3},
			input2: []string{"a", "b"},
			want:   []Pair[int, string]{{Fst: 1, Snd: "a"}, {Fst: 2, Snd: "b"}},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := Zip(tc.input1, tc.input2)
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("Zip mismatch (-got +want):\n%s", diff)
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
			name:      "pairs",
			input:     []Pair[int, string]{{Fst: 1, Snd: "a"}, {Fst: 2, Snd: "b"}},
			wantLeft:  []int{1, 2},
			wantRight: []string{"a", "b"},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			gotLeft, gotRight := UnZip(tc.input)
			if diff := cmp.Diff(gotLeft, tc.wantLeft); diff != "" {
				t.Errorf("UnZip left mismatch (-got +want):\n%s", diff)
			}
			if diff := cmp.Diff(gotRight, tc.wantRight); diff != "" {
				t.Errorf("UnZip right mismatch (-got +want):\n%s", diff)
			}
		})
	}
}

func TestConcat(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		input1 []int
		input2 []int
		want   []int
	}{
		{
			name:   "concat_slices",
			input1: []int{1, 2},
			input2: []int{3, 4},
			want:   []int{1, 2, 3, 4},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := Concat(tc.input1, tc.input2)
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("Concat mismatch (-got +want):\n%s", diff)
			}
		})
	}
}

func TestPartition(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name  string
		input []int
		size  int
		want  [][]int
	}{
		{
			name:  "size_2",
			input: []int{1, 2, 3, 4, 5},
			size:  2,
			want:  [][]int{{1, 2}, {3, 4}, {5}},
		},
		{
			name:  "size_0",
			input: []int{1, 2, 3},
			size:  0,
			want:  [][]int{},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := Partition(tc.input, tc.size)
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("Partition mismatch (-got +want):\n%s", diff)
			}
		})
	}
}

func TestCollect(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name      string
		seq       iter.Seq2[int, string]
		wantLeft  []int
		wantRight []string
	}{
		{
			name: "pairs_seq",
			seq: func(yield func(int, string) bool) {
				if !yield(1, "a") {
					return
				}
				yield(2, "b")
			},
			wantLeft:  []int{1, 2},
			wantRight: []string{"a", "b"},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			gotLeft, gotRight := Collect(tc.seq)
			if diff := cmp.Diff(gotLeft, tc.wantLeft); diff != "" {
				t.Errorf("Collect left mismatch (-got +want):\n%s", diff)
			}
			if diff := cmp.Diff(gotRight, tc.wantRight); diff != "" {
				t.Errorf("Collect right mismatch (-got +want):\n%s", diff)
			}
		})
	}
}
