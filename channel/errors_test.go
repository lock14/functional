package channel

import (
	"context"
	"errors"
	"testing"
)

func TestMapWithErr(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c, errs := MapWithErr(ctx, Of(1, 2, 3), func(i int) (int, error) {
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
	
	if len(mapped) != 2 || len(errList) != 1 {
		t.Errorf("MapWithErr failed: %v, %v", mapped, errList)
	}
}

func TestFlatMapWithErr(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c, errs := FlatMapWithErr(ctx, Of(1, 2), func(i int) (chan int, error) {
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
	
	if len(mapped) != 2 || len(errList) != 1 {
		t.Errorf("FlatMapWithErr failed: %v, %v", mapped, errList)
	}
}

func TestFilterWithErr(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c, errs := FilterWithErr(ctx, Of(1, 2, 3), func(i int) (bool, error) {
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
		t.Errorf("FilterWithErr failed: %v, %v", filtered, errList)
	}
}
