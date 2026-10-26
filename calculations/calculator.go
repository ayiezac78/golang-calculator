package calculator

import "errors"

// ErrDivideByZero is returned when an attempt is made to divide by zero.
var ErrDivideByZero = errors.New("cannot divide by zero")

// Calculator holds the two operands of a calculation.
type Calculator interface {
	Sum() float64
	Subtract() float64
	Multiply() float64
	Divide() (float64, error)
}

type Basic struct {
	Num1, Num2 float64
}

var _ Calculator = Basic{} // ensure Basic implements Calculator

// Sum returns Num1 + Num2.
func (c Basic) Sum() float64 {
	return c.Num1 + c.Num2
}

// Subtract returns Num1 - Num2.
func (c Basic) Subtract() float64 {
	return c.Num1 - c.Num2
}

// Multiply returns Num1 * Num2.
func (c Basic) Multiply() float64 {
	return c.Num1 * c.Num2
}

// Divide returns Num1 / Num2, or ErrDivideByZero when Num2 is zero.
func (c Basic) Divide() (float64, error) {
	if c.Num2 == 0 {
		return 0, ErrDivideByZero
	}
	return c.Num1 / c.Num2, nil
}
