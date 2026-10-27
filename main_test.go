package main

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/ayiezac78/golang-calculator/internal/testutil"
)

func TestApp(t *testing.T) {
	errRead := errors.New("read failed")

	tests := []struct {
		name       string
		stdin      io.Reader
		wantCode   int
		wantOut    []string // substrings expected on stdout, in order
		wantStderr string
	}{
		// positive
		{"add then quit", strings.NewReader("+\n2\n3\nn\n"), 0,
			[]string{"Welcome to Zairelle's Calculator", "Result: 5", "Goodbye!"}, ""},
		{"subtract", strings.NewReader("-\n10\n4\nn\n"), 0, []string{"Result: 6"}, ""},
		{"multiply", strings.NewReader("*\n6\n7\nn\n"), 0, []string{"Result: 42"}, ""},
		{"divide", strings.NewReader("/\n7\n2\nn\n"), 0, []string{"Result: 3.5"}, ""},
		{"two operations", strings.NewReader("+\n1\n1\ny\n*\n3\n3\nN\n"), 0,
			[]string{"Result: 2", "Result: 9", "Goodbye!"}, ""},
		// negative
		{"invalid operation retried", strings.NewReader("%\n+\n1\n2\nn\n"), 0,
			[]string{"Invalid operation. Please try again.", "Result: 3"}, ""},
		{"invalid number retried", strings.NewReader("+\nabc\n1\n2\nn\n"), 0,
			[]string{"Invalid number. Please try again.", "Result: 3"}, ""},
		{"divide by zero retried", strings.NewReader("/\n8\n0\n2\nn\n"), 0,
			[]string{"Error: cannot divide by zero", "Result: 4"}, ""},
		{"invalid y/n retried", strings.NewReader("+\n1\n1\nmaybe\nn\n"), 0,
			[]string{"Invalid input. Response should be (y/n)", "Goodbye!"}, ""},
		{"read error", iotest.ErrReader(errRead), 1,
			[]string{"Welcome to Zairelle's Calculator"}, "Error reading input: read failed"},
		// edge
		{"empty input exits cleanly", strings.NewReader(""), 0, []string{"Goodbye!"}, ""},
		{"input ends mid-calculation", strings.NewReader("+\n1\n"), 0, []string{"Goodbye!"}, ""},
		{"input ends at y/n prompt", strings.NewReader("+\n1\n1\n"), 0,
			[]string{"Result: 2", "Goodbye!"}, ""},
		{"negative numbers", strings.NewReader("+\n-3\n-4\nn\n"), 0, []string{"Result: -7"}, ""},
		{"whitespace around input", strings.NewReader("  *  \n 2 \n 5 \n n \n"), 0,
			[]string{"Result: 10"}, ""},
		{"line too long", strings.NewReader(strings.Repeat("1", bufio.MaxScanTokenSize+1)), 1,
			nil, "Error reading input: " + bufio.ErrTooLong.Error()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stderr bytes.Buffer
			var code int
			out := testutil.CaptureStdout(t, func() { code = app(tt.stdin, &stderr) })

			if code != tt.wantCode {
				t.Errorf("app() = %d, want %d", code, tt.wantCode)
			}

			rest := out
			for _, want := range tt.wantOut {
				i := strings.Index(rest, want)
				if i < 0 {
					t.Fatalf("stdout missing %q (in order); got:\n%s", want, out)
				}
				rest = rest[i+len(want):]
			}

			if got := strings.TrimSpace(stderr.String()); got != tt.wantStderr {
				t.Errorf("stderr = %q, want %q", got, tt.wantStderr)
			}
			if tt.wantCode != 0 && strings.Contains(out, "Goodbye!") {
				t.Errorf("stdout printed Goodbye! on failure:\n%s", out)
			}
		})
	}
}

func TestRun(t *testing.T) {
	errRead := errors.New("read failed")

	tests := []struct {
		name    string
		stdin   io.Reader
		wantErr error
	}{
		// positive
		{"stops on n", strings.NewReader("+\n1\n2\nn\n"), nil},
		{"loops on y then stops", strings.NewReader("+\n1\n2\ny\n-\n5\n3\nn\n"), nil},
		// negative
		{"read error", iotest.ErrReader(errRead), errRead},
		// edge
		{"no input", strings.NewReader(""), io.EOF},
		{"ends after y", strings.NewReader("+\n1\n2\ny\n"), io.EOF},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			testutil.CaptureStdout(t, func() { err = run(bufio.NewScanner(tt.stdin)) })

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("run() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

// TestRunCountsOperations checks that run performs exactly one calculation
// per "y" answer plus the first, and no more.
func TestRunCountsOperations(t *testing.T) {
	input := "+\n1\n1\ny\n+\n2\n2\ny\n+\n3\n3\nn\n"
	out := testutil.CaptureStdout(t, func() {
		if err := run(bufio.NewScanner(strings.NewReader(input))); err != nil {
			t.Errorf("run() error = %v", err)
		}
	})

	if got := strings.Count(out, "Result:"); got != 3 {
		t.Errorf("results printed = %d, want 3; got:\n%s", got, out)
	}
}
