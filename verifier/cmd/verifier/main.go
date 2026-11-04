// Command verifier is the offline standalone CLI an examiner runs against a
// ledger to verify chain integrity, Merkle proofs, and HSM-rooted seal
// signatures. See README.md for the design.
//
// The verifier takes only the ledger bytes and the institution's published
// seal-signing public key. It performs no network calls, ships as a single
// stdlib-only binary, and produces deterministic reports.
package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/mmpworks/ffiec/cliutil"
	"github.com/mmpworks/ffiec/verifier/internal/exitcodes"
)

const progName = "verifier"

const usageRoot = `verifier — offline chain-of-custody verification for examiners

Usage:
  verifier <command> [flags]

Commands:
  verify       Run the structural (and optional full) verification procedure,
               over a ledger, or over an NDJSON audit file with --audit-file.
  walk         Print every chain entry from a ledger in order.
  diff         Diff two ledger files; report any chain or seal divergence.
  gen-fixture  Generate a self-consistent demo ledger + seal-key pair.

Run "verifier <command> --help" for the flags of any command.
See README.md for the design.
`

func main() {
	os.Exit(Main(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

// Main is the testable entry point. Returns the process exit code:
//
//	0 on success
//	1 on a runtime error (including verification failure)
//	2 on a usage / argument error, or an audit file the walk cannot begin
//	3 on a key an audit-file walk cannot read (§10.12 configuration error)
func Main(args []string, _ io.Reader, stdout, stderr io.Writer) int {
	if len(args) < 1 {
		fmt.Fprint(stderr, usageRoot)
		return 2
	}
	err := dispatch(args, stdout, stderr)
	var walked *exitError
	if errors.As(err, &walked) {
		if walked.code == int(exitcodes.ConfigurationError) {
			fmt.Fprintf(stderr, "%s: %s\n", progName, walked.message)
		}
		return walked.code
	}
	return cliutil.ReportAndExit(progName, stderr, err)
}

func dispatch(args []string, stdout, stderr io.Writer) error {
	switch args[0] {
	case "verify":
		return runVerify(args[1:], stdout, stderr)
	case "walk":
		return runWalk(args[1:], stdout, stderr)
	case "diff":
		return runDiff(args[1:], stdout, stderr)
	case "gen-fixture":
		return runGenFixture(args[1:], stdout, stderr)
	case "-h", "--help", "help":
		fmt.Fprint(stdout, usageRoot)
		return nil
	default:
		return cliutil.Usagef("unknown command %q", args[0])
	}
}
