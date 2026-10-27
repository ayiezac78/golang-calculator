package basic

import (
	"errors"
	"math"
	"testing"
)

func TestOperations(t *testing.T) {
	tests := []struct {
		name                    string
		a, b                    float64
		sum, subtract, multiply float64
	}{
		// positive
		{"positive numbers", 7, 2, 9, 5, 14},
		// negative
		{"negative numbers", -3, -4, -7, 1, 12},
		{"mixed signs", -5, 3, -2, -8, -15},
		// edge
		{"zeros", 0, 0, 0, 0, 0},
		{"fractions", 0.5, 0.25, 0.75, 0.25, 0.125},
		{"overflow", math.MaxFloat64, math.MaxFloat64, math.Inf(1), 0, math.Inf(1)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Basic{Num1: tt.a, Num2: tt.b}

			if got := c.Sum(); got != tt.sum {
				t.Errorf("Sum() = %v, want %v", got, tt.sum)
			}
			if got := c.Subtract(); got != tt.subtract {
				t.Errorf("Subtract() = %v, want %v", got, tt.subtract)
			}
			if got := c.Multiply(); got != tt.multiply {
				t.Errorf("Multiply() = %v, want %v", got, tt.multiply)
			}
		})
	}
}

func TestDivide(t *testing.T) {
	tests := []struct {
		name    string
		a, b    float64
		want    float64
		wantErr error
	}{
		// positive
		{"whole", 8, 2, 4, nil},
		{"fraction", 7, 2, 3.5, nil},
		// negative
		{"negative", -9, 3, -3, nil},
		{"by zero", 5, 0, 0, ErrDivideByZero},
		// edge
		{"zero by zero", 0, 0, 0, ErrDivideByZero},
		{"zero by number", 0, 5, 0, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Basic{Num1: tt.a, Num2: tt.b}.Divide()
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Divide() error = %v, want %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("Divide() = %v, want %v", got, tt.want)
			}
		})
	}
}
