# `testkit` — shared test harnesses for ffiec command-line tools

A small workspace module that holds reusable test helpers. The first package
is `clitest` — a process-less harness for driving CLI tools through their
testable entry points.

## When to use it

Any tool in this workspace that ships a binary (`ledgerctl`, `verifier`,
`ledger`, future tools) gets a CLI test suite written against `testkit/clitest`.
That gives every tool the same harness, the same assertion vocabulary, and the
same fast in-process test loop.

## What `clitest` does

Runs your CLI's `Main(args, stdin, stdout, stderr) int` function directly,
captures stdout / stderr / exit code into a `Result`, and exposes assertion
helpers (`MustSucceed`, `MustExit`, `StdoutContains`, etc.) for tests.

No subprocess. No binary build per test. The test runs in the same process
the production binary uses, so coverage instrumentation reports honest numbers
and there is no per-test compile cost.

## How to plug a new tool in

Each tool exposes a single `Main` function with a fixed signature. The
production `main()` is one line that wires the OS streams to it.

### 1. Expose `Main` in your `cmd/<tool>/main.go`

```go
package main

import (
	"io"
	"os"
)

func main() {
	os.Exit(Main(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

// Main is the testable entry point. Returns the process exit code.
func Main(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	// ... dispatch, run, return 0 on success, 1 on runtime error,
	//     2 on usage error.
}
```

### 2. Write the tool's command handlers to take `io.Writer`s, not the OS streams

Every leaf function takes `(args []string, stdout, stderr io.Writer) error` —
or a subset if it doesn't need both streams. Avoid `os.Exit` inside the
handlers; return errors and let `Main` map them to exit codes. Avoid
`os.Stdout`, `os.Stderr`, and `fmt.Print*` (which target os streams) inside
the handlers; write to the supplied `io.Writer`s.

### 3. Add a `cmd/<tool>/main_test.go`

```go
package main

import (
	"testing"

	"github.com/mmpworks/ffiec/testkit/clitest"
)

func TestThing(t *testing.T) {
	clitest.Run(Main, "subcommand", "--flag", "value").
		MustSucceed(t).
		StdoutContains(t, "expected output")
}
```

### 4. Add the `testkit` dep to your tool's `go.mod`

```
require github.com/mmpworks/ffiec/testkit v0.0.0

replace github.com/mmpworks/ffiec/testkit => ../testkit
```

The `replace` directive matches the workspace convention used for the `core`
dep elsewhere in the repo.

## Reference plug-in

`ledger/cmd/ledgerctl/` is the first tool wired up this way. Read its
`main.go` (one-line wrapper) and `main_test.go` (full CLI suite) for a
working example.

## Adding more shared test packages

`clitest` is the first occupant. Future packages might include:

- `goldenfile` — diff a buffer against a checked-in golden file with an
  `--update` flag for regeneration.
- `chaintest` — synthetic ledger and synthetic chain-event fixtures shared
  across `core/`, `ledger/`, and `verifier/` tests.

Add them as siblings under `testkit/<package>/`. Keep the public surface
small; this module is for test code only.

## License

Apache 2.0. See `../LICENSE`.
