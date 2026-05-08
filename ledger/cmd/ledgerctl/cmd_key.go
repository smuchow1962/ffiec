package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/mmpworks/ffiec/ledger/internal/examiner"
)

func runKey(args []string, stdout, stderr io.Writer) error {
	if len(args) < 1 {
		return usagef("key: missing subcommand (generate)")
	}
	sub, rest := args[0], args[1:]
	switch sub {
	case "generate":
		return runKeyGenerate(rest, stdout, stderr)
	default:
		return usagef("key: unknown subcommand %q", sub)
	}
}

func runKeyGenerate(args []string, stdout, _ io.Writer) error {
	fs := flag.NewFlagSet("key generate", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	out := fs.String("out", "", "Output PEM path for the private key (required)")
	if stop, err := parseFlags(fs, args, stdout); err != nil {
		return err
	} else if stop {
		return nil
	}
	if *out == "" {
		return usagef("--out: required")
	}

	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return fmt.Errorf("generate: %w", err)
	}

	der, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return fmt.Errorf("marshal pkcs8: %w", err)
	}
	pemBytes := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der})
	if err := os.WriteFile(*out, pemBytes, 0o600); err != nil {
		return fmt.Errorf("write %s: %w", *out, err)
	}
	fmt.Fprintf(stdout, "wrote %s\nissuance_key_fingerprint: %s\n", *out, examiner.Fingerprint([]byte(pub)))
	return nil
}
