package main

import (
	"flag"
	"fmt"
	"io"

	"github.com/mmpworks/ffiec/cliutil"
	"github.com/mmpworks/ffiec/ledger/internal/examiner"
)

func runKey(args []string, stdout, stderr io.Writer) error {
	if len(args) < 1 {
		return cliutil.Usagef("key: missing subcommand (generate)")
	}
	sub, rest := args[0], args[1:]
	switch sub {
	case "generate":
		return runKeyGenerate(rest, stdout, stderr)
	default:
		return cliutil.Usagef("key: unknown subcommand %q", sub)
	}
}

func runKeyGenerate(args []string, stdout, _ io.Writer) error {
	fs := flag.NewFlagSet("key generate", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	out := fs.String("out", "", "Output PEM path for the private key (required)")
	if stop, err := cliutil.ParseFlags(fs, args, stdout, progName); err != nil {
		return err
	} else if stop {
		return nil
	}
	if *out == "" {
		return cliutil.Usagef("--out: required")
	}

	pub, err := cliutil.GenerateEd25519PKCS8(*out)
	if err != nil {
		return err
	}
	fmt.Fprintf(stdout, "wrote %s\nissuance_key_fingerprint: %s\n", *out, examiner.Fingerprint([]byte(pub)))
	return nil
}
