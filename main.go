package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	calculator "github.com/ayiezac78/golang-calculator/calculations"
)

// operations is the single source of truth for which operators are supported.
var operations = map[string]func(calculator.Calculator) (float64, error){
	"+": func(c calculator.Calculator) (float64, error) { return c.Sum(), nil },
	"-": func(c calculator.Calculator) (float64, error) { return c.Subtract(), nil },
	"*": func(c calculator.Calculator) (float64, error) { return c.Multiply(), nil },
	"/": calculator.Calculator.Divide,
}

func main() {
	os.Exit(app(os.Stdin, os.Stderr))
}

// app runs the calculator on stdin and returns the process exit code.
func app(stdin io.Reader, stderr io.Writer) int {
	fmt.Println("Welcome to Zairelle's Calculator")

	if err := run(bufio.NewScanner(stdin)); err != nil && !errors.Is(err, io.EOF) {
		fmt.Fprintln(stderr, "Error reading input:", err)
		return 1
	}

	fmt.Println("\nGoodbye!")
	return 0
}

// run keeps calculating until the user says no or input ends (io.EOF).
func run(in *bufio.Scanner) error {
	for {
		if err := calculate(in); err != nil {
			return err
		}

		again, err := wantsAnotherOperation(in)
		if err != nil || !again {
			return err
		}
	}
}

func calculate(in *bufio.Scanner) error {
	op, err := getOperation(in)
	if err != nil {
		return err
	}

	calc := calculator.Basic{}
	if calc.Num1, err = readNumber(in, "Input first number: "); err != nil {
		return err
	}
	if calc.Num2, err = readNumber(in, "Input second number: "); err != nil {
		return err
	}

	for {
		result, err := operations[op](calc)
		if errors.Is(err, calculator.ErrDivideByZero) {
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

func wantsAnotherOperation(in *bufio.Scanner) (bool, error) {
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
