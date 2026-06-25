package iterator

import (
	"testing"
)

func TestAllMatch(t *testing.T) {
	if !AllMatch(Range(1, 4), func(i int) bool { return i < 10 }) {
		t.Errorf("AllMatch failed (expected true)")
	}
	if AllMatch(Range(1, 4), func(i int) bool { return i < 2 }) {
		t.Errorf("AllMatch failed (expected false)")
	}
}

func TestAnyMatch(t *testing.T) {
	if !AnyMatch(Range(1, 4), func(i int) bool { return i == 2 }) {
		t.Errorf("AnyMatch failed (expected true)")
	}
	if AnyMatch(Range(1, 4), func(i int) bool { return i == 5 }) {
		t.Errorf("AnyMatch failed (expected false)")
	}
}

func TestCount(t *testing.T) {
	if Count(Range(1, 4)) != 3 {
		t.Errorf("Count failed")
	}
}

func TestConcat(t *testing.T) {
	c := Concat(Of(1, 2), Of(3, 4))
	var s []int
	for i := range c {
		s = append(s, i)
	}
	if len(s) != 4 || s[0] != 1 || s[3] != 4 {
		t.Errorf("Concat failed: %v", s)
	}
}

func TestOf(t *testing.T) {
	var s []int
	for i := range Of(1, 2, 3) {
		s = append(s, i)
	}
	if len(s) != 3 || s[0] != 1 || s[2] != 3 {
		t.Errorf("Of failed: %v", s)
	}
}

func TestPartition(t *testing.T) {
	c := Partition(Range(1, 6), 2)
	var slices [][]int
	for chunk := range c {
		var s []int
		for i := range chunk {
			s = append(s, i)
		}
		slices = append(slices, s)
	}
	if len(slices) != 3 || len(slices[0]) != 2 || slices[2][0] != 5 {
		t.Errorf("Partition failed: %v", slices)
	}
	
	c2 := Partition(Range(1, 6), 0)
	for range c2 {}
}

func TestEarlyExits(t *testing.T) {
	// Test early exit (yield returns false) to hit 100% on Map, Flatten, Filter, Iterate, Skip, Peek, Distinct
	
	var count int
	// Map early exit
	for _ = range Map(Range(1, 10), func(i int) int { return i * 2 }) {
		count++
		if count == 2 {
			break
		}
	}
	
	count = 0
	// Flatten early exit
	for _ = range Flatten(Of(Of(1, 2), Of(3, 4))) {
		count++
		if count == 2 {
			break
		}
	}
	
	count = 0
	// Filter early exit
	for _ = range Filter(Range(1, 10), func(i int) bool { return i%2 == 0 }) {
		count++
		if count == 2 {
			break
		}
	}
	
	count = 0
	// Iterate early exit
	for _ = range Iterate(1, func(i int) bool { return i < 10 }, func(i int) int { return i + 1 }) {
		count++
		if count == 2 {
			break
		}
	}
	
	count = 0
	// Skip early exit
	for _ = range Skip(Range(1, 10), 2) {
		count++
		if count == 2 {
			break
		}
	}
	
	count = 0
	// Peek early exit
	for _ = range Peek(Range(1, 10), func(i int) {}) {
		count++
		if count == 2 {
			break
		}
	}
	
	count = 0
	// Distinct early exit
	for _ = range Distinct(Of(1, 2, 2, 3, 4)) {
		count++
		if count == 2 {
			break
		}
	}
	
	count = 0
	// Partition early exit
	for chunk := range Partition(Range(1, 10), 2) {
		count++
		// internal chunk early exit
		for _ = range chunk {
			break // inner exit
		}
		if count == 2 {
			break // outer exit
		}
	}
}
