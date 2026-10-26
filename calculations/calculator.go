package calculator

import "errors"

// ErrDivideByZero is returned when an attempt is made to divide by zero.
var ErrDivideByZero = errors.New("cannot divide by zero")

// Calculations holds the two operands of a calculation.
type Calculations struct {
	Num1, Num2 float64
}

// Sum returns Num1 + Num2.
func (c Calculations) Sum() float64 {
	return c.Num1 + c.Num2
}

// Subtract returns Num1 - Num2.
func (c Calculations) Subtract() float64 {
	return c.Num1 - c.Num2
}

// Multiply returns Num1 * Num2.
func (c Calculations) Multiply() float64 {
	return c.Num1 * c.Num2
}

// Divide returns Num1 / Num2, or ErrDivideByZero when Num2 is zero.
func (c Calculations) Divide() (float64, error) {
	if c.Num2 == 0 {
		return 0, ErrDivideByZero
	}
	return c.Num1 / c.Num2, nil
}
