package basic

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

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

// operations is the single source of truth for which operators are supported.
var operations = map[string]func(Calculator) (float64, error){
	"+": func(c Calculator) (float64, error) { return c.Sum(), nil },
	"-": func(c Calculator) (float64, error) { return c.Subtract(), nil },
	"*": func(c Calculator) (float64, error) { return c.Multiply(), nil },
	"/": func(c Calculator) (float64, error) { return c.Divide() },
}

func Calculate(in *bufio.Scanner) error {
	op, err := getOperation(in)
	if err != nil {
		return err
	}

	calc := Basic{}
	if calc.Num1, err = readNumber(in, "Input first number: "); err != nil {
		return err
	}
	if calc.Num2, err = readNumber(in, "Input second number: "); err != nil {
		return err
	}

	for {
		result, err := operations[op](calc)
		if errors.Is(err, ErrDivideByZero) {
			fmt.Println("Error:", err)
			if calc.Num2, err = readNumber(in, "Re-enter second number: "); err != nil {
				return err
			}
			continue
		}
		if err != nil {
			return err
		}

		fmt.Println("Result:", result)
		return nil
	}
}

// readLine prints the prompt and returns the next trimmed line of input.
// It returns io.EOF when there is no more input.
func readLine(in *bufio.Scanner, prompt string) (string, error) {
	fmt.Print(prompt)
	if !in.Scan() {
		if err := in.Err(); err != nil {
			return "", err
		}
		return "", io.EOF
	}
	return strings.TrimSpace(in.Text()), nil
}

func getOperation(in *bufio.Scanner) (string, error) {
	for {
		op, err := readLine(in, "Choose operation (+, -, *, /): ")
		if err != nil {
			return "", err
		}
		if _, ok := operations[op]; ok {
			return op, nil
		}
		fmt.Println("Invalid operation. Please try again.")
	}
}

func readNumber(in *bufio.Scanner, prompt string) (float64, error) {
	for {
		text, err := readLine(in, prompt)
		if err != nil {
			return 0, err
		}
		num, err := strconv.ParseFloat(text, 64)
		if err == nil {
			return num, nil
		}
		fmt.Println("Invalid number. Please try again.")
	}
}

func WantsAnotherOperation(in *bufio.Scanner) (bool, error) {
	for {
		response, err := readLine(in, "Perform another operation? (y/n): ")
		if err != nil {
			return false, err
		}

		switch strings.ToLower(response) {
		case "y":
			return true, nil
		case "n":
			return false, nil
		default:
			fmt.Println("Invalid input. Response should be (y/n)")
		}
	}
}
