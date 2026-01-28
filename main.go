package main

import (
	"errors"
	"fmt"

	calculator "github.com/ayiezac78/golang-calculator/calculations"
)

func main() {
	var firstNum, secondNum int
	calc := calculator.Calculations{
		Num1: firstNum,
		Num2: secondNum,
	}

	appName := "Calculator"

	fmt.Printf("%s\n", appName)

	fmt.Print("Please choose operation: (+, -, *, /): ")
ops:
	for {
		var op string
		fmt.Scanln(&op)
		switch op {
		case "+":
			inputNum(&firstNum, &secondNum)
			fmt.Println("Result:", calc.Sum())
			break ops
		case "-":
			inputNum(&firstNum, &secondNum)
			fmt.Println("Result:", calc.Subtract())
			break ops
		case "*":
			inputNum(&firstNum, &secondNum)
			fmt.Println("Result:", calc.Multiply())
			break ops
		case "/":
			for {
				inputNum(&firstNum, &secondNum)
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
