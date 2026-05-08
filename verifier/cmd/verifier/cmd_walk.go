package main

import (
	"flag"
	"fmt"
	"io"

	"github.com/mmpworks/ffiec/cliutil"
	"github.com/mmpworks/ffiec/verifier/internal/verify"
)

func runWalk(args []string, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("walk", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	var (
		ledgerPath = fs.String("ledger", "", "Path to the ledger file (required)")
		runID      = fs.String("run-id", "", "Optional run id; reserved for future selective walks")
	)
	if stop, err := cliutil.ParseFlags(fs, args, stdout, progName); err != nil {
		return err
	} else if stop {
		return nil
	}
	if *ledgerPath == "" {
		return cliutil.Usagef("--ledger: required")
	}
	_ = runID // present for future filtering; currently every entry is printed

	led, err := verify.LoadLedger(*ledgerPath)
	if err != nil {
		return err
	}
	for i := range led.Entries {
		e := &led.Entries[i]
		fmt.Fprintf(stdout, "[%04d] %s seal_date=%s prev=%s entry=%s\n",
			e.MerkleLeafIndex, e.EntryID, e.SealDate, short(e.PrevHash), short(e.EntryHash))
	}
	fmt.Fprintf(stdout, "seal:    date=%s leaves=%d signer=%s\n",
		led.Seal.SealDate, led.Seal.LeafCount, led.Seal.SigningKeyFingerprint)
	return nil
}

// short trims a base64 hash for the walk's compact display. The whole hash
// is preserved on disk; the verify command shows the full bytes when the
// recomputation differs.
func short(s string) string {
	if len(s) <= 12 {
		return s
	}
	return s[:8] + "…" + s[len(s)-4:]
}
