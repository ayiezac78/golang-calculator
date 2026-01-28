package calculator

import "errors"

type Arithmetic interface {
	Sum() int
	Subtract() int
	Multiply() int
	Divide() (float32, error)
}

// ErrDivideByZero is returned when an attempt is made to divide by zero.
var ErrDivideByZero = errors.New("cannot divide by zero")

type Calculations struct {
	Num1, Num2 int
}

func (c *Calculations) Sum() int {
	return c.Num1 + c.Num2
}

func (c *Calculations) Subtract() int {
	return c.Num1 - c.Num2
}

func (c *Calculations) Multiply() int {
	return c.Num1 * c.Num2
}

func (c *Calculations) Divide() (float32, error) {
	if c.Num1 == 0 {
		return 0, ErrDivideByZero
	}
	return float32(c.Num1) / float32(c.Num2), nil
}
