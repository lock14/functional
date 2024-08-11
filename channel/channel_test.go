package channel

import (
	"github.com/google/go-cmp/cmp"
	"strconv"
	"testing"
)

func TestMap2(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name        string
		input       []int
		mappingFunc func(int) string
		want        []string
	}{
		{
			name:        "map_works",
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
