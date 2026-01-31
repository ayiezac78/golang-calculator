// main.go
package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	// "strings"

	calculator "github.com/ayiezac78/golang-calculator/calculations"
)

func main() {
	fmt.Println("Welcome to Azriel's Calculator")

	for {
		op := getOperation()
		num1, num2 := getNumbers()

		calc := calculator.Calculations{Num1: num1, Num2: num2}

		performOperation(op, &calc)

		wantsAnotherOperation()
	}
}

func getOperation() string {
	var op string
	for {
		fmt.Print("Choose operation (+, -, *, /): ")
		if _, err := fmt.Scanln(&op); err != nil {
			fmt.Println("Invalid input. Please try again.")
			continue
		}

		if op == "+" || op == "-" || op == "*" || op == "/" {
			return op
		}
		fmt.Println("Invalid operation. Please try again.")
	}
}

func getNumbers() (int, int) {
	var num1, num2 int
	fmt.Print("Input first number: ")
	if _, err := fmt.Scanln(&num1); err != nil {
		fmt.Println("Invalid input. Please try again.")
		getNumbers()
	}
	fmt.Print("Input second number: ")
	if _, err := fmt.Scanln(&num2); err != nil {
		fmt.Println("Invalid input. Please try again.")
		getNumbers()
	}
	return num1, num2
}

func performOperation(op string, calc *calculator.Calculations) {
	switch op {
	case "+":
		fmt.Println("Result:", calc.Sum())
	case "-":
		fmt.Println("Result:", calc.Subtract())
	case "*":
		fmt.Println("Result:", calc.Multiply())
	case "/":
		performDivision(calc)
	}
}

func performDivision(calc *calculator.Calculations) {
	result, err := calc.Divide()

	for errors.Is(err, calculator.ErrDivideByZero) {
		fmt.Println(err)
		fmt.Println("Please re-enter the numbers:")
		calc.Num1, calc.Num2 = getNumbers()
		result, err = calc.Divide()
	}

	if err != nil {
		fmt.Println("An unexpected error occurred:", err)
		return
	}

	fmt.Println("Result:", result)
}

func wantsAnotherOperation() bool {
	var response string

	for {
		fmt.Print("Perform another operation? (y/n): ")
		_, err := fmt.Scanln(&response)
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}

		// Convert to lowercase for easier comparison
		response = strings.ToLower(strings.TrimSpace(response))

		switch response {
		case "y":
			return true
		case "n":
			fmt.Println("Goodbye!")
			os.Exit(0)
		default:
			fmt.Println("Invalid input. Response should be (y/n)")
		}
	}
}
