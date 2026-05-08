// Command ledgerctl is the institution-side admin CLI for examiner credential
// provisioning. See README.md for the design.
//
// The package exposes a Main(args, stdin, stdout, stderr) int function so
// tests can drive the CLI in-process via testkit/clitest. The production
// main() is a one-line wrapper around it.
package main

import (
	"errors"
	"fmt"
	"io"
	"os"
)

const usageRoot = `ledgerctl — institution-side admin CLI for examiner credential provisioning

Usage:
  ledgerctl <command> [flags]

Commands:
  examiner issue   Mint a credential bundle for an arriving examiner.
  examiner revoke  Revoke an active bundle before its expiry.
  examiner list    List currently-active examination credentials.
  examiner audit   Reconstruct issuance/revocation history from the chain.

  key generate     Generate a fresh Ed25519 issuance keypair (PEM PKCS#8).

Run "ledgerctl <command> --help" for the flags of any command.
See README.md for the design and the credential-bundle schema.
`

// usageError tags errors that map to exit code 2 (wrong invocation).
// Runtime errors (everything else) map to exit code 1.
type usageError struct{ msg string }

func (u *usageError) Error() string { return u.msg }

func usagef(format string, a ...interface{}) error {
	return &usageError{msg: fmt.Sprintf(format, a...)}
}

func main() {
	os.Exit(Main(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

// Main is the testable entry point. It returns the process exit code:
//
//	0 on success
//	1 on a runtime error
//	2 on a usage / argument error
func Main(args []string, _ io.Reader, stdout, stderr io.Writer) int {
	if len(args) < 1 {
		fmt.Fprint(stderr, usageRoot)
		return 2
	}
	err := dispatch(args, stdout, stderr)
	if err == nil {
		return 0
	}
	var ue *usageError
	if errors.As(err, &ue) {
		fmt.Fprintf(stderr, "ledgerctl: %s\n", ue.msg)
		return 2
	}
	fmt.Fprintf(stderr, "ledgerctl: %v\n", err)
	return 1
}

func dispatch(args []string, stdout, stderr io.Writer) error {
	switch args[0] {
	case "examiner":
		return runExaminer(args[1:], stdout, stderr)
	case "key":
		return runKey(args[1:], stdout, stderr)
	case "-h", "--help", "help":
		fmt.Fprint(stdout, usageRoot)
		return nil
	default:
		return usagef("unknown command %q", args[0])
	}
}
