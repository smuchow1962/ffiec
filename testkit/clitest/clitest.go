// Package clitest is a process-less test harness for command-line tools in
// the ffiec workspace.
//
// Each tool exposes a Main(args, stdin, stdout, stderr) int function — the
// production main() is one line that wires the os streams to it. Tests then
// invoke Main directly through clitest.Run and assert on captured output.
//
// The harness is process-less by design — no subprocess spawn, no binary
// build per test run — so suites are fast and stay in the same process the
// production binary uses, which keeps coverage instrumentation honest.
//
// Usage shape for a tool's main_test.go:
//
//	func TestSomething(t *testing.T) {
//	    clitest.Run(Main, "examiner", "issue", "--out", path).
//	        MustSucceed(t).
//	        StdoutContains(t, "issued bnd_")
//	}
package clitest

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Main is the function signature every CLI tool exposes for testing.
// The production main() is a thin wrapper:
//
//	func main() {
//	    os.Exit(Main(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
//	}
type Main func(args []string, stdin io.Reader, stdout, stderr io.Writer) int

// Result captures the output of a single CLI invocation.
type Result struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

// Run invokes m with args (program name omitted) and an empty stdin.
func Run(m Main, args ...string) Result {
	return RunWithStdin(m, "", args...)
}

// RunWithStdin invokes m with the supplied stdin string.
func RunWithStdin(m Main, stdin string, args ...string) Result {
	var so, se bytes.Buffer
	code := m(args, strings.NewReader(stdin), &so, &se)
	return Result{Stdout: so.String(), Stderr: se.String(), ExitCode: code}
}

// MustSucceed asserts the run exited 0. Returns r for chaining.
func (r Result) MustSucceed(t *testing.T) Result {
	t.Helper()
	if r.ExitCode != 0 {
		t.Fatalf("expected exit 0, got %d\nstderr: %s\nstdout: %s", r.ExitCode, r.Stderr, r.Stdout)
	}
	return r
}

// MustFail asserts the run exited non-zero. Returns r for chaining.
func (r Result) MustFail(t *testing.T) Result {
	t.Helper()
	if r.ExitCode == 0 {
		t.Fatalf("expected non-zero exit, got 0\nstdout: %s\nstderr: %s", r.Stdout, r.Stderr)
	}
	return r
}

// MustExit asserts the run exited with the given code. Returns r for chaining.
func (r Result) MustExit(t *testing.T, code int) Result {
	t.Helper()
	if r.ExitCode != code {
		t.Fatalf("expected exit %d, got %d\nstderr: %s\nstdout: %s", code, r.ExitCode, r.Stderr, r.Stdout)
	}
	return r
}

// StdoutContains asserts stdout contains substr. Returns r for chaining.
func (r Result) StdoutContains(t *testing.T, substr string) Result {
	t.Helper()
	if !strings.Contains(r.Stdout, substr) {
		t.Fatalf("stdout does not contain %q\nstdout: %s", substr, r.Stdout)
	}
	return r
}

// StderrContains asserts stderr contains substr. Returns r for chaining.
func (r Result) StderrContains(t *testing.T, substr string) Result {
	t.Helper()
	if !strings.Contains(r.Stderr, substr) {
		t.Fatalf("stderr does not contain %q\nstderr: %s", substr, r.Stderr)
	}
	return r
}

// StdoutLacks asserts stdout does NOT contain substr. Returns r for chaining.
func (r Result) StdoutLacks(t *testing.T, substr string) Result {
	t.Helper()
	if strings.Contains(r.Stdout, substr) {
		t.Fatalf("stdout unexpectedly contains %q\nstdout: %s", substr, r.Stdout)
	}
	return r
}

// TempPath returns a path under t.TempDir, creating parent directories.
// The returned path is cleaned up when the test ends (t.TempDir contract).
func TempPath(t *testing.T, name string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), name)
	if dir := filepath.Dir(p); dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatalf("MkdirAll %s: %v", dir, err)
		}
	}
	return p
}

// ReadFile reads p and fails the test on error.
func ReadFile(t *testing.T, p string) []byte {
	t.Helper()
	raw, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("read %s: %v", p, err)
	}
	return raw
}

// WriteFile writes content to p with mode 0600 and fails the test on error.
func WriteFile(t *testing.T, p string, content []byte) {
	t.Helper()
	if dir := filepath.Dir(p); dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatalf("MkdirAll %s: %v", dir, err)
		}
	}
	if err := os.WriteFile(p, content, 0o600); err != nil {
		t.Fatalf("write %s: %v", p, err)
	}
}
