// Command ledgerctl is the institution-side admin CLI for examiner credential
// provisioning. See README.md for the design.
//
// The package exposes a Main(args, stdin, stdout, stderr) int function so
// tests can drive the CLI in-process via testkit/clitest. The production
// main() is a one-line wrapper around it.
package main

import (
	"fmt"
	"io"
	"os"

	"github.com/mmpworks/ffiec/cliutil"
)

const progName = "ledgerctl"

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
	return cliutil.ReportAndExit(progName, stderr, dispatch(args, stdout, stderr))
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
		return cliutil.Usagef("unknown command %q", args[0])
	}
}
