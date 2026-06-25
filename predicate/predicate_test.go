package predicate

import (
	"errors"
	"testing"
)

func TestIsNil(t *testing.T) {
	var p *int
	if !IsNil(p) {
		t.Errorf("Expected IsNil to be true for nil pointer")
	}
	i := 5
	if IsNil(&i) {
		t.Errorf("Expected IsNil to be false for non-nil pointer")
	}
	
	if IsNil(i) {
		t.Errorf("Expected IsNil to be false for non-pointer type")
	}
	
	var err error
	if !IsNil(err) {
		t.Errorf("Expected IsNil to be true for nil interface")
	}
	err = errors.New("err")
	if IsNil(err) {
		t.Errorf("Expected IsNil to be false for non-nil interface")
	}
}

func TestNotNil(t *testing.T) {
	var p *int
	if NotNil(p) {
		t.Errorf("Expected NotNil to be false for nil pointer")
	}
	i := 5
	if !NotNil(&i) {
		t.Errorf("Expected NotNil to be true for non-nil pointer")
	}
}

func TestNot(t *testing.T) {
	p := func(i int) bool { return i%2 == 0 }
	notP := Not(p)
	if notP(2) {
		t.Errorf("Expected Not(p) to be false for 2")
	}
	if !notP(3) {
		t.Errorf("Expected Not(p) to be true for 3")
	}
}

func TestTrue(t *testing.T) {
	if !True(5) {
		t.Errorf("Expected True to be true")
	}
}

func TestFalse(t *testing.T) {
	if False(5) {
		t.Errorf("Expected False to be false")
	}
}
