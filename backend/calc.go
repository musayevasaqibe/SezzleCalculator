package main

import (
	"fmt"
	"math"
)

func Compute(op string, a, b float64) (float64, error) {
	switch op {
	case "add", "+":
		return a + b, nil
	case "sub", "-":
		return a - b, nil
	case "mul", "*":
		return a * b, nil
	case "div", "/":
		if b == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		return a / b, nil
	case "pow":
		return math.Pow(a, b), nil
	case "sqrt":
		if a < 0 {
			return 0, fmt.Errorf("sqrt of negative number")
		}
		return math.Sqrt(a), nil
	case "pct": // a percent of b: a * b / 100
		return a * b / 100.0, nil
	default:
		return 0, fmt.Errorf("invalid operation: %s", op)
	}
}
