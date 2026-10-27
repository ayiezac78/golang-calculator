// Package testutil holds helpers shared by the project's tests.
package testutil

import (
	"io"
	"os"
	"testing"
	"time"
)

// CaptureStdout runs fn and returns everything it printed to stdout.
func CaptureStdout(t *testing.T, fn func()) string {
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

	// The input loops retry until input ends, so a bug there spins forever.
	// Abort the whole test binary instead of hanging.
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
