package channel

import (
	"context"
	"sync"
	"testing"
)

func TestIterate(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c := Iterate(ctx, 1, func(i int) bool { return i < 4 }, func(i int) int { return i + 1 })
	s := ToSlice(ctx, c)
	if len(s) != 3 || s[0] != 1 || s[2] != 3 {
		t.Errorf("Iterate failed: %v", s)
	}
}

func TestRange(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c := Range(ctx, 1, 4)
	s := ToSlice(ctx, c)
	if len(s) != 3 || s[0] != 1 || s[2] != 3 {
		t.Errorf("Range failed: %v", s)
	}
}

func TestRangeClosed(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c := RangeClosed(ctx, 1, 3)
	s := ToSlice(ctx, c)
	if len(s) != 3 || s[0] != 1 || s[2] != 3 {
		t.Errorf("RangeClosed failed: %v", s)
	}
}

func TestLimit(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c := Limit(ctx, Range(ctx, 1, 10), 3)
	s := ToSlice(ctx, c)
	if len(s) != 3 || s[0] != 1 || s[2] != 3 {
		t.Errorf("Limit failed: %v", s)
	}
}

func TestSkip(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c := Skip(ctx, Range(ctx, 1, 5), 2)
	s := ToSlice(ctx, c)
	if len(s) != 2 || s[0] != 3 || s[1] != 4 {
		t.Errorf("Skip failed: %v", s)
	}
}

func TestAllMatch(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if !AllMatch(ctx, Range(ctx, 1, 4), func(i int) bool { return i < 10 }) {
		t.Errorf("AllMatch failed (expected true)")
	}
	var callCount int
	if AllMatch(ctx, Range(ctx, 1, 100), func(i int) bool {
		callCount++
		return i < 2
	}) {
		t.Errorf("AllMatch failed (expected false)")
	}
	if callCount > 3 {
		t.Errorf("AllMatch did not short-circuit: got %d calls", callCount)
	}
}

func TestAnyMatch(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if !AnyMatch(ctx, Range(ctx, 1, 4), func(i int) bool { return i == 2 }) {
		t.Errorf("AnyMatch failed (expected true)")
	}
	if AnyMatch(ctx, Range(ctx, 1, 4), func(i int) bool { return i == 5 }) {
		t.Errorf("AnyMatch failed (expected false)")
	}
	var callCount int
	if !AnyMatch(ctx, Range(ctx, 1, 100), func(i int) bool {
		callCount++
		return i == 2
	}) {
		t.Errorf("AnyMatch failed (expected true)")
	}
	if callCount > 3 {
		t.Errorf("AnyMatch did not short-circuit: got %d calls", callCount)
	}
}

func TestClone(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	clones := Clone(ctx, Of(1, 2, 3), 3)
	if len(clones) != 3 {
		t.Fatalf("expected 3 clones, got %d", len(clones))
	}

	results := make([][]int, len(clones))
	var wg sync.WaitGroup
	for i, ch := range clones {
		wg.Add(1)
		go func(idx int, c chan int) {
			defer wg.Done()
			results[idx] = ToSlice(ctx, c)
		}(i, ch)
	}
	wg.Wait()

	for _, res := range results {
		if len(res) != 3 || res[0] != 1 || res[1] != 2 || res[2] != 3 {
			t.Errorf("unexpected clone content: %v", res)
		}
	}

	c0 := Clone(ctx, Of(1), 0)
	if len(c0) != 0 {
		t.Errorf("expected 0 clones, got %d", len(c0))
	}
}

func TestTakeWhile(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c := TakeWhile(ctx, Range(ctx, 1, 10), func(i int) bool { return i < 4 })
	s := ToSlice(ctx, c)
	if len(s) != 3 || s[0] != 1 || s[2] != 3 {
		t.Errorf("TakeWhile failed: %v", s)
	}
}

func TestCount(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if Count(ctx, Range(ctx, 1, 4)) != 3 {
		t.Errorf("Count failed")
	}
}

func TestConcat(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c := Concat(ctx, Of(1, 2), Of(3, 4))
	s := ToSlice(ctx, c)
	if len(s) != 4 || s[0] != 1 || s[3] != 4 {
		t.Errorf("Concat failed: %v", s)
	}
}

func TestForEach(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var sum int
	ForEach(ctx, Range(ctx, 1, 4), func(i int) { sum += i })
	if sum != 6 {
		t.Errorf("ForEach failed: %d", sum)
	}
}

func TestPartition(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c := Partition(ctx, Range(ctx, 1, 6), 2)
	var slices [][]int
	for chunk := range c {
		slices = append(slices, chunk)
	}
	if len(slices) != 3 || len(slices[0]) != 2 || slices[2][0] != 5 {
		t.Errorf("Partition failed: %v", slices)
	}

	c2 := Partition(ctx, Range(ctx, 1, 6), 0)
	<-c2 // wait for close
}

func TestStream(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	seq := func(yield func(int) bool) {
		yield(1)
		yield(2)
	}
	c := Stream(ctx, seq)
	s := ToSlice(ctx, c)
	if len(s) != 2 || s[0] != 1 || s[1] != 2 {
		t.Errorf("Stream failed: %v", s)
	}
}

func TestEarlyExits(t *testing.T) {
	// Cancel immediately to test ctx.Done() branches
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// Just executing them with cancelled context to hit the ctx.Done() paths
	<-Map(ctx, Of(1, 2), func(i int) int { return i })
	<-Flatten(ctx, Of(Of(1)))
	<-Filter(ctx, Of(1), func(i int) bool { return true })
	FoldLeft(ctx, Of(1), func(u, i int) int { return i }, 0)
	FoldRight(ctx, Of(1), func(i, u int) int { return i }, 0)
	Join(ctx, Of("a"), "")
	<-Zip(ctx, Of(1), Of(2))
	ts, _ := UnZip(ctx, Of(Pair[int, int]{1, 2}))
	<-ts
	<-Sorted(ctx, Of(1))
	<-Distinct(ctx, Of(1))
	<-Generate(ctx, func() int { return 1 })
	<-Iterate(ctx, 1, func(i int) bool { return true }, func(i int) int { return i })
	<-Limit(ctx, Of(1), 1)
	<-Skip(ctx, Of(1), 1)
	<-TakeWhile(ctx, Of(1), func(i int) bool { return true })
	<-Concat(ctx, Of(1), Of(2))
	<-Peek(ctx, Of(1), func(i int) {})
	ForEach(ctx, Of(1), func(i int) {})
	<-Partition(ctx, Of(1), 1)
	clones := Clone(ctx, Of(1), 1)
	<-clones[0]

	seq := func(yield func(int) bool) { yield(1) }
	<-Stream(ctx, seq)
}
