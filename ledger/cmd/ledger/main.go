// Command ledger is the institution-side reference ingest server. It
// receives OTLP from instrumented applications, verifies the HMAC chain on
// ingest, writes to an append-only ledger, computes the daily Merkle seal,
// and signs that seal in HSM custody.
//
// This is a bootstrap. The serve loop, OTLP receiver, storage backends,
// and HSM clients are scaffolded — the dispatch shape, config schema, and
// validation surface are real and exercised by the test suite.
package main

import (
	"fmt"
	"io"
	"os"

	"github.com/mmpworks/ffiec/cliutil"
)

const progName = "ledger"

const usageRoot = `ledger — institution-side reference ingest server

Usage:
  ledger <command> [flags]

Commands:
  serve            Start the OTLP receiver (bootstrap: scaffolded).
  config init      Write a default config file.
  config validate  Validate an existing config file.
  health           Validate config + readiness; exits 0 when ready to serve.

Run "ledger <command> --help" for the flags of any command.
See README.md for the design.
`

func main() {
	os.Exit(Main(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

// Main is the testable entry point. Returns the process exit code:
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
	case "serve":
		return runServe(args[1:], stdout, stderr)
	case "config":
		return runConfig(args[1:], stdout, stderr)
	case "health":
		return runHealth(args[1:], stdout, stderr)
	case "-h", "--help", "help":
		fmt.Fprint(stdout, usageRoot)
		return nil
	default:
		return cliutil.Usagef("unknown command %q", args[0])
	}
}
