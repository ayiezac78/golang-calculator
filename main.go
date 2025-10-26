package main

import (
	"fmt"
)

func main() {

	var operation string
	var firstNum int
	var secondNum int
	appName := "Calculator"

	fmt.Printf("%s\n", appName)
	fmt.Println("Please choose operation: (+, -, *, /)")
	fmt.Scanln(&operation)

	switch operation {
	case "+":
		fmt.Println("Input first number")
		fmt.Scanln(&firstNum)
		fmt.Println("Input second number")
		fmt.Scanln(&secondNum)
		fmt.Println(firstNum + secondNum)
	case "-":
		fmt.Println("Input first number")
		fmt.Scanln(&firstNum)
		fmt.Println("Input second number")
		fmt.Scanln(&secondNum)
		fmt.Println(firstNum - secondNum)
	case "*":
		fmt.Println("Input first number")
		fmt.Scanln(&firstNum)
		fmt.Println("Input second number")
		fmt.Scanln(&secondNum)
		fmt.Println(firstNum * secondNum)
	case "/":
		fmt.Println("Input first number")
		fmt.Scanln(&firstNum)
		fmt.Println("Input second number")
		fmt.Scanln(&secondNum)
		fmt.Println(float32(firstNum) / float32(secondNum))
	default:
		fmt.Println("Invalid operation")
	}
}
