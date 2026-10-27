package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"

	calculator "github.com/ayiezac78/golang-calculator/basic"
)

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
		if err := calculator.Calculate(in); err != nil {
			return err
		}

		again, err := calculator.WantsAnotherOperation(in)
		if err != nil || !again {
			return err
		}
	}
}
