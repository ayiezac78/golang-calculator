package calculations

import "errors"

// ErrDivideByZero is returned when an attempt is made to divide by zero.
var ErrDivideByZero = errors.New("cannot divide by zero")

type Calculator struct{}

func (c *Calculator) Sum(x, y int) int {
	return x + y
}

func (c *Calculator) Subtract(x, y int) int {
	return x - y
}

func (c *Calculator) Multiply(x, y int) int {
	return x * y
}

func (c *Calculator) Divide(x, y float32) (float32, error) {
	if y == 0 {
		return 0, ErrDivideByZero
	}
	return x / y, nil
}
