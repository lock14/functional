package channel

import (
	"context"
	"errors"
	"sort"
	"testing"
)

func TestParallelMapWithErr(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c, errs := ParallelMapWithErr(ctx, 4, Of(1, 2, 3), func(i int) (int, error) {
		if i == 2 {
			return 0, errors.New("err2")
		}
		return i * 10, nil
	})
	var mapped []int
	var errList []error
	
	done := make(chan struct{})
	go func() {
		for e := range errs {
			errList = append(errList, e)
		}
		close(done)
	}()
	
	for i := range c {
		mapped = append(mapped, i)
	}
	<-done
	
	sort.Ints(mapped)
	if len(mapped) != 2 || len(errList) != 1 {
		t.Errorf("ParallelMapWithErr failed: %v, %v", mapped, errList)
	}
	
	c0, _ := ParallelMapWithErr(ctx, 0, Of(1), func(i int) (int, error) { return i, nil })
	<-c0
}

func TestParallelFlatMapWithErr(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c, errs := ParallelFlatMapWithErr(ctx, 4, Of(1, 2), func(i int) (chan int, error) {
		if i == 2 {
			return nil, errors.New("err2")
		}
		return Of(10, 11), nil
	})
	var mapped []int
	var errList []error
	
	done := make(chan struct{})
	go func() {
		for e := range errs {
			errList = append(errList, e)
		}
		close(done)
	}()
	
	for i := range c {
		mapped = append(mapped, i)
	}
	<-done
	
	sort.Ints(mapped)
	if len(mapped) != 2 || len(errList) != 1 {
		t.Errorf("ParallelFlatMapWithErr failed: %v, %v", mapped, errList)
	}
}

func TestParallelFilterWithErr(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c, errs := ParallelFilterWithErr(ctx, 4, Of(1, 2, 3), func(i int) (bool, error) {
		if i == 2 {
			return false, errors.New("err2")
		}
		return i == 3, nil
	})
	var filtered []int
	var errList []error
	
	done := make(chan struct{})
	go func() {
		for e := range errs {
			errList = append(errList, e)
		}
		close(done)
	}()
	
	for i := range c {
		filtered = append(filtered, i)
	}
	<-done
	
	if len(filtered) != 1 || filtered[0] != 3 || len(errList) != 1 {
		t.Errorf("ParallelFilterWithErr failed: %v, %v", filtered, errList)
	}
	
	c0, _ := ParallelFilterWithErr(ctx, 0, Of(1), func(i int) (bool, error) { return true, nil })
	<-c0
}
