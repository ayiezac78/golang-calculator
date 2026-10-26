package main

import (
	"bufio"
	"errors"
	"io"
	"math"
	"os"
	"strings"
	"testing"
	"testing/iotest"
	"time"
)

// scanner returns a Scanner that reads the given input.
func scanner(input string) *bufio.Scanner {
	return bufio.NewScanner(strings.NewReader(input))
}

// captureStdout runs fn and returns everything it printed to stdout.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	orig := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = orig }()

	out := make(chan string)
	go func() {
		b, _ := io.ReadAll(r)
		out <- string(b)
	}()

	// The input loops retry until readLine returns an error, so a bug there
	// spins forever. Abort the whole test binary instead of hanging.
	done := make(chan struct{})
	go func() {
		fn()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(time.Second):
		os.Stdout = orig
		panic(t.Name() + ": did not finish; an input loop never ended")
	}

	w.Close()
	return <-out
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
			captureStdout(t, func() { got, err = readNumber(scanner(tt.input), "") })

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
	out := captureStdout(t, func() { readNumber(scanner("abc\n1\n"), "") })
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
			captureStdout(t, func() { got, err = getOperation(scanner(tt.input)) })

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
			captureStdout(t, func() { got, err = wantsAnotherOperation(scanner(tt.input)) })

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("wantsAnotherOperation() error = %v, want %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("wantsAnotherOperation() = %v, want %v", got, tt.want)
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
			out := captureStdout(t, func() { err = calculate(scanner(tt.input)) })

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("calculate() error = %v, want %v", err, tt.wantErr)
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
	out := captureStdout(t, func() { calculate(scanner("/\n1\n0\n1\n")) })

	for _, want := range []string{"Error: cannot divide by zero", "Re-enter second number: "} {
		if !strings.Contains(out, want) {
			t.Errorf("output = %q, want it to contain %q", out, want)
		}
	}
}

func TestRun(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantResults []string
		wantErr     error
	}{
		// positive
		{"single calculation then quit", "+\n1\n2\nn\n", []string{"Result: 3"}, nil},
		{"two calculations then quit", "+\n1\n2\ny\n*\n3\n4\nn\n", []string{"Result: 3", "Result: 12"}, nil},
		// negative
		{"input ends before answering", "+\n1\n2\n", []string{"Result: 3"}, io.EOF},
		{"input ends during second calculation", "+\n1\n2\ny\n-\n", []string{"Result: 3"}, io.EOF},
		{"invalid answer then quit", "+\n1\n2\nok\nn\n", []string{"Result: 3"}, nil},
		// edge
		{"empty input", "", nil, io.EOF},
		{"uppercase answers", "+\n1\n1\nY\n-\n1\n1\nN\n", []string{"Result: 2", "Result: 0"}, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			out := captureStdout(t, func() { err = run(scanner(tt.input)) })

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("run() error = %v, want %v", err, tt.wantErr)
			}
			if got := strings.Count(out, "Result:"); got != len(tt.wantResults) {
				t.Errorf("printed %d results, want %d; output = %q", got, len(tt.wantResults), out)
			}
			for _, want := range tt.wantResults {
				if !strings.Contains(out, want+"\n") {
					t.Errorf("output = %q, want %q", out, want)
				}
			}
		})
	}
}

func TestApp(t *testing.T) {
	tests := []struct {
		name       string
		stdin      io.Reader
		wantCode   int
		wantStderr string
		wantBye    bool
	}{
		// positive
		{"user quits", strings.NewReader("+\n1\n2\nn\n"), 0, "", true},
		// negative
		{"read error", iotest.ErrReader(errors.New("boom")), 1, "Error reading input: boom\n", false},
		// edge
		{"input ends (EOF)", strings.NewReader("+\n1\n2\n"), 0, "", true},
		{"empty input", strings.NewReader(""), 0, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stderr strings.Builder
			var code int
			out := captureStdout(t, func() { code = app(tt.stdin, &stderr) })

			if code != tt.wantCode {
				t.Errorf("app() = %d, want %d", code, tt.wantCode)
			}
			if stderr.String() != tt.wantStderr {
				t.Errorf("stderr = %q, want %q", stderr.String(), tt.wantStderr)
			}
			if got := strings.Contains(out, "Goodbye!"); got != tt.wantBye {
				t.Errorf("printed Goodbye = %v, want %v; output = %q", got, tt.wantBye, out)
			}
		})
	}
}

func TestRunReturnsScannerError(t *testing.T) {
	// A line longer than the scanner's buffer makes Scan fail with ErrTooLong,
	// which must be returned as-is rather than reported as io.EOF.
	in := scanner("+++++++++\n")
	in.Buffer(make([]byte, 4), 4)

	var err error
	captureStdout(t, func() { err = run(in) })

	if !errors.Is(err, bufio.ErrTooLong) {
		t.Fatalf("run() error = %v, want %v", err, bufio.ErrTooLong)
	}
}
