// Package cliutil holds shared building blocks for the ffiec command-line
// tools: usage-error types, flag-parsing helpers, PEM key loading, and the
// exit-code mapping every tool's Main() applies.
//
// Each tool keeps its own dispatch and its own usage banner; cliutil owns
// the small pieces that would otherwise be copy-pasted across tools.
package cliutil

import (
	"errors"
	"fmt"
	"io"
)

// UsageError tags errors that should map to exit code 2 (wrong invocation).
// Anything else maps to exit code 1.
type UsageError struct {
	Msg string
}

// Error returns the underlying message.
func (u *UsageError) Error() string { return u.Msg }

// Usagef returns a *UsageError formatted from format and args.
func Usagef(format string, a ...interface{}) error {
	return &UsageError{Msg: fmt.Sprintf(format, a...)}
}

// IsUsage reports whether err (or anything it wraps) is a *UsageError.
func IsUsage(err error) bool {
	if err == nil {
		return false
	}
	var ue *UsageError
	return errors.As(err, &ue)
}

// ExitCodeFor maps a returned error to the standard ffiec exit code:
//
//	nil          → 0
//	*UsageError  → 2
//	anything else → 1
func ExitCodeFor(err error) int {
	if err == nil {
		return 0
	}
	if IsUsage(err) {
		return 2
	}
	return 1
}

// ReportAndExit prints err to stderr (prefixed with prog) and returns the
// exit code Main should propagate. It is the standard final line of every
// tool's Main():
//
//	return cliutil.ReportAndExit("ledgerctl", stderr, err)
func ReportAndExit(prog string, stderr io.Writer, err error) int {
	if err == nil {
		return 0
	}
	var ue *UsageError
	if errors.As(err, &ue) {
		fmt.Fprintf(stderr, "%s: %s\n", prog, ue.Msg)
		return 2
	}
	fmt.Fprintf(stderr, "%s: %v\n", prog, err)
	return 1
}
