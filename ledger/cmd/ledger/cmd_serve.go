package main

import (
	"errors"
	"flag"
	"fmt"
	"io"

	"github.com/mmpworks/ffiec/cliutil"
	"github.com/mmpworks/ffiec/ledger/internal/serverconfig"
)

// runServe is the server entrypoint. At bootstrap it loads and validates
// the config and prints a "ready" line summarising what it would do, then
// returns a not-yet-implemented runtime error so an operator running it
// against a real environment knows the server loop hasn't landed.
//
// The shape is the same shape the production server will use: load config,
// build the ingestion + storage + sealing pipeline, run until signal. That
// keeps the wiring honest as the loop fills in.
func runServe(args []string, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("serve", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	cfgPath := fs.String("config", "", "Config file path (required)")
	dryRun := fs.Bool("dry-run", false, "Validate config + print readiness; do not start the receiver")
	if stop, err := cliutil.ParseFlags(fs, args, stdout, progName); err != nil {
		return err
	} else if stop {
		return nil
	}
	if *cfgPath == "" {
		return cliutil.Usagef("--config: required")
	}
	c, err := serverconfig.Load(*cfgPath)
	if err != nil {
		return err
	}
	printReadiness(stdout, c)
	if *dryRun {
		return nil
	}
	return errors.New("serve: server loop not yet implemented (bootstrap). Use --dry-run to validate config and exit cleanly.")
}

func printReadiness(w io.Writer, c *serverconfig.Config) {
	fmt.Fprintln(w, "ledger readiness:")
	fmt.Fprintf(w, "  config_version: %s\n", c.Version)
	fmt.Fprintf(w, "  otlp_grpc:      %s\n", emptyAsDash(c.OTLP.GRPCAddr))
	fmt.Fprintf(w, "  otlp_http:      %s\n", emptyAsDash(c.OTLP.HTTPAddr))
	fmt.Fprintf(w, "  storage:        driver=%s path=%q dsn=%q\n",
		c.Storage.Driver, c.Storage.Path, redactDSN(c.Storage.DSN))
	fmt.Fprintf(w, "  hsm:            driver=%s production=%t\n",
		c.HSM.Driver, c.HSM.ProductionMode)
	fmt.Fprintf(w, "  logging_level:  %s\n", c.Logging.Level)
}

func emptyAsDash(s string) string {
	if s == "" {
		return "-"
	}
	return s
}

// redactDSN replaces a non-empty DSN with the literal "<set>" so the
// readiness output never prints credentials. Empty stays empty so the
// operator sees the field is unset.
func redactDSN(dsn string) string {
	if dsn == "" {
		return ""
	}
	return "<set>"
}
