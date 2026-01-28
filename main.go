package main

import (
	"errors"
	"fmt"

	calculator "github.com/ayiezac78/golang-calculator/calculations"
)

var inputOperation string = "Please choose operation: (+, -, *, /): "

func main() {
	var firstNum, secondNum int
	var continueOperation = bool(false)
	calc := calculator.Calculations{
		Num1: firstNum,
		Num2: secondNum,
	}

	appName := "Welcome to Azriel's Calculator"

	fmt.Printf("%s\n", appName)

	fmt.Print(inputOperation)
ops:
	for {
		var op string
		fmt.Scanln(&op)
		switch op {
		case "+":
			inputNum(&calc.Num1, &calc.Num2)
			// calc.Num1 = firstNum
			// calc.Num2 = secondNum
			fmt.Println("Result:", calc.Sum())

			if !continueOperation {
				wantsAnotherOperation(inputOperation)
				continue
			}

			break ops
		case "-":
			inputNum(&calc.Num1, &calc.Num2)
			// calc.Num1 = firstNum
			// calc.Num2 = secondNum
			fmt.Println("Result:", calc.Subtract())

			break ops
		case "*":
			inputNum(&calc.Num1, &calc.Num2)
			fmt.Println("Result:", calc.Multiply())
			break ops
		case "/":
			for {
				inputNum(&calc.Num1, &calc.Num2)
				result, err := calc.Divide()
				if err != nil {
					if errors.Is(err, calculator.ErrDivideByZero) {
						fmt.Println(err)
						fmt.Println("Please re-input the numbers:")
						continue // prompt numbers again
					}
					// unexpected error
					fmt.Println("An unexpected error occurred:", err)
					break ops
				}
				fmt.Println("Result:", result)

				break ops
			}
		default:
			fmt.Println("Invalid operation")
			break ops
		}

	}
}

func inputNum(a, b *int) {
	fmt.Print("Input first number: ")
	fmt.Scanln(a)
	fmt.Print("Input second number: ")
	fmt.Scanln(b)
}

func wantsAnotherOperation(inputOperation string) bool {
	var response string
	fmt.Print("Do you want to perform another operation? (y/n): ")
	fmt.Scanln(&response)
	if response == "y" {
		fmt.Print(inputOperation)
		return true
	} else {
		return false
	}
	// return false
}
