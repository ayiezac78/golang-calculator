package main

import (
	"errors"
	"fmt"

	calculator "github.com/ayiezac78/golang-calculator/calculations"
)

var inputOperation string = "Please choose operation: (+, -, *, /): "
var firstNum, secondNum int

var calc = calculator.Calculations{
	Num1: firstNum,
	Num2: secondNum,
}

func main() {

	appName := "Welcome to Azriel's Calculator"

	fmt.Printf("%s\n", appName)

	fmt.Print(inputOperation)
ops:
	for {
		var op string
		fmt.Scanln(&op)
		switch op {
		case "+":
			sum(&calc.Num1, &calc.Num2)
			break ops
		case "-":
			subtract(&calc.Num1, &calc.Num2)
			break ops
		case "*":
			multiply(&calc.Num1, &calc.Num2)
			break ops
		case "/":
			for {
				divide(&calc.Num1, &calc.Num2)
				break ops
			}
		default:
			fmt.Println("Invalid operation")
			break ops
		}

	}
	// fmt.Println(inputOperation)

	if wantsAnotherOperation() {
		goto ops
	}
}

func inputNum(a, b *int) {
	fmt.Print("Input first number: ")
	fmt.Scanln(a)
	fmt.Print("Input second number: ")
	fmt.Scanln(b)
}

func sum(num1, num2 *int) {
	inputNum(num1, num2)
	fmt.Println("Result:", calc.Sum())
}

func subtract(num1, num2 *int) {
	inputNum(num1, num2)
	fmt.Println("Result:", calc.Subtract())
}

func multiply(num1, num2 *int) {
	inputNum(num1, num2)
	fmt.Println("Result:", calc.Multiply())
}

func divide(num1, num2 *int) (float32, error) {
	inputNum(num1, num2)
	result, err := calc.Divide()
	if err != nil {
		if errors.Is(err, calculator.ErrDivideByZero) {
			fmt.Println(err)
			fmt.Println("Please re-input the numbers:")
			return divide(num1, num2) // prompt numbers again
		}
		// unexpected error
		fmt.Println("An unexpected error occurred:", err)
		return 0, err
	}

	fmt.Println("Result:", result)
	return result, nil

}

func wantsAnotherOperation() bool {
	var response string
	fmt.Println("Do you want to perform another operation? (y/n) ")
	fmt.Println("Press Enter to exit.")
	fmt.Scanln(&response)
	if response == "y" {
		fmt.Print(inputOperation)
		return true
	}
	return false
}
