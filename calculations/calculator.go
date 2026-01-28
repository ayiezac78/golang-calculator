package calculator

import "errors"

type Arithmetic interface {
	Sum(x, y int) int
	Subtract(x, y int) int
	Multiply(x, y int) int
	Divide(x, y float32) (float32, error)
}

// ErrDivideByZero is returned when an attempt is made to divide by zero.
var ErrDivideByZero = errors.New("cannot divide by zero")

type Calculations struct {
	x, y int
}

func (c *Calculations) Sum(x, y int) int {
	return x + y
}

func (c *Calculations) Subtract(x, y int) int {
	return x - y
}

func (c *Calculations) Multiply(x, y int) int {
	return x * y
}

func (c *Calculations) Divide(x, y float32) (float32, error) {
	if y == 0 {
		return 0, ErrDivideByZero
	}
	return x / y, nil
}
