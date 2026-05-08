package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/mmpworks/ffiec/cliutil"
	"github.com/mmpworks/ffiec/ledger/internal/serverconfig"
)

func runConfig(args []string, stdout, stderr io.Writer) error {
	if len(args) < 1 {
		return cliutil.Usagef("config: missing subcommand (init|validate)")
	}
	sub, rest := args[0], args[1:]
	switch sub {
	case "init":
		return runConfigInit(rest, stdout, stderr)
	case "validate":
		return runConfigValidate(rest, stdout, stderr)
	default:
		return cliutil.Usagef("config: unknown subcommand %q", sub)
	}
}

func runConfigInit(args []string, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("config init", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	out := fs.String("out", "", "Output config file path (required)")
	overwrite := fs.Bool("overwrite", false, "Overwrite if --out already exists")
	if stop, err := cliutil.ParseFlags(fs, args, stdout, progName); err != nil {
		return err
	} else if stop {
		return nil
	}
	if *out == "" {
		return cliutil.Usagef("--out: required")
	}
	if !*overwrite {
		if _, err := os.Stat(*out); err == nil {
			return fmt.Errorf("%s already exists; pass --overwrite to replace it", *out)
		}
	}
	c := serverconfig.Default()
	if err := serverconfig.Save(*out, &c); err != nil {
		return err
	}
	fmt.Fprintf(stderr, "wrote default config to %s\n", *out)
	return nil
}

func runConfigValidate(args []string, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("config validate", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	in := fs.String("config", "", "Config file path (required)")
	if stop, err := cliutil.ParseFlags(fs, args, stdout, progName); err != nil {
		return err
	} else if stop {
		return nil
	}
	if *in == "" {
		return cliutil.Usagef("--config: required")
	}
	c, err := serverconfig.Load(*in)
	if err != nil {
		return err
	}
	fmt.Fprintf(stdout, "config valid: version=%s storage=%s hsm=%s production=%t logging=%s\n",
		c.Version, c.Storage.Driver, c.HSM.Driver, c.HSM.ProductionMode, c.Logging.Level)
	return nil
}
