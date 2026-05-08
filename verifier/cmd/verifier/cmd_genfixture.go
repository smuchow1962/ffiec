package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/mmpworks/ffiec/cliutil"
	"github.com/mmpworks/ffiec/verifier/internal/verify"
)

// runGenFixture writes a self-consistent demo ledger and the matching
// seal-signing public key. The verifier uses this for smoke-tests and
// auditor demos: a one-shot way to produce an artifact that round-trips
// through the verify subcommand without standing up a ledger server.
func runGenFixture(args []string, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("gen-fixture", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	var (
		outLedger = fs.String("out-ledger", "", "Output path for the generated ledger NDJSON (required)")
		outPub    = fs.String("out-pub", "", "Output path for the seal-signing public key PEM (required)")
		outIKM    = fs.String("out-ikm", "", "Optional output path for the master IKM bytes (32B raw)")
		entries   = fs.Int("entries", 5, "Number of chain entries to generate")
		sealDate  = fs.String("seal-date", "2026-06-15", "Seal date (YYYY-MM-DD)")
	)
	if stop, err := cliutil.ParseFlags(fs, args, stdout, progName); err != nil {
		return err
	} else if stop {
		return nil
	}
	if err := cliutil.RequireFlags(map[string]string{
		"out-ledger": *outLedger,
		"out-pub":    *outPub,
	}); err != nil {
		return err
	}
	if *entries <= 0 {
		return cliutil.Usagef("--entries must be > 0")
	}

	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return fmt.Errorf("generate seal key: %w", err)
	}
	ikm := make([]byte, 32)
	if _, err := rand.Read(ikm); err != nil {
		return fmt.Errorf("generate ikm: %w", err)
	}

	fb := &verify.FixtureBuilder{
		TenantBindingKDFLabel: "tenant=demo;env=demo",
		SealDate:              *sealDate,
		MasterIKM:             ikm,
		SealingPriv:           priv,
	}

	chain := make([]verify.ChainEntry, 0, *entries)
	prev := ""
	for i := 0; i < *entries; i++ {
		payload := []byte(fmt.Sprintf(`{"event":"demo","seq":%d}`, i))
		e, err := fb.AppendEvent(prev, i, fmt.Sprintf("e_%04d", i), payload)
		if err != nil {
			return err
		}
		chain = append(chain, e)
		prev = e.EntryHash
	}
	seal, err := fb.SealEntries(chain)
	if err != nil {
		return err
	}
	raw, err := verify.MarshalNDJSON(chain, seal)
	if err != nil {
		return err
	}
	if err := os.WriteFile(*outLedger, raw, 0o600); err != nil {
		return fmt.Errorf("write ledger: %w", err)
	}
	if err := cliutil.WriteEd25519PublicPEM(*outPub, pub); err != nil {
		return fmt.Errorf("write public key: %w", err)
	}
	if *outIKM != "" {
		if err := os.WriteFile(*outIKM, ikm, 0o600); err != nil {
			return fmt.Errorf("write ikm: %w", err)
		}
	}
	fmt.Fprintf(stderr, "wrote ledger=%s pub=%s ikm=%s entries=%d\n",
		*outLedger, *outPub, *outIKM, *entries)
	return nil
}
