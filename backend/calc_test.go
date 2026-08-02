package main

import "testing"

func approxEq(a, b float64) bool {
	const eps = 1e-9
	if a == b {
		return true
	}
	diff := a - b
	if diff < 0 {
		diff = -diff
	}
	return diff < eps
}

func TestComputeBasic(t *testing.T) {
	if v, _ := Compute("add", 2, 3); !approxEq(v, 5) {
		t.Fatalf("add expected 5 got %v", v)
	}
	if v, _ := Compute("sub", 5, 2); !approxEq(v, 3) {
		t.Fatalf("sub expected 3 got %v", v)
	}
	if v, _ := Compute("mul", 3, 4); !approxEq(v, 12) {
		t.Fatalf("mul expected 12 got %v", v)
	}
	if v, _ := Compute("div", 10, 2); !approxEq(v, 5) {
		t.Fatalf("div expected 5 got %v", v)
	}
}

func TestDivByZero(t *testing.T) {
	if _, err := Compute("div", 1, 0); err == nil {
		t.Fatalf("expected division by zero error")
	}
}

func TestSqrtNegative(t *testing.T) {
	if _, err := Compute("sqrt", -4, 0); err == nil {
		t.Fatalf("expected sqrt of negative error")
	}
}

func TestPow(t *testing.T) {
	if v, _ := Compute("pow", 2, 8); !approxEq(v, 256) {
		t.Fatalf("pow expected 256 got %v", v)
	}
}
