package calculator

import (
	"errors"
	"testing"
)

func TestOperations(t *testing.T) {
	c := Calculations{Num1: 7, Num2: 2}

	if got := c.Sum(); got != 9 {
		t.Errorf("Sum() = %v, want 9", got)
	}
	if got := c.Subtract(); got != 5 {
		t.Errorf("Subtract() = %v, want 5", got)
	}
	if got := c.Multiply(); got != 14 {
		t.Errorf("Multiply() = %v, want 14", got)
	}
}

func TestDivide(t *testing.T) {
	tests := []struct {
		name    string
		a, b    float64
		want    float64
		wantErr error
	}{
		{"whole", 8, 2, 4, nil},
		{"fraction", 7, 2, 3.5, nil},
		{"negative", -9, 3, -3, nil},
		{"by zero", 5, 0, 0, ErrDivideByZero},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Calculations{Num1: tt.a, Num2: tt.b}.Divide()
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Divide() error = %v, want %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("Divide() = %v, want %v", got, tt.want)
			}
		})
	}
}
