package basic

import (
	"bufio"
	"errors"
	"io"
	"math"
	"strings"
	"testing"

	"github.com/ayiezac78/golang-calculator/internal/testutil"
)

// scanner returns a Scanner that reads the given input.
func scanner(input string) *bufio.Scanner {
	return bufio.NewScanner(strings.NewReader(input))
}

func TestReadLine(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr error
	}{
		// positive
		{"plain line", "hello\n", "hello", nil},
		// negative
		{"no input", "", "", io.EOF},
		// edge
		{"trims whitespace", "  hi  \n", "hi", nil},
		{"empty line", "\n", "", nil},
		{"no trailing newline", "end", "end", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got string
			var err error
			out := testutil.CaptureStdout(t, func() { got, err = readLine(scanner(tt.input), "> ") })

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("readLine() error = %v, want %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("readLine() = %q, want %q", got, tt.want)
			}
			if out != "> " {
				t.Errorf("prompt = %q, want %q", out, "> ")
			}
		})
	}
}

func TestReadLineReturnsScannerError(t *testing.T) {
	// A line longer than the scanner's buffer makes Scan fail with ErrTooLong,
	// which must be returned as-is rather than reported as io.EOF.
	in := scanner("+++++++++\n")
	in.Buffer(make([]byte, 4), 4)

	var err error
	testutil.CaptureStdout(t, func() { _, err = readLine(in, "") })

	if !errors.Is(err, bufio.ErrTooLong) {
		t.Fatalf("readLine() error = %v, want %v", err, bufio.ErrTooLong)
	}
}

func TestReadNumber(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    float64
		wantErr error
	}{
		// positive
		{"integer", "42\n", 42, nil},
		{"decimal", "3.14\n", 3.14, nil},
		{"negative", "-7\n", -7, nil},
		// negative
		{"retries after letters", "abc\n5\n", 5, nil},
		{"retries after empty line", "\n5\n", 5, nil},
		{"no input", "", 0, io.EOF},
		{"only invalid input", "abc\nxyz\n", 0, io.EOF},
		// edge
		{"surrounding whitespace", "  8  \n", 8, nil},
		{"zero", "0\n", 0, nil},
		{"scientific notation", "1e3\n", 1000, nil},
		{"no trailing newline", "9", 9, nil},
		{"largest float", "1.7976931348623157e308\n", math.MaxFloat64, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got float64
			var err error
			testutil.CaptureStdout(t, func() { got, err = readNumber(scanner(tt.input), "") })

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("readNumber() error = %v, want %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("readNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadNumberPrintsRetryMessage(t *testing.T) {
	out := testutil.CaptureStdout(t, func() { readNumber(scanner("abc\n1\n"), "") })
	if !strings.Contains(out, "Invalid number. Please try again.") {
		t.Errorf("output = %q, want retry message", out)
	}
}

func TestGetOperation(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr error
	}{
		// positive
		{"add", "+\n", "+", nil},
		{"subtract", "-\n", "-", nil},
		{"multiply", "*\n", "*", nil},
		{"divide", "/\n", "/", nil},
		// negative
		{"retries after unknown operator", "%\n+\n", "+", nil},
		{"retries after word", "plus\n*\n", "*", nil},
		{"retries after doubled operator", "++\n-\n", "-", nil},
		{"no input", "", "", io.EOF},
		// edge
		{"surrounding whitespace", "  /  \n", "/", nil},
		{"retries after empty line", "\n+\n", "+", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got string
			var err error
			testutil.CaptureStdout(t, func() { got, err = getOperation(scanner(tt.input)) })

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("getOperation() error = %v, want %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("getOperation() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestWantsAnotherOperation(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    bool
		wantErr error
	}{
		// positive
		{"yes", "y\n", true, nil},
		{"no", "n\n", false, nil},
		// negative
		{"retries after word", "yes\nn\n", false, nil},
		{"retries after unrelated input", "maybe\ny\n", true, nil},
		{"no input", "", false, io.EOF},
		// edge
		{"uppercase Y", "Y\n", true, nil},
		{"uppercase N", "N\n", false, nil},
		{"surrounding whitespace", "  y  \n", true, nil},
		{"retries after empty line", "\nn\n", false, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got bool
			var err error
			testutil.CaptureStdout(t, func() { got, err = WantsAnotherOperation(scanner(tt.input)) })

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("WantsAnotherOperation() error = %v, want %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("WantsAnotherOperation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalculate(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantResult string // expected "Result: ..." line; empty when none is printed
		wantErr    error
	}{
		// positive
		{"add", "+\n2\n3\n", "Result: 5", nil},
		{"subtract", "-\n10\n4\n", "Result: 6", nil},
		{"multiply", "*\n6\n7\n", "Result: 42", nil},
		{"divide", "/\n7\n2\n", "Result: 3.5", nil},
		// negative
		{"divide by zero then valid", "/\n10\n0\n5\n", "Result: 2", nil},
		{"divide by zero repeatedly", "/\n10\n0\n0\n4\n", "Result: 2.5", nil},
		{"invalid inputs then valid", "x\n+\na\n1\nb\n2\n", "Result: 3", nil},
		{"input ends after operator", "+\n", "", io.EOF},
		{"input ends after first number", "+\n1\n", "", io.EOF},
		{"input ends after divide by zero", "/\n1\n0\n", "", io.EOF},
		// edge
		{"negative numbers", "*\n-3\n-4\n", "Result: 12", nil},
		{"zero divided by number", "/\n0\n5\n", "Result: 0", nil},
		{"subtract to negative", "-\n3\n5\n", "Result: -2", nil},
		{"overflow to infinity", "*\n1e308\n10\n", "Result: +Inf", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			out := testutil.CaptureStdout(t, func() { err = Calculate(scanner(tt.input)) })

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Calculate() error = %v, want %v", err, tt.wantErr)
			}
			if tt.wantResult == "" {
				if strings.Contains(out, "Result:") {
					t.Errorf("output = %q, want no result", out)
				}
				return
			}
			if !strings.Contains(out, tt.wantResult+"\n") {
				t.Errorf("output = %q, want %q", out, tt.wantResult)
			}
		})
	}
}

func TestCalculateDivideByZeroPrompts(t *testing.T) {
	out := testutil.CaptureStdout(t, func() { Calculate(scanner("/\n1\n0\n1\n")) })

	for _, want := range []string{"Error: cannot divide by zero", "Re-enter second number: "} {
		if !strings.Contains(out, want) {
			t.Errorf("output = %q, want it to contain %q", out, want)
		}
	}
}
