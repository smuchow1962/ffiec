package main

import (
	"flag"
	"fmt"
	"io"

	"github.com/mmpworks/ffiec/cliutil"
	"github.com/mmpworks/ffiec/ledger/internal/serverconfig"
)

// runHealth performs the same validation `serve` does, but never tries to
// start anything. Suitable as a Kubernetes startup probe or a
// pre-deployment smoke. Exits 0 if the config validates and would let the
// server come up; non-zero with a localised message otherwise.
func runHealth(args []string, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("health", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	cfgPath := fs.String("config", "", "Config file path (required)")
	if stop, err := cliutil.ParseFlags(fs, args, stdout, progName); err != nil {
		return err
	} else if stop {
		return nil
	}
	if *cfgPath == "" {
		return cliutil.Usagef("--config: required")
	}
	if _, err := serverconfig.Load(*cfgPath); err != nil {
		return err
	}
	fmt.Fprintln(stdout, "ok")
	return nil
}
